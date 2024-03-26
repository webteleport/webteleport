package websocket

import (
	"context"
	"net/http"

	"github.com/hashicorp/yamux"
	"github.com/webteleport/webteleport/transport"
)

type WebsocketSession struct {
	*yamux.Session
}

func (s *WebsocketSession) AcceptStream(context.Context) (transport.Stream, error) {
	stm, err := s.Session.AcceptStream()
	if err != nil {
		return nil, err
	}
	return &StreamConn{stm}, nil
}

func (s *WebsocketSession) OpenStream(context.Context) (transport.Stream, error) {
	stm, err := s.Session.OpenStream()
	if err != nil {
		return nil, err
	}
	return &StreamConn{stm}, nil
}

func (s *WebsocketSession) Close() error {
	s.Session.Close()
	return http.ErrServerClosed
}
