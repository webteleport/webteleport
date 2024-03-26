package websocket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/yamux"
	"github.com/webteleport/utils"
	"k0s.io/pkg/dial"
)

func Merge(addr string, relayURL *url.URL) (string, error) {
	u, err := url.Parse(utils.AsURL(addr))
	if err != nil {
		return "", err
	}
	u.Host = relayURL.Host
	u.Scheme = relayURL.Scheme
	u.Path = relayURL.Path
	u.RawPath = relayURL.RawPath
	params := relayURL.Query()
	params.Add("x-websocket-upgrade", "1")
	u.RawQuery = params.Encode()
	return u.String(), nil
}

func DialWebsocket(_ctx context.Context, addr string, hdr http.Header) (*yamux.Session, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}
	conn, err := dial.Dial(u)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (WS): %w", u.Hostname(), utils.UnwrapInnermost(err))
	}
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating yamux.Client session: %w", utils.UnwrapInnermost(err))
	}
	return session, nil
	// return &webtransportSession{session}, nil
}
