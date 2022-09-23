package main

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/btwiuse/h3/utils"
)

func main() {
	port := utils.EnvPort(":3000")
	log.Println("listening on TCP http://127.0.0.1" + port)
	ln, err := net.Listen("tcp4", port)
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, r.Host)
	})

	go func() {
		wts := webtransportServer(port)
		cert := utils.EnvCert("localhost.pem")
		key := utils.EnvKey("localhost-key.pem")
		log.Println("listening on UDP https://127.0.0.1" + port)
		log.Fatalln(wts.ListenAndServeTLS(cert, key))
	}()

	http.Serve(ln, http.DefaultServeMux)
}
