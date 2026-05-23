//go:build js

package webtransport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"syscall/js"
	"time"

	"github.com/webteleport/webteleport/tunnel"
)

var _ net.Conn = (*StreamConn)(nil)

var _ tunnel.Stream = (*StreamConn)(nil)

type StreamConn struct {
	stream js.Value
	reader js.Value
	writer js.Value
	addr   net.Addr

	mu            sync.Mutex
	readBuf       []byte
	readDeadline  time.Time
	writeDeadline time.Time
	closed        bool
}

func newStreamConn(stream js.Value, addr string) *StreamConn {
	return &StreamConn{
		stream: stream,
		reader: stream.Get("readable").Call("getReader"),
		writer: stream.Get("writable").Call("getWriter"),
		addr:   streamAddr{value: addr},
	}
}

func (sc *StreamConn) Read(p []byte) (int, error) {
	sc.mu.Lock()
	if sc.closed {
		sc.mu.Unlock()
		return 0, net.ErrClosed
	}
	if len(sc.readBuf) > 0 {
		n := copy(p, sc.readBuf)
		sc.readBuf = sc.readBuf[n:]
		sc.mu.Unlock()
		return n, nil
	}
	deadline := sc.readDeadline
	sc.mu.Unlock()

	ctx, cancel := deadlineContext(deadline)
	defer cancel()

	result, err := awaitPromise(ctx, sc.reader.Call("read"))
	if err != nil {
		return 0, normalizeStreamError(err)
	}
	if result.Get("done").Bool() {
		return 0, io.EOF
	}

	chunk, err := jsValueBytes(result.Get("value"))
	if err != nil {
		return 0, err
	}

	sc.mu.Lock()
	sc.readBuf = append(sc.readBuf[:0], chunk...)
	n := copy(p, sc.readBuf)
	sc.readBuf = sc.readBuf[n:]
	sc.mu.Unlock()
	return n, nil
}

func (sc *StreamConn) Write(p []byte) (int, error) {
	sc.mu.Lock()
	if sc.closed {
		sc.mu.Unlock()
		return 0, net.ErrClosed
	}
	deadline := sc.writeDeadline
	sc.mu.Unlock()

	ctx, cancel := deadlineContext(deadline)
	defer cancel()

	chunk := js.Global().Get("Uint8Array").New(len(p))
	js.CopyBytesToJS(chunk, p)
	if _, err := awaitPromise(ctx, sc.writer.Call("write", chunk)); err != nil {
		return 0, normalizeStreamError(err)
	}
	return len(p), nil
}

func (sc *StreamConn) Close() error {
	sc.mu.Lock()
	if sc.closed {
		sc.mu.Unlock()
		return nil
	}
	sc.closed = true
	sc.mu.Unlock()

	WebtransportConnsClosed.Add(1)
	_, _ = awaitPromise(context.Background(), sc.reader.Call("cancel"))
	_, _ = awaitPromise(context.Background(), sc.writer.Call("close"))
	return nil
}

func (sc *StreamConn) LocalAddr() net.Addr { return sc.addr }

func (sc *StreamConn) RemoteAddr() net.Addr { return sc.addr }

func (sc *StreamConn) SetDeadline(t time.Time) error {
	if err := sc.SetReadDeadline(t); err != nil {
		return err
	}
	return sc.SetWriteDeadline(t)
}

func (sc *StreamConn) SetReadDeadline(t time.Time) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.readDeadline = t
	return nil
}

func (sc *StreamConn) SetWriteDeadline(t time.Time) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.writeDeadline = t
	return nil
}

func (sc *StreamConn) CloseRead() error {
	_, err := awaitPromise(context.Background(), sc.reader.Call("cancel"))
	return err
}

func (sc *StreamConn) CloseWrite() error {
	_, err := awaitPromise(context.Background(), sc.writer.Call("close"))
	return err
}

type streamAddr struct {
	value string
}

func (a streamAddr) Network() string { return "webtransport" }

func (a streamAddr) String() string { return a.value }

type promiseResult struct {
	value js.Value
	err   error
}

func awaitPromise(ctx context.Context, promise js.Value) (js.Value, error) {
	resultc := make(chan promiseResult, 1)

	var thenFunc, catchFunc js.Func
	thenFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer thenFunc.Release()
		defer catchFunc.Release()
		value := js.Undefined()
		if len(args) > 0 {
			value = args[0]
		}
		resultc <- promiseResult{value: value}
		return nil
	})
	catchFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer thenFunc.Release()
		defer catchFunc.Release()
		err := fmt.Errorf("webtransport promise rejected")
		if len(args) > 0 {
			err = jsError(args[0])
		}
		resultc <- promiseResult{err: err}
		return nil
	})

	promise.Call("then", thenFunc).Call("catch", catchFunc)

	select {
	case result := <-resultc:
		return result.value, result.err
	case <-ctx.Done():
		return js.Undefined(), ctx.Err()
	}
}

func jsError(v js.Value) error {
	if v.IsUndefined() || v.IsNull() {
		return fmt.Errorf("webtransport promise rejected")
	}
	if name := v.Get("name"); !name.IsUndefined() && !name.IsNull() {
		if msg := v.Get("message"); !msg.IsUndefined() && !msg.IsNull() {
			return fmt.Errorf("%s: %s", name.String(), msg.String())
		}
		return errors.New(name.String())
	}
	return errors.New(v.String())
}

func jsValueBytes(v js.Value) ([]byte, error) {
	if v.IsUndefined() || v.IsNull() {
		return nil, io.EOF
	}
	length := v.Get("byteLength")
	if length.IsUndefined() || length.IsNull() {
		return nil, fmt.Errorf("unexpected chunk type %q", v.Type())
	}
	buf := make([]byte, length.Int())
	js.CopyBytesToGo(buf, v)
	return buf, nil
}

func normalizeStreamError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return os.ErrDeadlineExceeded
	}
	return err
}

func deadlineContext(deadline time.Time) (context.Context, context.CancelFunc) {
	if deadline.IsZero() {
		return context.WithCancel(context.Background())
	}
	return context.WithDeadline(context.Background(), deadline)
}
