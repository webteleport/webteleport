package client

import (
	"log"
	"net/http"

	"github.com/btwiuse/ufo"
)

func Run([]string) error {
	ln, err := ufo.Listen("https://ufo.k0s.io")
	if err != nil {
		return err
	}
	log.Println("listening on", ln.URL())
	return http.Serve(ln, http.FileServer(http.Dir(".")))
}
