package main

import (
	"log"
	"os"

	"github.com/btwiuse/ufo/gos"
)

func main() {
	log.Fatalln(gos.Run(os.Args[1:]))
}
