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
			dialErr  error
		)
		switch ep.Protocol {
		case "webtransport":
			dialAddr, dialErr = webtransport.DialAddr(ep.Addr, relayURL)
			tr = &webtransport.Transport{DialAddr: dialAddr}
		case "net-quic":
			if runtime.GOOS == "js" {
				dialErr = errors.New("net-quic unsupported on js/wasm")
				break
			}
			tr = &netquic.Transport{DialAddr: ep.Addr}
		case "quic", "quic-go":
			if runtime.GOOS == "js" {
				dialErr = errors.New("quic unsupported on js/wasm")
				break
			}
			tr = &quicgo.Transport{DialAddr: ep.Addr}
		case "tcp":
			if runtime.GOOS == "js" {
				dialErr = errors.New("tcp unsupported on js/wasm")
				break
			}
			tr = &tcp.Transport{DialAddr: ep.Addr}
		case "websocket":
			dialAddr, dialErr = websocket.DialAddr(ep.Addr, relayURL)
			tr = &websocket.Transport{DialAddr: dialAddr}
		}
		if dialErr != nil {
			slog.Warn("dial error", "protocol", ep.Protocol, "addr", ep.Addr, "error", dialErr)
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
