package websocket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/btwiuse/version"
	"github.com/coder/websocket"
	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/transport/common"
)

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
	ua := fmt.Sprintf("webteleport/%s (%s)", version.GitVersionString, version.GitCommitString)
	hdr.Set("User-Agent", ua)
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
