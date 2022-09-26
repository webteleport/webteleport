package hello

import (
	"io"
	"log"
	"net/http"

	"github.com/btwiuse/ufo"
)

func Run([]string) error {
	ln, err := ufo.Listen("https://ufo.k0s.io")
	if err != nil {
		return err
	}
	log.Println("🛸 listening on", ln.ClickableURL())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, UFO!\n")
	})
	return http.Serve(ln, nil)
}
