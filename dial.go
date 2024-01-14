package webteleport

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

func altsvcLines(txts []string) []string {
	const prefix = "Alt-Svc: "

	altsvcs := []string{}

	for _, txt := range txts {
		// Case insensitive prefix match. See Issue 22736.
		if len(txt) < len(prefix) || !strings.EqualFold(txt[:len(prefix)], prefix) {
			continue
		}
		altsvcs = append(altsvcs, txt[len(prefix):])
	}

	return altsvcs
}

func Dial(ctx context.Context, addr string) (quic.Connection, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		// NextProtos:         []string{"spdy/3", "h2", "hq-29"},
		// ClientSessionCache: tls.NewLRUClientSessionCache(1),
	}
	quicConf := &quic.Config{
		MaxIdleTimeout: time.Minute * 10080,
		// KeepAlive:             true,
		MaxIncomingStreams:    1000000,
		MaxIncomingUniStreams: 1000000,
		// TokenStore:            quicGo.NewLRUTokenStore(1, 1),
	}
	session, err := quic.DialAddr(ctx, addr, tlsConf, quicConf)
	return session, err
}
