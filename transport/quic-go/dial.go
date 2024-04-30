package quic

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/quic-go/quic-go"
)

// 2^60 == 1152921504606846976
var MaxIncomingStreams int64 = 1 << 60

var TLSConfig = &tls.Config{
	InsecureSkipVerify: true,
}

var QUICConfig = &quic.Config{
	EnableDatagrams:    true,
	MaxIncomingStreams: MaxIncomingStreams,
}

func Dial(ctx context.Context, addr string) (*QuicSession, error) {
	session, err := quic.DialAddr(ctx, addr, TLSConfig, QUICConfig)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (quic-go): %w", addr, err)
	}
	return &QuicSession{session}, nil
}
