package main

import (
	"log"
	"os"

	"github.com/webteleport/ufo/sse"
)

func main() {
	log.Fatalln(sse.Run(os.Args[1:]))
}
