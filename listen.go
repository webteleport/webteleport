package webteleport

import (
	"context"
	"errors"
	"net"
	"net/url"

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
func Listen(ctx context.Context, u string) (net.Listener, error) {
	URL, _ := url.Parse(u)
	endpoints := endpoint.Resolve(URL)
	if len(endpoints) == 0 {
		return nil, errors.New("service discovery failed: no webteleport endpoints found in Alt-Svc records / headers")
	}
	ep := endpoints[0]
	switch ep.Protocol {
	case "websocket":
		return websocket.Listen(ctx, ep.Addr)
	case "webtransport":
		return webtransport.Listen(ctx, ep.Addr)
	default:
		return webtransport.Listen(ctx, ep.Addr)
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
