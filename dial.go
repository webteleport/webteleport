package webteleport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
	"github.com/webteleport/utils"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

// Dial is a wrapper around webtransport.Dial with automatic HTTP/3 service discovery
func Dial(ctx context.Context, u *url.URL, hdr http.Header) (*webtransport.Session, error) {
	endpoints := Resolve(u)
	if len(endpoints) == 0 {
		return nil, errors.New("service discovery failed: no webteleport endpoints found in Alt-Svc records / headers")
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
	params := un.Query()
	params.Add("x-webteleport-upgrade", "1")
	un.RawQuery = params.Encode()
	_, session, err := d.Dial(ctx, un.String(), hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", un.String(), err)
	}
	return session, nil
}
