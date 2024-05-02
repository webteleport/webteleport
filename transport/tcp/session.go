package tcp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/yamux"
	"github.com/webteleport/webteleport/tunnel"
)

type TcpSession struct {
	*yamux.Session
}

func (s *TcpSession) Accept(context.Context) (tunnel.Stream, error) {
	stm, err := s.Session.AcceptStream()
	if err != nil {
		return nil, err
	}
	TcpConnsAccepted.Add(1)
	return &StreamConn{stm}, nil
}

func (s *TcpSession) Open(context.Context) (tunnel.Stream, error) {
	stm, err := s.Session.OpenStream()
	if err != nil {
		return nil, err
	}
	// once ctx got cancelled, err is nil but stream is empty too
	// add the check to avoid returning empty stream
	if stm == nil {
		return nil, fmt.Errorf("stream is empty")
	}
	TcpConnsOpened.Add(1)
	return &StreamConn{stm}, nil
}

func (s *TcpSession) Close() error {
	s.Session.Close()
	return http.ErrServerClosed
}
