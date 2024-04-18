//go:build js

package websocket

import (
	"net/http"

	"nhooyr.io/websocket"
)

func dialOptions(http.Header) *websocket.DialOptions {
	return &websocket.DialOptions{}
}
