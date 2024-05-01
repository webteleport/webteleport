package webtransport

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	wt "github.com/quic-go/webtransport-go"
	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/spec"
)

var _ spec.HTTPUpgrader = (*WebtransportUpgrader)(nil)

type WebtransportUpgrader struct {
	HOST string
	reqc chan *spec.Request
	*wt.Server
}

func (s *WebtransportUpgrader) Root() string {
	return s.HOST
}

func (s *WebtransportUpgrader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ssn, err := s.Server.Upgrade(w, r)
	if err != nil {
		slog.Warn(fmt.Errorf("webtransport upgrade failed: %w", err).Error())
	}

	tssn := &WebtransportSession{Session: ssn}
	tstm, err := tssn.Open(context.Background())
	if err != nil {
		slog.Warn(fmt.Errorf("webtransport stm0 init failed: %w", err).Error())
	}

	R := &spec.Request{
		Session: tssn,
		Stream:  tstm,
		Path:    r.URL.Path,
		Values:  r.URL.Query(),
		RealIP:  utils.RealIP(r),
	}
	s.reqc <- R
}

func (s *WebtransportUpgrader) Upgrade() (*spec.Request, error) {
	if s.reqc == nil {
		s.reqc = make(chan *spec.Request, 10)
	}
	r, ok := <-s.reqc
	if !ok {
		return nil, io.EOF
	}
	return r, nil
}
