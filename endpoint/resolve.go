package endpoint

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ebi-yade/altsvc-go"
	"github.com/webteleport/utils"
)

type Endpoint struct {
	Protocol string // "websocket" or "webtransport"
	Addr     string // host:port
}

// ExtractAltSvcH3Endpoints reads Alt-Svc value
// returns a list of [host]:port endpoints
func ExtractAltSvcEndpoints(hostname, line, protocolId string) (endpoints []Endpoint) {
	svcs, err := altsvc.Parse(line)
	if err != nil {
		return
	}
	for _, svc := range svcs {
		if svc.ProtocolID != protocolId {
			continue
		}
		// host could be empty, port must not
		addr := svc.AltAuthority.Host + ":" + svc.AltAuthority.Port
		ep := Endpoint{
			Protocol: "webtransport",
			Addr:     utils.Graft(hostname, addr),
		}
		endpoints = append(endpoints, ep)
	}
	return
}

// Resolve gets all webteleport endpoints from Alt-Svc dns records / headers
func Resolve(u *url.URL) (endpoints []Endpoint) {
	endpoints = append(endpoints, eps(u.Hostname(), ENV("ALT_SVC"))...)
	endpoints = append(endpoints, eps(u.Hostname(), TXT(u.Host))...)
	endpoints = append(endpoints, eps(u.Hostname(), HEAD(u.String()))...)
	return
}

func eps(hostname string, altsvcs []string) (endpoints []Endpoint) {
	for _, v := range altsvcs {
		es := ExtractAltSvcEndpoints(hostname, v, "webteleport")
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
func HEAD(s string) (v []string) {
	resp, err := http.Head(utils.AsURL(s))
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
