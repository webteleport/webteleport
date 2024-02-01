package websocket

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/webteleport/utils"
)

var _ net.Listener = (*WebsocketListener)(nil)

func Listen(ctx context.Context, ep string, relayURL *url.URL) (*WebsocketListener, error) {
	session, err := DialWebsocket(ctx, ep, relayURL, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	stm0, err := session.AcceptStream()
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

	ln := &WebsocketListener{
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

// WebsocketListener implements [net.Listener]
type WebsocketListener struct {
	session *yamux.Session
	stm0    *yamux.Stream
	scheme  string
	host    string
	port    string
}

// calling Accept returns a new [net.Conn]
func (l *WebsocketListener) Accept() (net.Conn, error) {
	stream, err := l.session.AcceptStream()
	if err != nil {
		return nil, fmt.Errorf("accept: %w", err)
	}
	return NewAcceptedConn(stream), nil
}

func (l *WebsocketListener) Close() error {
	l.session.Close()
	return http.ErrServerClosed
}

// Addr returns Listener itself which is an implementor of [net.Addr]
func (l *WebsocketListener) Addr() net.Addr {
	return &WebsocketAddr{l}
}

type WebsocketAddr struct {
	*WebsocketListener
}

// Network returns the protocol scheme, either http or https
func (addr *WebsocketAddr) Network() string {
	return addr.WebsocketListener.scheme
}

// String returns the host(:port) address of Listener, forcing ASCII
func (addr *WebsocketAddr) String() string {
	return utils.ToIdna(addr.WebsocketListener.host) + addr.WebsocketListener.port
}

// Display returns the host(:port) address of Listener, Unicode is kept inact
func Display(l *WebsocketListener) string {
	return l.host + l.port
}

// AsciiURL returns the public accessible address of the Listener
func AsciiURL(l *WebsocketListener) string {
	return l.Addr().Network() + "://" + l.Addr().String()
}

// HumanURL returns the human readable URL
func HumanURL(l *WebsocketListener) string {
	return l.Addr().Network() + "://" + Display(l)
}

// ClickableURL returns a clickable url the URL
//
//	when link == text, it displays `link[link]`
//	when link != text, it displays `text ([link](link))`
func ClickableURL(l *WebsocketListener) string {
	disp, link := HumanURL(l), AsciiURL(l)
	if disp == link {
		return utils.MaybeHyperlink(link)
	}
	return fmt.Sprintf("%s ( %s )", disp, utils.MaybeHyperlink(link))
}
