package common

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func ReadLine(conn io.Reader) (string, error) {
	// do multiple read to get the first line
	b := make([]byte, 1)
	var buf bytes.Buffer
	for {
		_, err := conn.Read(b)
		if err != nil {
			return "", fmt.Errorf("read line error: %w", err)
		}
		if b[0] == '\n' {
			break
		}
		buf.Write(b)
	}
	return strings.TrimSpace(buf.String()), nil
}
