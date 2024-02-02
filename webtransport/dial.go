package webtransport

import (
	"context"
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

func DialWebtransport(ctx context.Context, addr string, relayURL *url.URL, hdr http.Header) (*webtransport.Session, error) {
	d := &webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			QuicConfig: &quic.Config{
				MaxIncomingStreams: MaxIncomingStreams,
			},
		},
	}
	u, err := url.Parse(utils.AsURL(addr))
	if err != nil {
		return nil, err
	}
	// we are dialing an HTTP/3 address, so it is guaranteed to be https
	u.Scheme = "https"
	// u.Host = relayURL.Host
	u.Path = relayURL.Path
	u.RawPath = relayURL.RawPath
	params := relayURL.Query()
	params.Add("x-webtransport-upgrade", "1")
	u.RawQuery = params.Encode()
	_, session, err := d.Dial(ctx, u.String(), hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", u.Hostname(), utils.UnwrapInnermost(err))
	}
	return session, nil
	// return &webtransportSession{session}, nil
}
