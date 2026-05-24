//go:build js

package webtransportjs

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
)

var _ net.Conn = (*Conn)(nil)

type Conn struct {
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

func newConn(stream js.Value, addr string) *Conn {
	return &Conn{
		stream: stream,
		reader: stream.Get("readable").Call("getReader"),
		writer: stream.Get("writable").Call("getWriter"),
		addr:   streamAddr{value: addr},
	}
}

func (c *Conn) Read(p []byte) (int, error) {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return 0, net.ErrClosed
	}
	if len(c.readBuf) > 0 {
		n := copy(p, c.readBuf)
		c.readBuf = c.readBuf[n:]
		c.mu.Unlock()
		return n, nil
	}
	deadline := c.readDeadline
	c.mu.Unlock()

	ctx, cancel := deadlineContext(deadline)
	defer cancel()

	result, err := awaitPromise(ctx, c.reader.Call("read"))
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

	c.mu.Lock()
	c.readBuf = chunk
	n := copy(p, c.readBuf)
	c.readBuf = c.readBuf[n:]
	c.mu.Unlock()
	return n, nil
}

func (c *Conn) Write(p []byte) (int, error) {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return 0, net.ErrClosed
	}
	deadline := c.writeDeadline
	c.mu.Unlock()

	ctx, cancel := deadlineContext(deadline)
	defer cancel()

	chunk := js.Global().Get("Uint8Array").New(len(p))
	js.CopyBytesToJS(chunk, p)
	if _, err := awaitPromise(ctx, c.writer.Call("write", chunk)); err != nil {
		return 0, normalizeStreamError(err)
	}
	return len(p), nil
}

func (c *Conn) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	c.mu.Unlock()

	_, errCancel := awaitPromise(context.Background(), c.reader.Call("cancel"))
	_, errClose := awaitPromise(context.Background(), c.writer.Call("close"))
	return errors.Join(normalizeStreamError(errCancel), normalizeStreamError(errClose))
}

func (c *Conn) LocalAddr() net.Addr { return c.addr }

func (c *Conn) RemoteAddr() net.Addr { return c.addr }

func (c *Conn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.readDeadline = t
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writeDeadline = t
	return nil
}

func (c *Conn) CloseRead() error {
	_, err := awaitPromise(context.Background(), c.reader.Call("cancel"))
	return err
}

func (c *Conn) CloseWrite() error {
	_, err := awaitPromise(context.Background(), c.writer.Call("close"))
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
	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			thenFunc.Release()
			catchFunc.Release()
		})
	}

	thenFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer cleanup()
		value := js.Undefined()
		if len(args) > 0 {
			value = args[0]
		}
		resultc <- promiseResult{value: value}
		return nil
	})
	catchFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer cleanup()
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
