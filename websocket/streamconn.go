package websocket

import (
	"net"

	"github.com/hashicorp/yamux"
)

var _ net.Conn = (*StreamConn)(nil)

// StreamsConn wraps *yamux.Stream into net.Conn
type StreamConn struct {
	*yamux.Stream
}

func NewConn(s *yamux.Stream) net.Conn {
	WebsocketConnsOpened.Add(1)
	return &StreamConn{s}
}

func (sc *StreamConn) Close() error {
	WebsocketConnsClosed.Add(1)
	return sc.Stream.Close()
}
