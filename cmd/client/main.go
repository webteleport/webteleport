// this is an example consumer of the skynet package
// it listens on a random skynet with registered handlers for webtransport && websocket connections
// currently websocket works fine
// while webtransport is broken because reverseproxy doesn't support it yet

package main

import (
	"log"
	"os"

	"github.com/btwiuse/skynet/client"
)

// curl3 https://7.skynet.k0s.io:300 --http3 -H "Host: 7.skynet.k0s.io"
func main() {
	log.Fatalln(client.Run(os.Args[1:]))
}
