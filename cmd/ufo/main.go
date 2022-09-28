package main

import (
	"log"
	"os"

	"github.com/btwiuse/multicall"
	"github.com/webteleport/webteleport/echo"
	"github.com/webteleport/webteleport/gos"
	"github.com/webteleport/webteleport/hdr"
	"github.com/webteleport/webteleport/hello"
	"github.com/webteleport/webteleport/nc"
	"github.com/webteleport/webteleport/rp"
	"github.com/webteleport/webteleport/server"
	"github.com/webteleport/webteleport/sse"
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
