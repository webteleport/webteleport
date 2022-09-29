package server

import (
	"github.com/marten-seemann/webtransport-go"
)

type Session struct {
	*webtransport.Session
	Candidates []string
}
