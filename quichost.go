package quichost

import (
	"context"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/marten-seemann/webtransport-go"
)

func Listen(u string) (net.Listener, error) {
	up, err := url.Parse(u)
	if err != nil {
		return nil, nil
	}
	log.Println(up.Scheme)
	log.Println(up.Host)
	log.Println("dialing", u)
	ctx, _ := context.WithTimeout(context.TODO(), 3*time.Second)
	session, err := dial(ctx, up)
	return &listener{session}, nil
}

type listener struct {
	session *webtransport.Session
}

func (l *listener) Accept() (net.Conn, error) {
	return nil, nil
}

func (l *listener) Close() error {
	return nil
}

func (l *listener) Addr() net.Addr {
	return l
}

func (l *listener) Network() string {
	return "https"
}

func (l *listener) String() string {
	return "TODO.quichost.k0s.io:https"
}
