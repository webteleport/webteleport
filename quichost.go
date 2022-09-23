package quichost

import (
	"log"
	"net"
	"net/url"
)

func Listen(u string) (net.Listener, error) {
	up, err := url.Parse(u)
	if err != nil {
		return nil, nil
	}
	log.Println(up.Scheme)
	log.Println(up.Host)
	return &listener{}, nil
}

type listener struct{}

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
