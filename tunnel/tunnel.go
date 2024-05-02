package tunnel

import (
	"context"
	"net"
)

type Dialer interface {
	Dial(ctx context.Context, addr string) (Session, error)
}

type Listener interface {
	Listen(ctx context.Context, addr string) (net.Listener, error)
}

type Session interface {
	Accept(context.Context) (Stream, error)
	Open(context.Context) (Stream, error)
	Close() error
}

type Stream interface {
	net.Conn
}

type Transport interface {
	Dialer
	Listener
}
