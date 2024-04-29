package quic

import (
	"context"
	"crypto/tls"

	"github.com/webtransport/quic"
)

func Dial(ctx context.Context, addr string) (*QuicSession, error) {
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
	return &QuicSession{session}, err
}
