package webteleport

import (
	"net"

	"github.com/quic-go/webtransport-go"
)

var _ net.Conn = (*StreamConn)(nil)

// StreamsConn wraps webtransport.Stream into net.Conn
//
// TODO this should be part of github.com/webtransport/webtransport
type StreamConn struct {
	webtransport.Stream
	Session *webtransport.Session
}

// LocalAddr is required to impl net.Conn
func (sc *StreamConn) LocalAddr() net.Addr { return sc.Session.LocalAddr() }

// RemoteAddr is required to impl net.Conn
func (sc *StreamConn) RemoteAddr() net.Addr { return sc.Session.RemoteAddr() }
