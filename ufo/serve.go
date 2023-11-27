package ufo

import (
	"log"
	"net/http"
	"time"

	"github.com/webteleport/auth"
)

// DefaultTimeout is the default dialing timeout for the UFO server.
var DefaultTimeout = 10 * time.Second

// Serve starts a UFO server on the given station URL.
func Serve(stationURL string, handler http.Handler) error {
	// Create the URL with query parameters
	u, err := createURLWithQueryParams(stationURL)
	if err != nil {
		return err
	}

	// listen on the station URL with a timeout
	ln, err := listenWithTimeout(DefaultTimeout, u.String())
	if err != nil {
		return err
	}

	log.Println("🛸 listening on", ln.ClickableURL())

	// use the default serve mux if nil handler is provided
	if handler == nil {
		handler = http.DefaultServeMux
	}

	// configure password authentication middleware
	lm := &auth.LoginMiddleware{
		Password: u.Fragment,
	}

	// wrap the handler with password authentication
	if lm.IsPasswordRequired() {
		handler = lm.Wrap(handler)
		log.Println("🔒 secured by password authentication")
	} else {
		log.Println("🔓 publicly accessible without a password")
	}

	return http.Serve(ln, handler)
}
