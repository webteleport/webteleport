package quic

import (
	"context"
	"fmt"

	"github.com/webteleport/webteleport/tunnel"
	"golang.org/x/net/quic"
)

var _ tunnel.Session = (*QuicSession)(nil)

type QuicSession struct {
	Session *quic.Conn
}

func (s *QuicSession) Accept(ctx context.Context) (tunnel.Stream, error) {
	stm, err := s.Session.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	GoQuicConnsAccepted.Add(1)
	return &StreamConn{stm, s.Session}, nil
}

func (s *QuicSession) Open(ctx context.Context) (tunnel.Stream, error) {
	stm, err := s.Session.NewStream(ctx)
	if err != nil {
		return nil, err
	}
	// once ctx got cancelled, err is nil but stream is empty too
	// add the check to avoid returning empty stream
	if stm == nil {
		return nil, fmt.Errorf("stream is empty")
	}
	GoQuicConnsOpened.Add(1)
	return &StreamConn{stm, s.Session}, nil
}

func (s *QuicSession) Close() error {
	return s.Session.Close()
}

func (s *QuicSession) Context() context.Context {
	return context.TODO()
}
