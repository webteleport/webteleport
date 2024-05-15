package websocket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/transport/common"
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
	params.Add(UpgradeQuery, "1")
	u.RawQuery = params.Encode()
	return u.String(), nil
}

func Dial(ctx context.Context, addr string, hdr http.Header) (*WebsocketSession, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", addr, err)
	}

	conn, err := DialConn(ctx, addr, ModifyHeader(hdr))
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (WS): %w", u.Hostname(), utils.UnwrapInnermost(err))
	}

	session, err := common.YamuxServer(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating yamux server session: %w", utils.UnwrapInnermost(err))
	}

	return &WebsocketSession{session}, nil
}

func ModifyHeader(hdr http.Header) http.Header {
	if hdr == nil {
		hdr = make(http.Header)
	}
	hdr.Set(UpgradeHeader, "1")
	return hdr
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
