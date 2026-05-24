package webteleport

import (
	"context"
	"log/slog"
	"net"
	"net/url"

	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/endpoint"
	"github.com/webteleport/webteleport/transport/websocket"
	"github.com/webteleport/webteleport/transport/webtransport"
	"github.com/webteleport/webteleport/tunnel"
)

// Listen calls [Dial] to create a [Listener], which is essentially a wrapper struct
// around a webtransport session, which in turn is able to spawn arbitrary number of streams
// that implements [net.Conn]
//
// It is modelled after [net.Listen], however it doesn't require the caller to be able to
// bind to a local port.
//
// The returned Listener can be imagined to be bound to a remote [net.Addr], which can be obtained
// using the [Listener.Addr] method
func Listen(ctx context.Context, relayAddr string) (net.Listener, error) {
	relayURL, err := url.Parse(utils.AsURL(relayAddr))
	if err != nil {
		return nil, err
	}

	// try to find ALT_SVC records in ENV/DNS/HEAD, see endpoint.Resolve
	var (
		dialAddr string
		tr       tunnel.Transport
	)
	for _, ep := range endpoint.Resolve(relayURL) {
		switch ep.Protocol {
		case "webtransport":
			dialAddr, err = webtransport.DialAddr(ep.Addr, relayURL)
			if err != nil {
				slog.Warn("webtransport dial error", "addr", ep.Addr, "error", err)
				continue
			}
			tr = &webtransport.Transport{}
		case "websocket":
			dialAddr, err = websocket.DialAddr(ep.Addr, relayURL)
			if err != nil {
				return nil, err
			}
			tr = &websocket.Transport{}
		}
		break
	}
	return tr.Listen(ctx, dialAddr)
}
