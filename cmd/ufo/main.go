package main

import (
	"log"
	"os"

	"github.com/btwiuse/multicall"
	"github.com/webteleport/ufo/echo"
	"github.com/webteleport/ufo/gos"
	"github.com/webteleport/ufo/hdr"
	"github.com/webteleport/ufo/hello"
	"github.com/webteleport/ufo/nc"
	"github.com/webteleport/ufo/rp"
	"github.com/webteleport/ufo/server"
	"github.com/webteleport/ufo/sse"
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
