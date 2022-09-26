package main

import (
	"log"
	"os"

	"github.com/btwiuse/multicall"
	"github.com/btwiuse/ufo/echo"
	"github.com/btwiuse/ufo/gos"
	"github.com/btwiuse/ufo/hdr"
	"github.com/btwiuse/ufo/hello"
	"github.com/btwiuse/ufo/nc"
	"github.com/btwiuse/ufo/rp"
	"github.com/btwiuse/ufo/server"
	"github.com/btwiuse/ufo/sse"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	err := Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}

var cmdRun multicall.RunnerFuncMap = map[string]multicall.RunnerFunc{
	"hello":  hello.Run,
	"echo":   echo.Run,
	"gos":    gos.Run,
	"hdr":    hdr.Run,
	"nc":     nc.Run,
	"server": server.Run,
	"sse":    sse.Run,
	"rp":     rp.Run,
}

func Run(args []string) error {
	return cmdRun.Run(os.Args[1:])
}
