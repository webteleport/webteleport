package webtransport

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/webtransport-go"
	"github.com/webteleport/utils"
)

// 2^60 == 1152921504606846976
var MaxIncomingStreams int64 = 1 << 60

var QUICConfig = &quic.Config{
	EnableDatagrams:    true,
	MaxIncomingStreams: MaxIncomingStreams,
}

func DialAddr(addr string, relayURL *url.URL) (string, error) {
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

func Dial(ctx context.Context, addr string, hdr http.Header) (*WebtransportSession, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}
	dialer := &webtransport.Dialer{
		QUICConfig: QUICConfig,
	}
	_, session, err := dialer.Dial(ctx, addr, ModifyHeader(hdr))
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (UDP): %w", u.Hostname(), utils.UnwrapInnermost(err))
	}
	return &WebtransportSession{session}, nil
}

func ModifyHeader(hdr http.Header) http.Header {
	if hdr == nil {
		hdr = make(http.Header)
	}
	hdr.Set(UpgradeHeader, "1")
	return hdr
}
