package websocket

import (
	"net"

	"github.com/hashicorp/yamux"
	"github.com/webteleport/webteleport/tunnel"
)

var _ net.Conn = (*StreamConn)(nil)

var _ tunnel.Stream = (*StreamConn)(nil)

// StreamsConn wraps *yamux.Stream into net.Conn
type StreamConn struct {
	*yamux.Stream
}

func (sc *StreamConn) Close() error {
	WebsocketConnsClosed.Add(1)
	return sc.Stream.Close()
}

func (sc *StreamConn) CloseRead() error {
	return nil
}

func (sc *StreamConn) CloseWrite() error {
	return nil
}
