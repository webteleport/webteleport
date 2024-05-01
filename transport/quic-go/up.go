package quic

import (
	"context"
	"fmt"
	"net/url"

	"github.com/quic-go/quic-go"
	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/spec"
	"github.com/webteleport/webteleport/transport/common"
)

var _ spec.Upgrader = (*QuicGoUpgrader)(nil)

type QuicGoUpgrader struct {
	*quic.Listener
	HOST string
}

func (s *QuicGoUpgrader) Root() string {
	return s.HOST
}

func (s *QuicGoUpgrader) Upgrade() (*spec.Request, error) {
	conn, err := s.Listener.Accept(context.Background())
	if err != nil {
		return nil, fmt.Errorf("accept error: %w", err)
	}

	tssn := &QuicSession{conn}

	stm0, err := tssn.Accept(context.Background())
	if err != nil {
		return nil, fmt.Errorf("accept stm0 error: %w", err)
	}

	ruri, err := common.ReadLine(stm0)
	if err != nil {
		return nil, fmt.Errorf("read request uri error: %w", err)
	}

	u, err := url.ParseRequestURI(ruri)
	if err != nil {
		return nil, fmt.Errorf("parse request uri error: %w", err)
	}

	R := &spec.Request{
		Session: tssn,
		Stream:  stm0,
		Path:    u.Path,
		Values:  u.Query(),
		RealIP:  utils.StripPort(conn.RemoteAddr().String()),
	}
	return R, nil
}
