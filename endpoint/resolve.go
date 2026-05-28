package endpoint

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"slices"

	"github.com/ebi-yade/altsvc-go"
	"github.com/webteleport/utils"
)

type Endpoint struct {
	Protocol string // "tcp", "quic", "quic-go", "websocket", "webtransport", "net-quic", etc.
	Addr     string // host:port
}

// ExtractWebteleport parses Alt-Svc header lines, keeps only
// entries with protocol ID "webteleport", and returns endpoints.
func ExtractWebteleport(hostname string, lines ...string) (endpoints []Endpoint) {
	for _, line := range lines {
		svcs, err := altsvc.Parse(line)
		if err != nil {
			slog.Warn("altsvc parse error", "line", line, "error", err)
			continue
		}
		for _, svc := range svcs {
			if svc.ProtocolID != "webteleport" {
				continue
			}
			addr := svc.AltAuthority.Host + ":" + svc.AltAuthority.Port
			ep := Endpoint{
				Protocol: "webtransport",
				Addr:     utils.Graft(hostname, addr),
			}
			endpoints = append(endpoints, ep)
		}
	}
	return
}

// Resolve discovers webteleport endpoints from Alt-Svc via env (ALT_SVC) and HTTP HEAD request.
// Always appends a websocket endpoint on the original host as the final option.
// For raw protocols (tcp, quic, etc.), callers can construct endpoints directly.
func Resolve(ctx context.Context, u *url.URL) (endpoints []Endpoint) {
	if ctx == nil {
		ctx = context.Background()
	}
	endpoints = ExtractWebteleport(
		u.Hostname(),
		slices.Concat(
			AltSvcFromEnv("ALT_SVC"),
			AltSvcFromHEAD(ctx, u.String()),
		)...,
	)
	endpoints = append(endpoints, Endpoint{
		Protocol: "websocket",
		Addr:     u.Hostname(),
	})
	return
}

// AltSvcFromEnv reads the Alt-Svc value from the given environment variable.
func AltSvcFromEnv(key string) []string {
	v, ok := os.LookupEnv(key)
	if ok {
		return []string{v}
	}
	return nil
}

// AltSvcFromHEAD fetches Alt-Svc headers from the given URL via an HTTP HEAD request.
func AltSvcFromHEAD(ctx context.Context, rawurl string) []string {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawurl, nil)
	if err != nil {
		slog.Warn("http req error", "error", utils.UnwrapInnermost(err))
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		slog.Warn("http head error", "error", utils.UnwrapInnermost(err))
		return nil
	}
	defer resp.Body.Close()
	return resp.Header.Values("Alt-Svc")
}
