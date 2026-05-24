//go:build js

package webtransportjs

import (
	"context"
	"fmt"
	"net/url"
	"syscall/js"
)

func Dial(ctx context.Context, addr string) (*Session, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}

	webTransport := js.Global().Get("WebTransport")
	if webTransport.IsUndefined() || webTransport.IsNull() {
		return nil, fmt.Errorf("error dialing %s (WebTransport): WebTransport API is unavailable", u.Hostname())
	}

	transport := webTransport.New(addr)
	if _, err := awaitPromise(ctx, transport.Get("ready")); err != nil {
		return nil, fmt.Errorf("error dialing %s (WebTransport): %w", u.Hostname(), err)
	}

	return newSession(transport, addr)
}
