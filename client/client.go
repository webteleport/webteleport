package main

import (
	"fmt"
	"log"
	"net"

	"github.com/btwiuse/quichost"
)

func main() {
	ln, err := quichost.Listen("https://quichost.k0s.io")
	if err != nil {
		log.Fatalln(err)
	}
	addr := ln.Addr()
	network := addr.Network()
	hostport := addr.String()
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("/%s/%s/%s\n", network, host, port)
}
