package quic

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/webtransport/quic"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

func Dial(ctx context.Context, addr string) (*QuicSession, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	quicConf := &quic.Config{
		TLSConfig:            tlsConfig,
		MaxBidiRemoteStreams: MaxIncomingStreams,
	}
	l, err := quic.Listen("udp", ":0", quicConf)
	if err != nil {
		return nil, err
	}
	session, err := l.Dial(ctx, "udp", addr, quicConf)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (go-quic): %w", addr, err)
	}
	return &QuicSession{session}, nil
}
