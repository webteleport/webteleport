package websocket

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/yamux"
	"github.com/webteleport/webteleport/tunnel"
)

var _ tunnel.Session = (*WebsocketSession)(nil)

type WebsocketSession struct {
	*yamux.Session
}

func (s *WebsocketSession) Accept(context.Context) (tunnel.Stream, error) {
	stm, err := s.Session.AcceptStream()
	if err != nil {
		return nil, err
	}
	WebsocketConnsAccepted.Add(1)
	return &StreamConn{stm}, nil
}

func (s *WebsocketSession) Open(context.Context) (tunnel.Stream, error) {
	stm, err := s.Session.OpenStream()
	if err != nil {
		return nil, err
	}
	// once ctx got cancelled, err is nil but stream is empty too
	// add the check to avoid returning empty stream
	if stm == nil {
		return nil, fmt.Errorf("stream is empty")
	}
	WebsocketConnsOpened.Add(1)
	return &StreamConn{stm}, nil
}

func (s *WebsocketSession) Close() error {
	s.Session.Close()
	return http.ErrServerClosed
}

func (s *WebsocketSession) Context() context.Context {
	return context.TODO()
}
