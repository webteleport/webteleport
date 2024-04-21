package transport

import (
	"context"
	"io"
	"net"
)

type Dialer interface {
	Dial(ctx context.Context, addr string) (Session, error)
}

type Listener interface {
	Listen(ctx context.Context, addr string) (net.Listener, error)
}

type Transport interface {
	Dialer
	Listener
}

type Session interface {
	Accept(context.Context) (Stream, error)
	Open(context.Context) (Stream, error)
	io.Closer
}

type Stream interface {
	net.Conn
}
