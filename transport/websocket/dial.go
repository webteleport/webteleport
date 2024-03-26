package websocket

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/hashicorp/yamux"
	"github.com/webteleport/utils"
	"nhooyr.io/websocket"
)

func DialAddr(addr string, relayURL *url.URL) (string, error) {
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

func DialWebsocket(ctx context.Context, addr string, hdr http.Header) (*WebsocketSession, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}
	conn, err := Dial(ctx, addr, hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (WS): %w", u.Hostname(), utils.UnwrapInnermost(err))
	}
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating yamux.Client session: %w", utils.UnwrapInnermost(err))
	}
	return &WebsocketSession{session}, nil
}

func Dial(ctx context.Context, addr string, hdr http.Header) (conn net.Conn, err error) {
	wsconn, _, err := websocket.Dial(
		ctx,
		addr,
		&websocket.DialOptions{
			HTTPClient: &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
			HTTPHeader: hdr,
		},
	)
	if err != nil {
		return nil, err
	}

	return websocket.NetConn(context.Background(), wsconn, websocket.MessageBinary), nil
}
