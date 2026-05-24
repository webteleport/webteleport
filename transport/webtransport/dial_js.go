//go:build js

package webtransport

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"syscall/js"

	"github.com/webteleport/utils"
)

func DialAddr(addr string, relayURL *url.URL) (string, error) {
	u, err := url.Parse(utils.AsURL(addr))
	if err != nil {
		return "", err
	}
	u.Scheme = "https"
	u.Path = relayURL.Path
	u.RawPath = relayURL.RawPath
	params := relayURL.Query()
	params.Set(UpgradeQuery, "1")
	u.RawQuery = params.Encode()
	return u.String(), nil
}

func Dial(ctx context.Context, addr string, hdr http.Header) (*WebtransportSession, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}

	webTransport := js.Global().Get("WebTransport")
	if webTransport.IsUndefined() || webTransport.IsNull() {
		return nil, fmt.Errorf("error dialing %s (WebTransport): WebTransport API is unavailable", u.Hostname())
	}

	addr = applyHeaderQuery(addr, hdr)
	transport := webTransport.New(addr)
	if _, err := awaitPromise(ctx, transport.Get("ready")); err != nil {
		return nil, fmt.Errorf("error dialing %s (WebTransport): %w", u.Hostname(), err)
	}

	return newWebtransportSession(transport, addr)
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
