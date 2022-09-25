package gos

import (
	"log"
	"net/http"

	"github.com/btwiuse/ufo"
	"k0s.io/pkg/middleware"
)

func Arg0(args []string, fallback string) string {
	if len(args) > 0 {
		return args[0]
	}
	return fallback
}

func Run(args []string) error {
	ln, err := ufo.Listen(Arg0(args, "https://ufo.k0s.io"))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("listening on", ln.URL())
	fs := middleware.LoggingMiddleware(http.FileServer(http.Dir(".")))
	return http.Serve(ln, fs)
}
