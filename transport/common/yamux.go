package common

import (
	"io"
	"os"

	"github.com/hashicorp/yamux"
)

func YamuxConfig() *yamux.Config {
	c := yamux.DefaultConfig()
	if os.Getenv("YAMUX_LOG") == "" {
		c.LogOutput = io.Discard
	}
	c.EnableKeepAlive = false
	return c
}

// relay
func YamuxClient(conn io.ReadWriteCloser) (*yamux.Session, error) {
	config := YamuxConfig()
	session, err := yamux.Client(conn, config)
	return session, err
}

// client
func YamuxServer(conn io.ReadWriteCloser) (*yamux.Session, error) {
	config := YamuxConfig()
	session, err := yamux.Server(conn, config)
	return session, err
}
