package websocket

import (
	"context"
	"net"

	"github.com/webteleport/webteleport/transport"
)

// WebsocketTransport is a transport that uses Websocket.
type WebsocketTransport struct {
}

func New() transport.Transport {
	return &WebsocketTransport{}
}

var _ transport.Transport = &WebsocketTransport{}

func (t *WebsocketTransport) Dial(ctx context.Context, addr string) (transport.Session, error) {
	websocketSession, err := DialWebsocket(ctx, addr, nil)
	if err != nil {
		return nil, err
	}
	return websocketSession, nil
}

func (t *WebsocketTransport) Listen(ctx context.Context, addr string) (net.Listener, error) {
	_ = ctx
	_ = addr
	return nil, nil
}
