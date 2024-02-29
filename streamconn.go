package webteleport

import (
	"net"
	"time"

	"golang.org/x/net/quic"
)

var _ net.Conn = (*StreamConn)(nil)

// StreamsConn wraps webtransport.Stream into net.Conn
//
// TODO this should be part of github.com/webtransport/webtransport
type StreamConn struct {
	*quic.Stream
	Session *quic.Conn
}

// LocalAddr is required to impl net.Conn
func (sc *StreamConn) LocalAddr() net.Addr { return sc.Session.LocalAddr() }

// RemoteAddr is required to impl net.Conn
func (sc *StreamConn) RemoteAddr() net.Addr { return sc.Session.RemoteAddr() }

// SetDeadline is required to impl net.Conn
func (sc *StreamConn) SetDeadline(time.Time) error { return nil }

// SetReadDeadline is required to impl net.Conn
func (sc *StreamConn) SetReadDeadline(time.Time) error { return nil }

// SetWriteDeadline is required to impl net.Conn
func (sc *StreamConn) SetWriteDeadline(time.Time) error { return nil }
