package hello

import (
	"context"
	"io"
	"log"
	"net/http"

	ufo "github.com/webteleport/webteleport"
)

func Run([]string) error {
	ln, err := ufo.Listen(context.Background(), "https://ufo.k0s.io")
	if err != nil {
		return err
	}
	log.Println("🛸 listening on", ln.ClickableURL())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, UFO!\n")
	})
	return http.Serve(ln, nil)
}
