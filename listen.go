package webteleport

import (
	"context"
	"net"
	"net/url"
	"os"

	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/endpoint"
	"github.com/webteleport/webteleport/websocket"
	"github.com/webteleport/webteleport/webtransport"
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
	endpoints := endpoint.Resolve(relayURL)

	// use websocket transport when no ALT_SVC records found, or env WEBSOCKET is set
	if len(endpoints) == 0 || os.Getenv("WEBSOCKET") != "" {
		return websocket.Listen(ctx, relayURL.Hostname(), relayURL)
	}

	// otherwise use whatever protocol specified in the first endpoint
	ep := endpoints[0]
	switch {
	case ep.Protocol == "websocket":
		return websocket.Listen(ctx, ep.Addr, relayURL)
	default:
		return webtransport.Listen(ctx, ep.Addr, relayURL)
	}
}

// String returns the host(:port) address of Listener, Unicode is kept inact
func Display(i interface{}) string {
	switch l := i.(type) {
	case *webtransport.WebtransportListener:
		return webtransport.Display(l)
	case *websocket.WebsocketListener:
		return websocket.Display(l)
	}
	return "<unknown>"
}

// AsciiURL returns the public accessible address of the Listener
func AsciiURL(i interface{}) string {
	switch l := i.(type) {
	case *webtransport.WebtransportListener:
		return webtransport.AsciiURL(l)
	case *websocket.WebsocketListener:
		return websocket.AsciiURL(l)
	}
	return "<unknown>"
}

// HumanURL returns the human readable URL
func HumanURL(i interface{}) string {
	switch l := i.(type) {
	case *webtransport.WebtransportListener:
		return webtransport.HumanURL(l)
	case *websocket.WebsocketListener:
		return websocket.HumanURL(l)
	}
	return "<unknown>"
}

// AutoURL returns a clickable url the URL
//
//	when link == text, it displays `link[link]`
//	when link != text, it displays `text ([link](link))`
func ClickableURL(i interface{}) string {
	switch l := i.(type) {
	case *webtransport.WebtransportListener:
		return webtransport.ClickableURL(l)
	case *websocket.WebsocketListener:
		return websocket.ClickableURL(l)
	}
	return "<unknown>"
}
