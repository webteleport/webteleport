package websocket

import (
	"net"

	"github.com/hashicorp/yamux"
)

var _ net.Conn = (*StreamConn)(nil)

// StreamsConn wraps *yamux.Stream into net.Conn
//
// TODO this should be part of github.com/webtransport/webtransport
type StreamConn struct {
	*yamux.Stream
}

func (sc *StreamConn) Close() error {
	WebsocketConnsClosed.Add(1)
	return sc.Stream.Close()
}
