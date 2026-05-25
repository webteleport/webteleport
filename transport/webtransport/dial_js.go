//go:build js

package webtransport

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/webtransport/webtransport"
)

func Dial(ctx context.Context, addr string, hdr http.Header) (*WebtransportSession, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}

	addr = applyHeaderQuery(addr, hdr)
	session, err := webtransport.Dial(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (WebTransport): %w", u.Hostname(), err)
	}

	return &WebtransportSession{Session: session}, nil
}

func ModifyHeader(hdr http.Header) http.Header {
	if hdr == nil {
		hdr = make(http.Header)
	}
	hdr.Set(UpgradeHeader, "1")
	return hdr
}

func applyHeaderQuery(addr string, hdr http.Header) string {
	hdr = ModifyHeader(hdr)
	u, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	q := u.Query()
	for key, values := range hdr {
		if key != UpgradeHeader || len(values) == 0 {
			continue
		}
		q.Set(UpgradeQuery, values[0])
	}
	u.RawQuery = q.Encode()
	return u.String()
}
