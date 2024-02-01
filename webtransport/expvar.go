package webtransport

import (
	"expvar"
)

var (
	WebtransportConnsAccepted = expvar.NewInt("webtransport_conns_accepted")
	WebtransportConnsOpened   = expvar.NewInt("webtransport_conns_opened")
	WebtransportConnsClosed   = expvar.NewInt("webtransport_conns_closed")
)
