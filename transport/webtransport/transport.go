package webtransport

import (
	"context"
	"net"

	"github.com/webteleport/webteleport/transport"
)

type Transport struct{}

var _ transport.Transport = (*Transport)(nil)

func (t *Transport) Dial(ctx context.Context, addr string) (transport.Session, error) {
	tssn, err := DialWebtransport(ctx, addr, nil)
	if err != nil {
		return nil, err
	}
	return tssn, nil
}

func (t *Transport) Listen(ctx context.Context, addr string) (net.Listener, error) {
	return Listen(ctx, addr)
}
