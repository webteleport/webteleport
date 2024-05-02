package edge

import (
	"net/http"
	"net/url"

	"github.com/webteleport/transport"
)

// Edge is a transport agnostic edge object
type Edge struct {
	Session transport.Session
	Stream  transport.Stream
	Path    string
	Values  url.Values
	Header  http.Header
	RealIP  string
}
