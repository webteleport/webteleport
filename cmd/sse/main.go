package main

import (
	"log"
	"os"

	"github.com/btwiuse/ufo/sse"
)

func main() {
	log.Fatalln(sse.Run(os.Args[1:]))
}
