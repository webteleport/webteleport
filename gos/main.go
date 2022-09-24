package gos

import (
	"log"
	"net/http"

	"github.com/btwiuse/quichost"
	"k0s.io/pkg/middleware"
)

func Run([]string) error {
	ln, err := quichost.Listen("https://quichost.k0s.io")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("listening on", ln.URL())
	fs := middleware.LoggingMiddleware(http.FileServer(http.Dir(".")))
	return http.Serve(ln, fs)
}
