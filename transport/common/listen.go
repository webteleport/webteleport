package common

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/webteleport/webteleport/tunnel"
)

var _ net.Addr = (*Listener)(nil)
var _ net.Listener = (*Listener)(nil)

// Listener implements [net.Listener]
type Listener struct {
	Session tunnel.Session
	Scheme  string
	Address string // host[:port][/path/]
	Relay   *url.URL
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

// Addr returns Listener's full address.
// When Relay is set, Addr().String() joins Address with Relay.Host.
func (l *Listener) Addr() net.Addr {
	if l.Relay == nil {
		return l
	}
	return &subdomainAddr{l}
}

// subdomainAddr joins Address and Relay.Host with a dot.
type subdomainAddr struct{ l *Listener }

func (a *subdomainAddr) Network() string { return a.l.Network() }
func (a *subdomainAddr) String() string  { return a.l.Address + "." + a.l.Relay.Host }

// Network returns the protocol scheme, either http or https
func (l *Listener) Network() string {
	return l.Scheme
}

// String returns the host(:port) address of Listener
func (l *Listener) String() string {
	return l.Address
}
