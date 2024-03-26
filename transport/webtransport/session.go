package webtransport

import (
	"context"
	"net/http"

	"github.com/quic-go/webtransport-go"
	"github.com/webteleport/webteleport/transport"
)

type WebtransportSession struct {
	*webtransport.Session
}

func (s *WebtransportSession) AcceptStream(ctx context.Context) (transport.Stream, error) {
	stm, err := s.Session.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}
	return &StreamConn{stm, s.Session}, nil
}

func (s *WebtransportSession) OpenStream(ctx context.Context) (transport.Stream, error) {
	stm, err := s.Session.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	return &StreamConn{stm, s.Session}, nil
}

func (s *WebtransportSession) Close() error {
	s.Session.CloseWithError(1337, "foobar")
	return http.ErrServerClosed
}
