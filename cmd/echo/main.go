// this is an example consumer of the quichost package
// it listens on a random quichost with registered handlers for webtransport && websocket connections
// currently websocket works fine
// while webtransport is broken because reverseproxy doesn't support it yet

package main

import (
	"log"
	"os"

	"github.com/btwiuse/quichost/echo"
)

// curl3 https://7.quichost.k0s.io:300 --http3 -H "Host: 7.quichost.k0s.io"
func main() {
	log.Fatalln(echo.Run(os.Args[1:]))
}
