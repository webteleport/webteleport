package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/webteleport/webteleport"
)

func relay() string {
	addr := os.Getenv("RELAY")
	if addr == "" {
		return "localhost:8080"
	}
	return addr
}

func main() {
	ln, err := webteleport.Listen(context.Background(), relay())
	if err != nil {
		log.Fatal(err)
	}
	endpoint := fmt.Sprintf("%s://%s", ln.Addr().Network(), ln.Addr().String())
	log.Println("Listening on", endpoint)
	defer ln.Close()
	http.Serve(ln, http.HandlerFunc(teapot))

}

func teapot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(418), 418)
}
