package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	for i := 0; ; i++ {
		x(i)
	}
}

func x(i int) {
	// log.Println("listening on", "tcp4", "0.0.0.0:0")
	ln, err := net.Listen("tcp4", "0.0.0.0:0")
	if err != nil {
		log.Fatalln(err)
	}
	network := ln.Addr().Network()
	host, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%d /%s/%s/%s\n", i, network, host, port)
}
