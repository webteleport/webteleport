package main

import (
	"context"
	"log"
	"net/http"

	"github.com/webteleport/webteleport"
)

func main() {
	ln, err := webteleport.Listen(context.TODO(), "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on", ln.Addr())
	defer ln.Close()
	http.Serve(ln, http.HandlerFunc(teapot))

}

func teapot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(418), 418)
}
