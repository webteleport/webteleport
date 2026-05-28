package webtransport

import (
	"context"
	"fmt"
	"net/url"

	"github.com/webteleport/webteleport/transport/common"
)

func Listen(ctx context.Context, addr string) (*common.Listener, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	session, err := Dial(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	stm0, err := session.Accept(context.Background())
	if err != nil {
		return nil, fmt.Errorf("stm0: %w", err)
	}

	hostport, err := common.ReadHandshake(stm0)
	if err != nil {
		return nil, err
	}

	return &common.Listener{
		Session: session,
		Scheme:  u.Scheme,
		Address: hostport,
		Relay:   u,
	}, nil

}
