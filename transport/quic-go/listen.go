package quic

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"

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
	stm0, err := session.Accept(context.Background())
	if err != nil {
		return nil, fmt.Errorf("stm0: %w", err)
	}
	io.WriteString(stm0, fmt.Sprintf("%s\n", u.RequestURI()))
	errchan := make(chan string)
	hostchan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stm0)
		for scanner.Scan() {
			line := scanner.Text()
			// ignore server pings
			if line == "" || line == "PING" {
				continue
			}
			if strings.HasPrefix(line, "HOST ") {
				hostchan <- strings.TrimPrefix(line, "HOST ")
				continue
			}
			if strings.HasPrefix(line, "ERR ") {
				errchan <- strings.TrimPrefix(line, "ERR ")
				continue
			}
			slog.Warn(fmt.Sprintf("stm0: unknown command: %s", line))
		}
	}()

	ln := &common.Listener{
		Session: session,
		Scheme:  u.Scheme,
	}
	select {
	case emsg := <-errchan:
		return nil, fmt.Errorf("server: %s", emsg)
	case hostport := <-hostchan:
		ln.Address = hostport
		return ln, nil
	}
}
