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

func Merge(addr string, relayURL *url.URL) (string, error) {
	u, err := url.Parse(utils.AsURL(addr))
	if err != nil {
		return "", err
	}
	// we are dialing an HTTP/3 address, so it is guaranteed to be https
	u.Scheme = "https"
	// u.Host = relayURL.Host
	u.Path = relayURL.Path
	u.RawPath = relayURL.RawPath
	params := relayURL.Query()
	params.Add("x-webtransport-upgrade", "1")
	u.RawQuery = params.Encode()
	return u.String(), nil
}

func DialWebtransport(ctx context.Context, addr string, hdr http.Header) (*webtransport.Session, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}
	dialer := &webtransport.Dialer{
		RoundTripper: &http3.RoundTripper{
			QuicConfig: &quic.Config{
				MaxIncomingStreams: MaxIncomingStreams,
			},
		},
	}
	_, session, err := dialer.Dial(ctx, addr, hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", u.Hostname(), utils.UnwrapInnermost(err))
	}
	return session, nil
	// return &webtransportSession{session}, nil
}
