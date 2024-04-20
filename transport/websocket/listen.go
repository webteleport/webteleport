package websocket

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"strings"

	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/transport"
)

var _ net.Listener = (*WebsocketListener)(nil)

func Listen(ctx context.Context, addr string) (*WebsocketListener, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	session, err := DialWebsocket(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	stm0, err := session.AcceptStream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("stm0: %w", err)
	}
	errchan := make(chan string)
	hostchan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stm0)
		for scanner.Scan() {
			line := scanner.Text()
			// ignore server pings
			if line == "" || line == "PING" {
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

	ln := &WebsocketListener{
		session: session,
		scheme:  u.Scheme,
	}
	select {
	case emsg := <-errchan:
		return nil, fmt.Errorf("server: %s", emsg)
	case hostport := <-hostchan:
		// TODO handle host + port + path
		if host, port, err := net.SplitHostPort(hostport); err == nil {
			ln.host = host
			ln.port = port
		} else {
			ln.host = hostport
		}
	}
	return ln, nil
}

// WebsocketListener implements [net.Listener]
type WebsocketListener struct {
	session transport.Session
	scheme  string
	host    string
	port    string
}

// calling Accept returns a new [net.Conn]
func (l *WebsocketListener) Accept() (net.Conn, error) {
	streamConn, err := l.session.AcceptStream(context.Background())
	if err != nil {
		return nil, fmt.Errorf("accept: %w", err)
	}
	return streamConn, nil
}

func (l *WebsocketListener) Close() error {
	return l.session.Close()
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
