package quichost

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/marten-seemann/webtransport-go"
)

type Listener interface {
	Accept() (Conn, error)
	Close() error
	Addr() net.Addr
}

type Conn interface {
	io.Reader
	io.Writer
	io.Closer
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}

func Listen(u string) (Listener, error) {
	up, err := url.Parse(u)
	if err != nil {
		return nil, nil
	}
	log.Println(up.Scheme)
	log.Println(up.Host)
	log.Println("dialing", u)
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	session, err := dial(ctx, up)
	if err != nil {
		return nil, nil
	}
	stm0, err := session.OpenStream()
	if err != nil {
		return nil, nil
	}
	hostchan := make(chan string)
	ln := &listener{
		session: session,
		stm0:    stm0,
	}
	// go io.Copy(os.Stdout, stm0)
	go func() {
		scanner := bufio.NewScanner(stm0)
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
			if strings.HasPrefix(line, "HOST ") {
				hostchan <- strings.TrimPrefix(line, "HOST ")
			}
		}
	}()
	// go io.Copy(stm0, os.Stdin)
	ln.host = <-hostchan
	return ln, nil
}

type listener struct {
	session *webtransport.Session
	stm0    webtransport.Stream
	host    string
}

func (l *listener) Accept() (Conn, error) {
	return l.session.AcceptStream(context.Background())
}

func (l *listener) Close() error {
	return l.session.Close()
}

func (l *listener) Addr() net.Addr {
	return l
}

func (l *listener) Network() string {
	return "https"
}

func (l *listener) String() string {
	return l.host
}
