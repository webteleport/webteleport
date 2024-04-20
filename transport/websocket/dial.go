package websocket

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

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

func YamuxConfig() *yamux.Config {
	c := yamux.DefaultConfig()
	if os.Getenv("YAMUX_LOG") == "" {
		c.LogOutput = io.Discard
	}
	return c
}

func YamuxReverseGender(conn io.ReadWriteCloser, r *http.Request) (string, *yamux.Session, error) {
	// default gender of new clients is "server", so reverse it to "client"
	config := YamuxConfig()
	session, err := yamux.Client(conn, config)
	return "client", session, err
}

func Dial(ctx context.Context, addr string, hdr http.Header) (*WebsocketSession, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}

	conn, err := DialConn(ctx, addr, hdr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (WS): %w", u.Hostname(), utils.UnwrapInnermost(err))
	}

	config := YamuxConfig()
	session, err := yamux.Server(conn, config)
	if err != nil {
		return nil, fmt.Errorf("error creating yamux server session: %w", utils.UnwrapInnermost(err))
	}

	return &WebsocketSession{session}, nil
}

func DialConn(ctx context.Context, addr string, hdr http.Header) (conn net.Conn, err error) {
	wsconn, _, err := websocket.Dial(
		ctx,
		addr,
		dialOptions(hdr),
	)
	if err != nil {
		return nil, err
	}

	return websocket.NetConn(context.Background(), wsconn, websocket.MessageBinary), nil
}
