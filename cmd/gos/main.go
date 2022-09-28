package main

import (
	"log"
	"os"

	"github.com/webteleport/webteleport/gos"
)

func main() {
	log.Fatalln(gos.Run(os.Args[1:]))
}
