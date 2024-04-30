package tcp

import (
	"expvar"
)

var (
	TcpConnsAccepted = expvar.NewInt("tcp_conns_accepted")
	TcpConnsOpened   = expvar.NewInt("tcp_conns_opened")
	TcpConnsClosed   = expvar.NewInt("tcp_conns_closed")
)
