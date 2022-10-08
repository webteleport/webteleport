package server

import (
	"context"
	"fmt"
	"net"

	"github.com/marten-seemann/webtransport-go"
	"github.com/webteleport/webteleport"
)

type Session struct {
	*webtransport.Session
	Controller net.Conn
}

func (ssn *Session) InitController(ctx context.Context) error {
	if ssn.Controller != nil {
		return nil
	}
	stm0, err := ssn.OpenConn(ctx)
	if err != nil {
		return err
	}
	ssn.Controller = stm0
	return nil
}

func (ssn *Session) OpenConn(ctx context.Context) (net.Conn, error) {
	// when there is a timeout, it still panics before MARK
	//
	// ctx, _ = context.WithTimeout(ctx, 3*time.Second)
	//
	// turns out the stream is empty so need to check stream == nil
	stream, err := ssn.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	// once ctx got cancelled, err is nil but stream is empty too
	// add the check to avoid returning empty stream
	if stream == nil {
		return nil, fmt.Errorf("stream is empty")
	}
	// log.Println(`MARK`, stream)
	// MARK
	conn := &webteleport.StreamConn{stream, ssn.Session}
	return conn, nil
}
