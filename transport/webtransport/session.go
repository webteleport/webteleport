package webtransport

import (
	"context"
	"fmt"
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
	WebtransportConnsAccepted.Add(1)
	return &StreamConn{stm, s.Session}, nil
}

func (s *WebtransportSession) OpenStream(ctx context.Context) (transport.Stream, error) {
	stm, err := s.Session.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	// once ctx got cancelled, err is nil but stream is empty too
	// add the check to avoid returning empty stream
	if stm == nil {
		return nil, fmt.Errorf("stream is empty")
	}
	WebtransportConnsOpened.Add(1)
	return &StreamConn{stm, s.Session}, nil
}

func (s *WebtransportSession) Close() error {
	s.Session.CloseWithError(1337, "foobar")
	return http.ErrServerClosed
}
