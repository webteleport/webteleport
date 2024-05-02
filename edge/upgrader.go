package edge

import (
	"net/http"
)

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
