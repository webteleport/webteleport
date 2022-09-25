package ufo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ebi-yade/altsvc-go"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

var MaxIncomingStreams = 65535

func extractH3(h http.Header) ([]string, bool) {
	line := h.Get("Alt-Svc")
	if line == "" {
		return []string{}, false
	}
	svcs, err := altsvc.Parse(line)
	if err != nil {
		log.Println(err)
		return []string{}, false
	}
	results := []string{}
	for _, svc := range svcs {
		if svc.ProtocolID == "h3" {
			results = append(results, svc.AltAuthority.Host+":"+svc.AltAuthority.Port)
		}
	}
	return results, len(results) > 0
}

// graft returns Host(base):Port(alt)
//
// assuming
// - base is host[:port]
// - alt is [host]:port
func graft(base, alt string) string {
	althost, altport, _ := strings.Cut(alt, ":")
	if altport == "" {
		// altport not found
		// it should never happen
		return base
	}
	if althost != "" {
		// alt is host:port
		// it is rare
		return alt
	}
	basehost, _, _ := strings.Cut(base, ":")
	return basehost + ":" + altport
}

// Dial is a wrapper around webtransport.Dial with automatic HTTP/3 service discovery
func Dial(ctx context.Context, u *url.URL, hdr http.Header) (*webtransport.Session, error) {
	resp, err := http.Head(u.String())
	if err != nil {
		return nil, err
	}
	endpoints, ok := extractH3(resp.Header)
	if !ok {
		return nil, errors.New("HTTP/3 service discovery failed: no Alt-Svc header found")
	}
	alt := endpoints[0]
	addr := graft(u.Host, alt)
	d := &webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			EnableDatagrams: true,
			QuicConfig: &quic.Config{
				MaxIncomingStreams: MaxIncomingStreams,
			},
		},
	}
	// we are dialing an HTTP/3 address, so it is guaranteed to be https://
	uri := "https://" + addr
	_, session, err := d.Dial(ctx, uri, hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", uri, err)
	}
	return session, nil
}
