package webtransport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/quic-go/webtransport-go"
	"github.com/webteleport/utils"
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

var _ net.Listener = (*WebtransportListener)(nil)

func Listen(ctx context.Context, ep string, relayURL *url.URL) (*WebtransportListener, error) {
	// localhost:3000 will be parsed by net/url as URL{Scheme: localhost, Port: 3000}
	// hence the hack
	if !strings.Contains(ep, "://") {
		ep = "http://" + ep
	}
	session, err := DialWebtransport(ctx, ep, relayURL, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	stm0, err := session.AcceptStream(ctx)
	if err != nil {
		return nil, fmt.Errorf("stm0: %w", err)
	}
	errchan := make(chan string)
	hostchan := make(chan string)
	// go io.Copy(os.Stdout, stm0)
	go func() {
		scanner := bufio.NewScanner(stm0)
		for scanner.Scan() {
			line := scanner.Text()
			// ignore server pings
			if line == "PING" {
				continue
			}
			if strings.HasPrefix(line, "HOST ") {
				hostchan <- strings.TrimPrefix(line, "HOST ")
				continue
			}
			if strings.HasPrefix(line, "ERR ") {
				errchan <- strings.TrimPrefix(line, "ERR ")
				continue
			}
			slog.Warn(fmt.Sprintf("stm0: unknown command: %s", line))
		}
	}()
	// go io.Copy(stm0, os.Stdin)
	// Start a goroutine to gracefully handle close signal
	// TODO: cancel when server exits
	go func() {
		signalChannel := make(chan os.Signal, 1)

		// Notify the channel for specific signals
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		// Wait for the signal
		<-signalChannel

		// Print "bye" when the program exits
		_, err := io.WriteString(stm0, "CLOSE\n")
		if err != nil {
			slog.Warn(fmt.Sprintf("close error: %v", err))
		}
		slog.Info("terminating in 1 second")
		time.Sleep(time.Second)

		os.Exit(0)
	}()

	ln := &WebtransportListener{
		session: session,
		stm0:    stm0,
		scheme:  relayURL.Scheme,
		port:    utils.ExtractURLPort(relayURL),
	}
	select {
	case emsg := <-errchan:
		return nil, fmt.Errorf("server: %s", emsg)
	case ln.host = <-hostchan:
		return ln, nil
	}
}

type Session interface {
	AcceptStream(context.Context) (Stream, error)
}

type Stream interface {
	io.Reader
	io.Writer
}

// WebtransportListener implements [net.Listener]
type WebtransportListener struct {
	session *webtransport.Session // Session // webtransportSession,
	stm0    webtransport.Stream   // Stream  // webtransport.Stream
	scheme  string
	host    string
	port    string
}

// calling Accept returns a new [net.Conn]
func (l *WebtransportListener) Accept() (net.Conn, error) {
	stream, err := l.session.AcceptStream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("accept: %w", err)
	}
	return NewAcceptedConn(stream, l.session), nil
}

func (l *WebtransportListener) Close() error {
	return l.session.CloseWithError(1337, "foobar")
}

// Addr returns Listener itself which is an implementor of [net.Addr]
func (l *WebtransportListener) Addr() net.Addr {
	return &WebtransportAddr{l}
}

type WebtransportAddr struct {
	*WebtransportListener
}

// Network returns the protocol scheme, either http or https
func (addr *WebtransportAddr) Network() string {
	return addr.WebtransportListener.scheme
}

// String returns the host(:port) address of Listener, forcing ASCII
func (addr *WebtransportAddr) String() string {
	return utils.ToIdna(addr.WebtransportListener.host) + addr.WebtransportListener.port
}

// String returns the host(:port) address of Listener, Unicode is kept inact
func Display(l *WebtransportListener) string {
	return l.host + l.port
}

// AsciiURL returns the public accessible address of the Listener
func AsciiURL(l *WebtransportListener) string {
	return l.Addr().Network() + "://" + l.Addr().String()
}

// HumanURL returns the human readable URL
func HumanURL(l *WebtransportListener) string {
	return l.Addr().Network() + "://" + Display(l)
}

// AutoURL returns a clickable url the URL
//
//	when link == text, it displays `link[link]`
//	when link != text, it displays `text ([link](link))`
func ClickableURL(l *WebtransportListener) string {
	disp, link := HumanURL(l), AsciiURL(l)
	if disp == link {
		return utils.MaybeHyperlink(link)
	}
	return fmt.Sprintf("%s ( %s )", disp, utils.MaybeHyperlink(link))
}
