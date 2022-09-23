package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/btwiuse/quichost"
)

func Run([]string) error {
	ln, err := quichost.Listen("https://quichost.k0s.io")
	if err != nil {
		return err
	}
	addr := ln.Addr()
	location := fmt.Sprintf("%s://%s", addr.Network(), addr.String())
	log.Println("listening on", location)
	return http.Serve(ln, http.FileServer(http.Dir(".")))
}
