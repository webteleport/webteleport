package transport

import (
	"context"
	"io"
	"net"
)

type Transport interface {
	Dial(ctx context.Context, addr string) (Session, error)
	Listen(ctx context.Context, addr string) (net.Listener, error)
}

type Session interface {
	AcceptStream(context.Context) (Stream, error)
	io.Closer
}

type Stream interface {
	net.Conn
}
