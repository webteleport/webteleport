package main

import (
	"log"
	"os"

	"github.com/btwiuse/skynet/gos"
)

func main() {
	log.Fatalln(gos.Run(os.Args[1:]))
}
