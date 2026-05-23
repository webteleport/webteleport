//go:build js

package webtransport

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/webteleport/webteleport/tunnel"
)

var _ tunnel.Session = (*WebtransportSession)(nil)

type WebtransportSession struct {
	transport js.Value
	incoming  js.Value
	addr      string
}

func newWebtransportSession(transport js.Value, addr string) (*WebtransportSession, error) {
	incoming := transport.Get("incomingBidirectionalStreams")
	if incoming.IsUndefined() || incoming.IsNull() {
		return nil, fmt.Errorf("incoming bidirectional streams are unavailable")
	}
	return &WebtransportSession{
		transport: transport,
		incoming:  incoming.Call("getReader"),
		addr:      addr,
	}, nil
}

func (s *WebtransportSession) Accept(ctx context.Context) (tunnel.Stream, error) {
	result, err := awaitPromise(ctx, s.incoming.Call("read"))
	if err != nil {
		return nil, err
	}
	if result.Get("done").Bool() {
		return nil, context.Canceled
	}
	WebtransportConnsAccepted.Add(1)
	return newStreamConn(result.Get("value"), s.addr), nil
}

func (s *WebtransportSession) Open(ctx context.Context) (tunnel.Stream, error) {
	stream, err := awaitPromise(ctx, s.transport.Call("createBidirectionalStream"))
	if err != nil {
		return nil, err
	}
	if stream.IsUndefined() || stream.IsNull() {
		return nil, fmt.Errorf("stream is empty")
	}
	WebtransportConnsOpened.Add(1)
	return newStreamConn(stream, s.addr), nil
}

func (s *WebtransportSession) Close() error {
	s.transport.Call("close")
	if !(s.incoming.IsUndefined() || s.incoming.IsNull()) {
		_, _ = awaitPromise(context.Background(), s.incoming.Call("cancel"))
	}
	return nil
}

func (s *WebtransportSession) Context() context.Context {
	return context.TODO()
}
