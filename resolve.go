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

// HEAD gets all altsvcs from Alt-Svc header values of given url
func HEAD(s string) []string {
	resp, err := http.Head(s)
	if err != nil {
		slog.Warn(fmt.Sprintf("http head error: %v", err))
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
