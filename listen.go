package webteleport

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/url"
	"runtime"

	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/endpoint"
	quicgo "github.com/webteleport/webteleport/transport/quic-go"
	netquic "github.com/webteleport/webteleport/transport/net-quic"
	"github.com/webteleport/webteleport/transport/tcp"
	"github.com/webteleport/webteleport/transport/websocket"
	"github.com/webteleport/webteleport/transport/webtransport"
	"github.com/webteleport/webteleport/tunnel"
)

func fromEndpoints(eps []endpoint.Endpoint, relayURL *url.URL) []tunnel.Transport {
	var ts []tunnel.Transport
	for _, ep := range eps {
		var (
			dialAddr string
			tr       tunnel.Transport
			err      error
		)
		switch ep.Protocol {
		case "webtransport":
			dialAddr, err = webtransport.ResolveAddr(ep.Addr, relayURL)
			tr = &webtransport.Transport{DialAddr: dialAddr}
		case "net-quic":
			if runtime.GOOS == "js" {
				err = errors.New("net-quic unsupported on js/wasm")
				break
			}
			dialAddr, err = netquic.ResolveAddr(ep.Addr, relayURL)
			tr = &netquic.Transport{DialAddr: dialAddr}
		case "quic", "quic-go":
			if runtime.GOOS == "js" {
				err = errors.New("quic unsupported on js/wasm")
				break
			}
			dialAddr, err = quicgo.ResolveAddr(ep.Addr, relayURL)
			tr = &quicgo.Transport{DialAddr: dialAddr}
		case "tcp":
			if runtime.GOOS == "js" {
				err = errors.New("tcp unsupported on js/wasm")
				break
			}
			dialAddr, err = tcp.ResolveAddr(ep.Addr, relayURL)
			tr = &tcp.Transport{DialAddr: dialAddr}
		case "websocket":
			dialAddr, err = websocket.ResolveAddr(ep.Addr, relayURL)
			tr = &websocket.Transport{DialAddr: dialAddr}
		}
		if err != nil {
			slog.Warn("skip transport", "protocol", ep.Protocol, "addr", ep.Addr, "error", err)
			continue
		}
		ts = append(ts, tr)
	}
	return ts
}

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

	var lastErr error = errors.New("no endpoints available to attempt connection")
	for _, tr := range fromEndpoints(endpoint.Resolve(ctx, relayURL), relayURL) {
		l, err := tr.Listen(ctx, "")
		if err != nil {
			lastErr = err
			slog.Warn("listen error", "error", err)
			continue
		}
		return l, nil
	}
	return nil, lastErr
}
