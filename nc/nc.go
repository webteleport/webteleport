package nc

import (
	"io"
	"log"
	"net/http"
	"os"

	ufo "github.com/webteleport/webteleport"
	"k0s.io/pkg/middleware"
	"k0s.io/pkg/wrap"
)

func Arg0(args []string, fallback string) string {
	if len(args) > 0 {
		return args[0]
	}
	return fallback
}

func Run(args []string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := wrap.Hijack(w)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = io.Copy(conn, io.TeeReader(conn, os.Stderr))
		if err != nil {
			log.Println(err)
		}
	})
	return ufo.Serve(Arg0(args, "https://ufo.k0s.io"), middleware.LoggingMiddleware(http.DefaultServeMux))
}
