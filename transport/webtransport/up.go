package webtransport

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	wt "github.com/quic-go/webtransport-go"
	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/edge"
	"github.com/webteleport/webteleport/transport/common"
)

var (
	UpgradeQuery  = "x-webtransport-upgrade"
	UpgradeHeader = http.CanonicalHeaderKey(UpgradeQuery)
)

var _ edge.HTTPUpgrader = (*Upgrader)(nil)

type Upgrader struct {
	reqc chan *edge.Edge
	*wt.Server
	common.RootPatterns
}

func (s *Upgrader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ssn, err := s.Server.Upgrade(w, r)
	if err != nil {
		slog.Warn(fmt.Errorf("webtransport upgrade failed: %w", err).Error())
	}

	tssn := &WebtransportSession{Session: ssn}
	tstm, err := tssn.Open(context.Background())
	if err != nil {
		slog.Warn(fmt.Errorf("webtransport stm0 init failed: %w", err).Error())
	}

	R := &edge.Edge{
		Session: tssn,
		Stream:  tstm,
		Path:    r.URL.Path,
		Header:  r.Header,
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
