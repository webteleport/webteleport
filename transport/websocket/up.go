package websocket

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/btwiuse/wsconn"
	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/edge"
	"github.com/webteleport/webteleport/transport/common"
)

var _ edge.HTTPUpgrader = (*Upgrader)(nil)

type Upgrader struct {
	HOST string
	reqc chan *edge.Edge
}

func (s *Upgrader) Root() string {
	return s.HOST
}

func (s *Upgrader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsconn.Wrconn(w, r)
	if err != nil {
		slog.Warn(fmt.Errorf("websocket upgrade failed: %w", err).Error())
		return
	}

	ssn, err := common.YamuxClient(conn)
	if err != nil {
		slog.Warn(fmt.Errorf("websocket creating yamux client failed: %w", err).Error())
		return
	}

	tssn := &WebsocketSession{Session: ssn}
	tstm, err := tssn.Open(context.Background())
	if err != nil {
		slog.Warn(fmt.Errorf("websocket stm0 init failed: %w", err).Error())
		return
	}

	R := &edge.Edge{
		Session: tssn,
		Stream:  tstm,
		Path:    r.URL.Path,
		Values:  r.URL.Query(),
		RealIP:  utils.RealIP(r),
	}
	s.reqc <- R
}

func (s *Upgrader) Upgrade() (*edge.Edge, error) {
	if s.reqc == nil {
		s.reqc = make(chan *edge.Edge, 10)
	}
	r, ok := <-s.reqc
	if !ok {
		return nil, io.EOF
	}
	return r, nil
}
