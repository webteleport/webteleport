//go:build js

package webtransport

import (
	"net"

	"github.com/webteleport/webteleport/tunnel"
	"github.com/webteleport/webteleport/webtransportjs"
)

var _ net.Conn = (*StreamConn)(nil)

var _ tunnel.Stream = (*StreamConn)(nil)

type StreamConn struct {
	*webtransportjs.Conn
}

func (sc *StreamConn) Close() error {
	WebtransportConnsClosed.Add(1)
	return sc.Conn.Close()
}
