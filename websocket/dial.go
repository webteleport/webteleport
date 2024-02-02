package websocket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/yamux"
	"k0s.io/pkg/dial"
)

func DialWebsocket(_ctx context.Context, addr string, relayURL *url.URL, hdr http.Header) (*yamux.Session, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	// we are dialing an HTTP/3 address, so it is guaranteed to be https
	u.Host = relayURL.Host
	u.Scheme = relayURL.Scheme
	u.Path = relayURL.Path
	u.RawPath = relayURL.RawPath
	params := relayURL.Query()
	params.Add("x-websocket-upgrade", "1")
	u.RawQuery = params.Encode()
	conn, err := dial.Dial(u)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (WS): %w", u.String(), err)
	}
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating yamux.Client session: %w", err)
	}
	return session, nil
	// return &webtransportSession{session}, nil
}
