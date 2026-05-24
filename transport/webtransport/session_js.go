//go:build js

package webtransport

import (
	"context"
	"fmt"

	"github.com/webteleport/webteleport/tunnel"
	"github.com/webtransport/webtransport"
)

var _ tunnel.Session = (*WebtransportSession)(nil)

type WebtransportSession struct {
	Session *webtransport.Session
}

func (s *WebtransportSession) Accept(ctx context.Context) (tunnel.Stream, error) {
	conn, err := s.Session.Accept(ctx)
	if err != nil {
		return nil, err
	}
	stm, ok := conn.(*webtransport.Conn)
	if !ok {
		return nil, fmt.Errorf("unexpected webtransport stream type %T", conn)
	}
	WebtransportConnsAccepted.Add(1)
	return &StreamConn{Conn: stm}, nil
}

func (s *WebtransportSession) Open(ctx context.Context) (tunnel.Stream, error) {
	conn, err := s.Session.Open(ctx)
	if err != nil {
		return nil, err
	}
	stm, ok := conn.(*webtransport.Conn)
	if !ok {
		return nil, fmt.Errorf("unexpected webtransport stream type %T", conn)
	}
	WebtransportConnsOpened.Add(1)
	return &StreamConn{Conn: stm}, nil
}

func (s *WebtransportSession) Close() error {
	return s.Session.Close()
}

func (s *WebtransportSession) Context() context.Context {
	return s.Session.Context()
}
