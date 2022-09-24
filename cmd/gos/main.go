package main

import (
	"log"
	"os"

	"github.com/btwiuse/skynet/echo"
)

func main() {
	log.Fatalln(echo.Run(os.Args[1:]))
}
