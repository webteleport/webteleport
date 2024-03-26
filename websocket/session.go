package websocket

import (
	"context"

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
