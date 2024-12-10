package quic

import (
	"net"
	"time"

	"github.com/webteleport/webteleport/tunnel"
	"golang.org/x/net/quic"
)

var _ net.Conn = (*StreamConn)(nil)

var _ tunnel.Stream = (*StreamConn)(nil)

// StreamsConn wraps quic.Stream into net.Conn
type StreamConn struct {
	*quic.Stream
	Session *quic.Conn
}

// LocalAddr is required to impl net.Conn
func (sc *StreamConn) LocalAddr() net.Addr {
	return &netAddr{
		network:  "udp",
		hostport: sc.Session.LocalAddr().String(),
	}
}

// RemoteAddr is required to impl net.Conn
func (sc *StreamConn) RemoteAddr() net.Addr {
	return &netAddr{
		network:  "udp",
		hostport: sc.Session.RemoteAddr().String(),
	}
}

// SetDeadline is required to impl net.Conn
func (sc *StreamConn) SetDeadline(time.Time) error { return nil }

// SetReadDeadline is required to impl net.Conn
func (sc *StreamConn) SetReadDeadline(time.Time) error { return nil }

// SetWriteDeadline is required to impl net.Conn
func (sc *StreamConn) SetWriteDeadline(time.Time) error { return nil }

var _ net.Addr = (*netAddr)(nil)

type netAddr struct {
	network  string
	hostport string
}

func (n *netAddr) Network() string {
	return n.network
}

func (n *netAddr) String() string {
	return n.hostport
}
