package quic

import (
	"expvar"
)

var (
	GoQuicConnsAccepted = expvar.NewInt("go_quic_conns_accepted")
	GoQuicConnsOpened   = expvar.NewInt("go_quic_conns_opened")
	GoQuicConnsClosed   = expvar.NewInt("go_quic_conns_closed")
)
