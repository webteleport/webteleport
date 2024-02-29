package webteleport

import (
	"context"
	"crypto/tls"
	"strings"

	"github.com/webtransport/quic"
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

func Dial(ctx context.Context, addr string) (*quic.Conn, error) {
	quicConf := &quic.Config{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS13,
			NextProtos:         []string{"hq-interop"},
		},
	}
	l, err := quic.Listen("udp", ":0", quicConf)
	if err != nil {
		return nil, err
	}
	session, err := l.Dial(ctx, "udp", addr, quicConf)
	return session, err
}
