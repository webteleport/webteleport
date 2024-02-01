package websocket

import (
	"context"
	// "errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
	// "github.com/webteleport/utils"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

func DialWebsocket(ctx context.Context, addr string, hdr http.Header) (*webtransport.Session, error) {
	d := &webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			QuicConfig: &quic.Config{
				MaxIncomingStreams: MaxIncomingStreams,
			},
		},
	}
	un, _ := url.Parse(addr)
	// we are dialing an HTTP/3 address, so it is guaranteed to be https
	un.Scheme = "https"
	params := un.Query()
	params.Add("x-webteleport-upgrade", "1")
	un.RawQuery = params.Encode()
	_, session, err := d.Dial(ctx, un.String(), hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", un.String(), err)
	}
	return session, nil
	// return &webtransportSession{session}, nil
}
