package edge

import (
	"net/http"
	"net/url"

	"github.com/webteleport/webteleport/tunnel"
)

// Edge is a transport agnostic edge object
type Edge struct {
	Session tunnel.Session
	Stream  tunnel.Stream
	Path    string
	Values  url.Values
	Header  http.Header
	RealIP  string
}
