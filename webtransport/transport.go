package webtransport

import (
	"context"
	"net"

	"github.com/webteleport/webteleport/transport"
)

// WebtransportTransport is a transport that uses Webtransport.
type WebtransportTransport struct {
}

func NewTransport() transport.Transport {
	return &WebtransportTransport{}
}

var _ transport.Transport = &WebtransportTransport{}

func (t *WebtransportTransport) Dial(ctx context.Context, addr string) (transport.Session, error) {
	webtransportSession, err := DialWebtransport(ctx, addr, nil)
	if err != nil {
		return nil, err
	}
	return webtransportSession, nil
}

func (t *WebtransportTransport) Listen(ctx context.Context, addr string) (net.Listener, error) {
	return Listen(ctx, addr)
}
