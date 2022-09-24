package client

import (
	"log"
	"net/http"

	"github.com/btwiuse/skynet"
)

func Run([]string) error {
	ln, err := skynet.Listen("https://skynet.k0s.io")
	if err != nil {
		return err
	}
	log.Println("listening on", ln.URL())
	return http.Serve(ln, http.FileServer(http.Dir(".")))
}
