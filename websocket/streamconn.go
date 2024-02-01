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

func NewAcceptedConn(s *yamux.Stream) net.Conn {
	WebsocketConnsAccepted.Add(1)
	return &StreamConn{s}
}

func NewOpenedConn(s *yamux.Stream) net.Conn {
	WebsocketConnsOpened.Add(1)
	return &StreamConn{s}
}

func (sc *StreamConn) Close() error {
	WebsocketConnsClosed.Add(1)
	return sc.Stream.Close()
}
