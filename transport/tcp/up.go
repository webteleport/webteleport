package tcp

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/spec"
	"github.com/webteleport/webteleport/transport/common"
)

var _ spec.Upgrader = (*Upgrader)(nil)

type Upgrader struct {
	net.Listener
	HOST string
}

func (s *Upgrader) Root() string {
	return s.HOST
}

func (s *Upgrader) Upgrade() (*spec.Request, error) {
	conn, err := s.Listener.Accept()
	if err != nil {
		return nil, err
	}

	ssn, err := common.YamuxClient(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("tcp creating yamux client failed: %w", err)
	}

	tssn := &TcpSession{ssn}

	stm0, err := tssn.Open(context.Background())
	if err != nil {
		return nil, fmt.Errorf("open stm0 error: %w", err)
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
