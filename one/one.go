// this tool is used to reproduce server crash when MaxIncomingStreams is too small
//
// don't use it in prod
//
// should make ufo.MaxIncomingStreams a const once a remedy is found

package one

import (
	"io"
	"log"
	"net/http"

	"github.com/btwiuse/ufo"
)

func Arg0(args []string, fallback string) string {
	if len(args) > 0 {
		return args[0]
	}
	return fallback
}

func Run(args []string) error {
	ufo.MaxIncomingStreams = 1 + 1 // one for stm0, another req0
	ln, err := ufo.Listen(Arg0(args, "https://ufo.k0s.io"))
	if err != nil {
		return err
	}
	log.Println("listening on", ln.URL())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, UFO! (this should not appear after a refresh)\n")
	})
	return http.Serve(ln, nil)
}
