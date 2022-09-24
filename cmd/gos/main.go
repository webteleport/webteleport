package main

import (
	"log"
	"os"

	"github.com/btwiuse/quichost/echo"
)

func main() {
	log.Fatalln(echo.Run(os.Args[1:]))
}
