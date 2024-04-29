package quic

import (
	"net"

	"github.com/quic-go/quic-go"
	"github.com/webteleport/transport"
)

var _ net.Conn = (*StreamConn)(nil)

var _ transport.Stream = (*StreamConn)(nil)

// StreamsConn wraps quic.Stream into net.Conn
type StreamConn struct {
	quic.Stream
	Session quic.Connection
}

// Close calls CancelRead to avoid memory leak, see
// - https://github.com/quic-go/quic-go/issues/3558
// - https://pkg.go.dev/github.com/quic-go/quic-go#Stream
func (sc *StreamConn) Close() error {
	sc.Stream.CancelRead(CancelRead)
	QuicGoConnsClosed.Add(1)
	return sc.Stream.Close()
}

var CancelRead quic.StreamErrorCode = 3558
var CancelWrite quic.StreamErrorCode = 3559

// LocalAddr is required to impl net.Conn
func (sc *StreamConn) LocalAddr() net.Addr { return sc.Session.LocalAddr() }

// RemoteAddr is required to impl net.Conn
func (sc *StreamConn) RemoteAddr() net.Addr { return sc.Session.RemoteAddr() }

func (sc *StreamConn) CloseRead() error {
	sc.Stream.CancelRead(CancelRead)
	return nil
}

func (sc *StreamConn) CloseWrite() error {
	sc.Stream.CancelWrite(CancelWrite)
	return nil
}
