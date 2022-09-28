package webteleport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

// Dial is a wrapper around webtransport.Dial with automatic HTTP/3 service discovery
func Dial(ctx context.Context, u *url.URL, hdr http.Header) (*webtransport.Session, error) {
	resp, err := http.Head(u.String())
	if err != nil {
		return nil, err
	}
	endpoints := ExtractAltSvcEndpoints(resp.Header, "webteleport")
	if len(endpoints) == 0 {
		return nil, errors.New("webteleport service discovery failed: no endpoint in Alt-Svc header found")
	}
	alt := endpoints[0]
	addr := Graft(u.Host, alt)
	d := &webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			QuicConfig: &quic.Config{
				MaxIncomingStreams: MaxIncomingStreams,
			},
		},
	}
	// we are dialing an HTTP/3 address, so it is guaranteed to be https://
	uri := "https://" + addr + u.Path
	_, session, err := d.Dial(ctx, uri, hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", uri, err)
	}
	return session, nil
}
