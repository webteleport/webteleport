package common

import (
	"context"
	"fmt"
	"net"

	"github.com/webteleport/transport"
)

var _ net.Addr = (*Listener)(nil)
var _ net.Listener = (*Listener)(nil)

// Listener implements [net.Listener]
type Listener struct {
	Session transport.Session
	Scheme  string
	Address string // host[:port][/path/]
}

// calling Accept returns a new [net.Conn]
func (l *Listener) Accept() (net.Conn, error) {
	streamConn, err := l.Session.Accept(context.Background())
	if err != nil {
		return nil, fmt.Errorf("accept: %w", err)
	}
	return streamConn, nil
}

func (l *Listener) Close() error {
	return l.Session.Close()
}

// Addr returns Listener itself which is an implementor of [net.Addr]
func (l *Listener) Addr() net.Addr {
	return l
}

// Network returns the protocol scheme, either http or https
func (l *Listener) Network() string {
	return l.Scheme
}

// String returns the host(:port) address of Listener
func (l *Listener) String() string {
	return l.Address
}
