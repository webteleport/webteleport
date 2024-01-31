package webtransport

import (
	"expvar"
)

var ConnsAccepted = expvar.NewInt("connsAccepted")
var ConnsOpened = expvar.NewInt("connsOpened")
var ConnsClosed = expvar.NewInt("connsClosed")
