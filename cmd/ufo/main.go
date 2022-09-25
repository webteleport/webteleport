package main

import (
	"log"
	"os"

	"github.com/btwiuse/multicall"
	"github.com/btwiuse/ufo/client"
	"github.com/btwiuse/ufo/echo"
	"github.com/btwiuse/ufo/gos"
	"github.com/btwiuse/ufo/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	err := Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

var cmdRun multicall.RunnerFuncMap = map[string]multicall.RunnerFunc{
	"client": client.Run,
	"echo":   echo.Run,
	"gos":    gos.Run,
	"server": server.Run,
}

func Run(args []string) error {
	return cmdRun.Run(os.Args[1:])
}
