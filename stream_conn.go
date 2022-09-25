package ufo

import (
	"net"

	"github.com/marten-seemann/webtransport-go"
)

var _ net.Conn = (*StreamConn)(nil)

type StreamConn struct {
	webtransport.Stream
	Session *webtransport.Session
}

func (sc *StreamConn) LocalAddr() net.Addr  { return sc.Session.LocalAddr() }
func (sc *StreamConn) RemoteAddr() net.Addr { return sc.Session.RemoteAddr() }
