package webteleport

import (
	// "bufio"
	"context"
	"errors"
	// "fmt"
	"io"
	// "log/slog"
	"net"
	"net/url"
	// "os"
	// "os/signal"
	// "strings"
	// "syscall"
	// "time"

	// "github.com/quic-go/webtransport-go"
	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/webtransport"
)

// Listen calls [Dial] to create a [Listener], which is essentially a wrapper struct
// around a webtransport session, which in turn is able to spawn arbitrary number of streams
// that implements [net.Conn]
//
// It is modelled after [net.Listen], however it doesn't require the caller to be able to
// bind to a local port.
//
// The returned Listener can be imagined to be bound to a remote [net.Addr], which can be obtained
// using the [Listener.Addr] method
func Listen(ctx context.Context, u string) (net.Listener, error) {
	URL, _:= url.Parse(u)
	endpoints := Resolve(URL)
	if len(endpoints) == 0 {
		return nil, errors.New("service discovery failed: no webteleport endpoints found in Alt-Svc records / headers")
	}
	alt := endpoints[0]
	addr := utils.Graft(URL.Host, alt)
	return webtransport.ListenWebtransport(ctx, addr)
}

type Session interface {
	AcceptStream(context.Context) (Stream, error)
}

type Stream interface {
	io.Reader
	io.Writer
}
