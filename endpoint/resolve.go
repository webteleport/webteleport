package endpoint

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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
			// TODO: support ALT_SVC keys like "webteleport-ws" for specifying ws endpoints on non-bootstrap ports
			// Until then, we assume that all endpoints are websockets by default
			// Protocol: "websocket",
			Protocol: "webtransport",
			Addr:     utils.Graft(hostname, addr),
		}
		endpoints = append(endpoints, ep)
	}
	return
}

// Resolve gets all webteleport endpoints from Alt-Svc dns records / headers
// fallback to websocket so that resolve always returns at least one endpoint
func Resolve(u *url.URL) (endpoints []Endpoint) {
	endpoints = append(endpoints, eps(u.Hostname(), ENV("ALT_SVC"))...)
	endpoints = append(endpoints, eps(u.Hostname(), HEAD(u.Host))...)
	if len(endpoints) == 0 {
		endpoints = append(endpoints, Endpoint{
			Protocol: "websocket",
			Addr:     u.Hostname(),
		})
	}
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
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodHead, utils.AsURL(s), nil)
	if err != nil {
		slog.Warn(fmt.Sprintf("http req error: %v", err))
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Warn(fmt.Sprintf("http head error: %v", err))
		return
	}
	return resp.Header.Values("Alt-Svc")
}

// TXT gets all altsvcs from TXT records of given host
func TXT(h string) []string {
	// skip common local addresses
	switch h {
	case "localhost", "127.0.0.1", "::1":
		return []string{}
	}
	txts, err := utils.LookupHostTXT(h, "1.1.1.1:53")
	if err != nil {
		slog.Warn(fmt.Sprintf("dns lookup error: %s: %v", h, err))
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
