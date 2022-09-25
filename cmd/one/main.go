package main

import (
	"log"
	"os"

	"github.com/btwiuse/ufo/one"
)

func main() {
	log.Fatalln(one.Run(os.Args[1:]))
}
