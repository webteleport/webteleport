package quic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/webteleport/webteleport/tunnel"
)

var _ tunnel.Session = (*QuicSession)(nil)

type QuicSession struct {
	Session quic.Connection
}

func (s *QuicSession) Accept(ctx context.Context) (tunnel.Stream, error) {
	stm, err := s.Session.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	QuicGoConnsAccepted.Add(1)
	return &StreamConn{stm, s.Session}, nil
}

func (s *QuicSession) Open(ctx context.Context) (tunnel.Stream, error) {
	stm, err := s.Session.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	// once ctx got cancelled, err is nil but stream is empty too
	// add the check to avoid returning empty stream
	if stm == nil {
		return nil, fmt.Errorf("stream is empty")
	}
	QuicGoConnsOpened.Add(1)
	return &StreamConn{stm, s.Session}, nil
}

func (s *QuicSession) Close() error {
	s.Session.CloseWithError(1337, "foobar")
	return http.ErrServerClosed
}

func (s *QuicSession) Context() context.Context {
	return context.TODO()
}
