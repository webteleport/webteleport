package websocket

import (
	"context"
	// "errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/yamux"
	"k0s.io/pkg/dial"
)

func DialWebsocket(_ctx context.Context, addr string, hdr http.Header) (*yamux.Session, error) {
	un, _ := url.Parse(addr)
	// we are dialing an HTTP/3 address, so it is guaranteed to be https
	un.Scheme = "https"
	params := un.Query()
	params.Add("x-websocket-upgrade", "1")
	un.RawQuery = params.Encode()
	conn, err := dial.Dial(un)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (WS): %w", un.String(), err)
	}
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating yamux.Client session: %w", err)
	}
	return session, nil
	// return &webtransportSession{session}, nil
}
