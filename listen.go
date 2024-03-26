package webteleport

import (
	"context"
	"net"
	"net/url"

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
	ep := endpoint.Resolve(relayURL)[0]

	switch ep.Protocol {
	case "webtransport":
		addr, err := webtransport.Merge(ep.Addr, relayURL)
		if err != nil {
			return nil, err
		}
		return webtransport.Listen(ctx, addr)
	default:
		addr, err := websocket.Merge(ep.Addr, relayURL)
		if err != nil {
			return nil, err
		}
		return websocket.Listen(ctx, addr)
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
