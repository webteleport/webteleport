// this is an example consumer of the ufo package
// it listens on a random ufo with registered handlers for webtransport && websocket connections
// currently websocket works fine
// while webtransport is broken because reverseproxy doesn't support it yet

package main

import (
	"log"
	"os"

	"github.com/btwiuse/ufo/client"
)

// curl3 https://7.ufo.k0s.io:300 --http3 -H "Host: 7.ufo.k0s.io"
func main() {
	log.Fatalln(client.Run(os.Args[1:]))
}
