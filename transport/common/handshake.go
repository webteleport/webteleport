package common

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
)

var HandshakeTimeout = 10 * time.Second

// ReadHandshake reads the HOST/ERR handshake protocol from stm0.
// Empty lines and "PING" are ignored. "HOST <addr>" returns the address.
// "ERR <msg>" returns an error. Unknown lines are logged as a warning.
// Returns an error if no valid handshake is received within HandshakeTimeout.
func ReadHandshake(stm0 io.Reader) (string, error) {
	errchan := make(chan string)
	hostchan := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stm0)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || line == "PING" {
				continue
			}
			if v, ok := strings.CutPrefix(line, "HOST "); ok {
				hostchan <- v
				return
			}
			if v, ok := strings.CutPrefix(line, "ERR "); ok {
				errchan <- v
				return
			}
			slog.Warn("stm0: unknown command", "command", line)
		}
	}()
	select {
	case emsg := <-errchan:
		return "", fmt.Errorf("server: %s", emsg)
	case hostport := <-hostchan:
		return hostport, nil
	case <-time.After(HandshakeTimeout):
		return "", fmt.Errorf("handshake timeout after %v", HandshakeTimeout)
	}
}
