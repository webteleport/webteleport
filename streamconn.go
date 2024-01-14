package webteleport

import (
	"net"

	"github.com/quic-go/quic-go"
)

var _ net.Conn = (*StreamConn)(nil)

// StreamsConn wraps webtransport.Stream into net.Conn
//
// TODO this should be part of github.com/webtransport/webtransport
type StreamConn struct {
	quic.Stream
	Session quic.Connection
}

// Close calls CancelRead to avoid memory leak, see
// - https://github.com/quic-go/quic-go/issues/3558
// - https://pkg.go.dev/github.com/quic-go/webtransport-go#Stream
func (sc *StreamConn) Close() error {
	sc.Stream.CancelRead(CancelRead)
	ConnsClosed.Add(1)
	return sc.Stream.Close()
}

var CancelRead quic.StreamErrorCode = 3558

// LocalAddr is required to impl net.Conn
func (sc *StreamConn) LocalAddr() net.Addr { return sc.Session.LocalAddr() }

// RemoteAddr is required to impl net.Conn
func (sc *StreamConn) RemoteAddr() net.Addr { return sc.Session.RemoteAddr() }
