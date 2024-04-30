package quic

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/quic-go/quic-go"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

func Dial(ctx context.Context, addr string) (*QuicSession, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}
	quicConf := &quic.Config{
		EnableDatagrams:    true,
		MaxIncomingStreams: MaxIncomingStreams,
	}
	session, err := quic.DialAddr(ctx, addr, tlsConf, quicConf)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (quic-go): %w", addr, err)
	}
	return &QuicSession{session}, nil
}
