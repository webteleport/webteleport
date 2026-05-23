package webtransport

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

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
	errchan := make(chan string)
	hostchan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stm0)
		for scanner.Scan() {
			line := scanner.Text()
			kind, payload, ok := parseControlLine(line)
			if !ok {
				if line == "" || line == "PING" {
					continue
				}
				slog.Warn(fmt.Sprintf("stm0: unknown command: %s", line))
				continue
			}
			if kind == "HOST" {
				hostchan <- payload
				continue
			}
			if kind == "ERR" {
				errchan <- payload
				continue
			}
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

func parseControlLine(line string) (kind, payload string, ok bool) {
	if line == "" || line == "PING" {
		return "", "", false
	}
	if strings.HasPrefix(line, "HOST ") {
		return "HOST", strings.TrimPrefix(line, "HOST "), true
	}
	if strings.HasPrefix(line, "ERR ") {
		return "ERR", strings.TrimPrefix(line, "ERR "), true
	}
	return "", "", false
}
