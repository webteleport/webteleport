//go:build js

package webtransport

import (
	"net"

	"github.com/webteleport/webteleport/tunnel"
	"github.com/webtransport/webtransport"
)

var _ net.Conn = (*StreamConn)(nil)

var _ tunnel.Stream = (*StreamConn)(nil)

type StreamConn struct {
	*webtransport.Conn
}

func (sc *StreamConn) Close() error {
	StreamMetrics.Closed.Add(1)
	return sc.Conn.Close()
}
