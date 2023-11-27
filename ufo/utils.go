package ufo

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/webteleport/webteleport"
)

// listen with a timeout
func listenWithTimeout(timeout time.Duration, addr string) (*webteleport.Listener, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return webteleport.Listen(ctx, addr)
}

// createURLWithQueryParams creates a URL with query parameters
func createURLWithQueryParams(stationURL string) (*url.URL, error) {
	// parse the station URL
	u, err := url.Parse(stationURL)
	if err != nil {
		return nil, err
	}

	// attach extra info to the query string
	q := u.Query()
	q.Add("client", "ufo")
	for _, arg := range os.Args {
		q.Add("args", arg)
	}
	u.RawQuery = q.Encode()

	return u, nil
}
