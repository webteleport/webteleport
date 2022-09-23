package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/btwiuse/quichost"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	ln, err := quichost.Listen("https://quichost.k0s.io")
	if err != nil {
		log.Fatalln(err)
	}
	addr := ln.Addr()
	location := fmt.Sprintf("%s://%s", addr.Network(), addr.String())
	log.Println("listening on", location)
	http.Serve(ln, http.FileServer(http.Dir(".")))
}
