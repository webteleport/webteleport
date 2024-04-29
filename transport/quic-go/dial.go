package quic

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/quic-go/quic-go"
)

// 2^60 == 1152921504606846976
const MaxIncomingStreams int64 = 1 << 60

func Dial(ctx context.Context, addr string) (*QuicSession, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}
	quicConf := &quic.Config{
		MaxIdleTimeout:        time.Minute * 10080,
		MaxIncomingStreams:    1000000,
		MaxIncomingUniStreams: 1000000,
	}
	session, err := quic.DialAddr(ctx, addr, tlsConf, quicConf)
	return &QuicSession{session}, err
}
