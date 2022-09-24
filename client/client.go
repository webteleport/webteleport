package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/btwiuse/skynet"
)

func Run([]string) error {
	ln, err := skynet.Listen("https://skynet.k0s.io")
	if err != nil {
		return err
	}
	addr := ln.Addr()
	location := fmt.Sprintf("%s://%s", addr.Network(), addr.String())
	log.Println("listening on", location)
	return http.Serve(ln, http.FileServer(http.Dir(".")))
}
