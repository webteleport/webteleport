package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/webteleport/webteleport/transport/common"
)

func Dial(ctx context.Context, addr string) (*TcpSession, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s (TCP): %w", addr, err)
	}

	session, err := common.YamuxServer(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating yamux server session: %w", err)
	}

	return &TcpSession{session}, nil
}
