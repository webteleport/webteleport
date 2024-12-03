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
	IsRoot(string) bool
	Upgrade() (*Edge, error)
}
