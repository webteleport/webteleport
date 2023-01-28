package webteleport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/webteleport/utils"
	"github.com/webtransport/quic-go"
	"github.com/webtransport/quic-go/http3"
	"github.com/webtransport/webtransport-go"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

// Dial is a wrapper around webtransport.Dial with automatic HTTP/3 service discovery
func Dial(ctx context.Context, u *url.URL, hdr http.Header) (*webtransport.Session, error) {
	resp, err := http.Head(u.String())
	if err != nil {
		return nil, err
	}
	endpoints := utils.ExtractAltSvcEndpoints(resp.Header, "webteleport")
	if len(endpoints) == 0 {
		return nil, errors.New("webteleport service discovery failed: no endpoint in Alt-Svc header found")
	}
	alt := endpoints[0]
	addr := utils.Graft(u.Host, alt)
	d := &webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			QuicConfig: &quic.Config{
				MaxIncomingStreams: MaxIncomingStreams,
			},
		},
	}
	// we are dialing an HTTP/3 address, so it is guaranteed to be https
	un, _ := url.Parse(u.String())
	un.Scheme = "https"
	un.Host = addr
	_, session, err := d.Dial(ctx, un.String(), hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", un.String(), err)
	}
	return session, nil
}
