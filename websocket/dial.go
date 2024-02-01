package websocket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/yamux"
	"k0s.io/pkg/dial"
)

func DialWebsocket(_ctx context.Context, addr string, up *url.URL, hdr http.Header) (*yamux.Session, error) {
	un, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	// we are dialing an HTTP/3 address, so it is guaranteed to be https
	un.Host = up.Host
	un.Scheme = up.Scheme
	un.Path = up.Path
	un.RawPath = up.RawPath
	params := up.Query()
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
