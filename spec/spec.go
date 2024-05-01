package spec

import (
	"net/http"
	"net/url"

	"github.com/webteleport/transport"
)

// Request is a transport agnostic request object
type Request struct {
	Session transport.Session
	Stream  transport.Stream
	Path    string
	Values  url.Values
	Header  http.Header
	RealIP  string
}

// Upgrade incoming requests via HTTP
type HTTPUpgrader interface {
	Upgrader
	http.Handler
}

// Upgrade incoming requests
type Upgrader interface {
	Root() string
	Upgrade() (*Request, error)
}

// Subscribe to incoming requests
type Subscriber interface {
	Subscribe(upgrader Upgrader)
}
