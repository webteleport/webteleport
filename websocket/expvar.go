package websocket

import (
	"expvar"
)

var (
	WebsocketConnsAccepted = expvar.NewInt("websocket_conns_accepted")
	WebsocketConnsOpened   = expvar.NewInt("websocket_conns_opened")
	WebsocketConnsClosed   = expvar.NewInt("websocket_conns_closed")
)
