package gos

import (
	"log"
	"net/http"

	"github.com/btwiuse/ufo"
	"k0s.io/pkg/middleware"
)

func Run([]string) error {
	ln, err := ufo.Listen("https://ufo.k0s.io")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("listening on", ln.URL())
	fs := middleware.LoggingMiddleware(http.FileServer(http.Dir(".")))
	return http.Serve(ln, fs)
}
