package webteleport

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/marten-seemann/webtransport-go"
)

var _ net.Listener = (*listener)(nil)

var DefaultTimeout = 5 * time.Second

func Serve(u string, handler http.Handler) error {
	ctx, _ := context.WithTimeout(context.Background(), DefaultTimeout)
	ln, err := Listen(ctx, u)
	if err != nil {
		return err
	}
	log.Println("🛸 listening on", ln.ClickableURL())
	if handler == nil {
		handler = http.DefaultServeMux
	}
	return http.Serve(ln, handler)
}

func Listen(ctx context.Context, u string) (*listener, error) {
	// localhost:3000 will be parsed by net/url as URL{Scheme: localhost, Port: 3000}
	// hence the hack
	if !strings.Contains(u, "://") {
		u = "http://" + u
	}
	up, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	session, err := Dial(ctx, up, nil)
	if err != nil {
		return nil, err
	}
	stm0, err := session.AcceptStream(ctx)
	if err != nil {
		return nil, err
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
			log.Println("stm0: unknown command:", line)
		}
	}()
	// go io.Copy(stm0, os.Stdin)
	ln := &listener{
		session: session,
		stm0:    stm0,
		scheme:  up.Scheme,
		port:    ExtractURLPort(up),
	}
	select {
	case emsg := <-errchan:
		return nil, fmt.Errorf("server: %s", emsg)
	case ln.host = <-hostchan:
		return ln, nil
	}
}

type listener struct {
	session *webtransport.Session
	stm0    webtransport.Stream
	scheme  string
	host    string
	port    string
}

func (l *listener) Accept() (net.Conn, error) {
	stream, err := l.session.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}
	return &StreamConn{stream, l.session}, nil
}

func (l *listener) Close() error {
	return l.session.Close()
}

// Addr returns listener itself which is an implementor of net.Addr
func (l *listener) Addr() net.Addr {
	return l
}

// Network returns the protocol scheme, either http or https
func (l *listener) Network() string {
	return l.scheme
}

// String returns the host(:port) address of listener, forcing ASCII
func (l *listener) String() string {
	return ToIdna(l.host) + l.port
}

// String returns the host(:port) address of listener, Unicode is kept inact
func (l *listener) Display() string {
	return l.host + l.port
}

// AsciiURL returns the public accessible address of the listener
func (l *listener) AsciiURL() string {
	return l.Network() + "://" + l.String()
}

// HumanURL returns the human readable URL
func (l *listener) HumanURL() string {
	return l.Network() + "://" + l.Display()
}

// AutoURL returns a clickable url the URL
//   when link == text, it displays `link[link]`
//   when link != text, it displays `text ([link](link))`
func (l *listener) ClickableURL() string {
	disp, link := l.HumanURL(), l.AsciiURL()
	if disp == link {
		return MaybeHyperlink(link)
	}
	return fmt.Sprintf("%s ( %s )", disp, MaybeHyperlink(link))
}
