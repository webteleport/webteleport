package quic

import (
	"expvar"
)

var (
	QuicGoConnsAccepted = expvar.NewInt("quic_go_conns_accepted")
	QuicGoConnsOpened   = expvar.NewInt("quic_go_conns_opened")
	QuicGoConnsClosed   = expvar.NewInt("quic_go_conns_closed")
)
