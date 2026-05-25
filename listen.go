package webteleport

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/url"

	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/endpoint"
	quicgo "github.com/webteleport/webteleport/transport/quic-go"
	netquic "github.com/webteleport/webteleport/transport/net-quic"
	"github.com/webteleport/webteleport/transport/tcp"
	"github.com/webteleport/webteleport/transport/websocket"
	"github.com/webteleport/webteleport/transport/webtransport"
	"github.com/webteleport/webteleport/tunnel"
)

type candidate struct {
	dialAddr string
	tr       tunnel.Transport
}

func fromEndpoints(eps []endpoint.Endpoint, relayURL *url.URL) []candidate {
	var cs []candidate
	for _, ep := range eps {
		var (
			dialAddr string
			tr       tunnel.Transport
			err      error
		)
		switch ep.Protocol {
		case "webtransport":
			dialAddr, err = webtransport.DialAddr(ep.Addr, relayURL)
			tr = &webtransport.Transport{}
		case "net-quic":
			dialAddr = ep.Addr
			tr = &netquic.Transport{}
		case "quic", "quic-go":
			dialAddr = ep.Addr
			tr = &quicgo.Transport{}
		case "tcp":
			dialAddr = ep.Addr
			tr = &tcp.Transport{}
		case "websocket":
			dialAddr, err = websocket.DialAddr(ep.Addr, relayURL)
			tr = &websocket.Transport{}
		}
		if err != nil {
			slog.Warn("dial error", "protocol", ep.Protocol, "addr", ep.Addr, "error", err)
			continue
		}
		cs = append(cs, candidate{dialAddr: dialAddr, tr: tr})
	}
	return cs
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
	for _, c := range fromEndpoints(endpoint.Resolve(ctx, relayURL), relayURL) {
		l, err := c.tr.Listen(ctx, c.dialAddr)
		if err != nil {
			lastErr = err
			slog.Warn("listen error", "dialAddr", c.dialAddr, "error", err)
			continue
		}
		return l, nil
	}
	return nil, lastErr
}
