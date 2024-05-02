package spec

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

// Upgrade incoming edge requests via HTTP
type HTTPUpgrader interface {
	Upgrader
	http.Handler
}

// Upgrade incoming edge requests
type Upgrader interface {
	Root() string
	Upgrade() (*Edge, error)
}

// Subscribe to incoming requests
type Subscriber interface {
	Subscribe(upgrader Upgrader)
}
