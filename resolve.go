package webteleport

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/webteleport/utils"
)

// Resolve gets all webteleport endpoints from Alt-Svc dns records / headers
func Resolve(u *url.URL) (endpoints []string) {
	endpoints = append(endpoints, eps(ENV("ALT_SVC"))...)
	endpoints = append(endpoints, eps(TXT(u.Host))...)
	endpoints = append(endpoints, eps(HEAD(u.String()))...)
	return
}

func eps(altsvcs []string) (endpoints []string) {
	for _, v := range altsvcs {
		es := utils.ExtractAltSvcEndpoints(v, "webteleport")
		endpoints = append(endpoints, es...)
	}
	return
}

// ENV gets altsvc value from env
func ENV(s string) []string {
	v, ok := os.LookupEnv(s)
	if ok {
		return []string{v}
	}
	return []string{}
}

// AsURL expands :port and hostname to http://localhost:port & http://hostname respectively
func AsURL(s string) string {
	if strings.HasPrefix(s, ":") {
		s = "http://localhost" + s
	} else if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		s = "http://" + s
	}
	return s
}

// HEAD gets all altsvcs from Alt-Svc header values of given url
func HEAD(s string) (v []string) {
	resp, err := http.Head(AsURL(s))
	if err != nil {
		slog.Warn(fmt.Sprintf("http head error: %v", err))
		return
	}
	return resp.Header.Values("Alt-Svc")
}

// TXT gets all altsvcs from TXT records of given host
func TXT(h string) []string {
	txts, err := utils.LookupHostTXT(h, "1.1.1.1:53")
	if err != nil {
		slog.Warn(fmt.Sprintf("dns lookup error: %v", err))
	}
	return altsvcLines(txts)
}

func altsvcLines(txts []string) []string {
	const prefix = "Alt-Svc: "

	altsvcs := []string{}

	for _, txt := range txts {
		// Case insensitive prefix match. See Issue 22736.
		if len(txt) < len(prefix) || !strings.EqualFold(txt[:len(prefix)], prefix) {
			continue
		}
		altsvcs = append(altsvcs, txt[len(prefix):])
	}

	return altsvcs
}
