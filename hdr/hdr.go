package hdr

import (
	"io"
	"net/http"

	"github.com/btwiuse/pretty"
	"github.com/webteleport/ufo"
	"k0s.io/pkg/middleware"
)

func Arg0(args []string, fallback string) string {
	if len(args) > 0 {
		return args[0]
	}
	return fallback
}

func Run(args []string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, pretty.JSONStringLine(r.Header))
	})
	return ufo.Serve(Arg0(args, "https://ufo.k0s.io"), middleware.LoggingMiddleware(http.DefaultServeMux))
}
