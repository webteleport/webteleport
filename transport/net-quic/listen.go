package quic

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/webteleport/utils"
	"github.com/webteleport/webteleport/transport/common"
)

func Listen(ctx context.Context, addr string) (*common.Listener, error) {
	u, err := url.Parse(utils.AsURL(addr))
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	session, err := Dial(ctx, u.Host)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	stm0, err := session.Open(context.Background())
	if err != nil {
		return nil, fmt.Errorf("stm0: %w", err)
	}
	io.WriteString(stm0, fmt.Sprintf("%s\n", u.RequestURI()))

	hostport, err := common.ReadHandshake(stm0)
	if err != nil {
		return nil, err
	}

	return &common.Listener{
		Session: session,
		Scheme:  u.Scheme,
		Address: hostport,
	}, nil

}
