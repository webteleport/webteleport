package ufo

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/webteleport/webteleport"
)

// listen with a timeout
func listenWithTimeout(addr string, timeout time.Duration) (*webteleport.Listener, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

// logServerStatus logs the status of the server.
func logServerStatus(ln *webteleport.Listener, u *url.URL) {
	log.Println("🛸 listening on", ln.ClickableURL())

	if u.Fragment == "" {
		log.Println("🔓 publicly accessible without a password")
	} else {
		log.Println("🔒 secured by password authentication")
	}
}

// parseQuietParam parses the 'quiet' query parameter.
func parseQuietParam(query url.Values) (bool, error) {
	q := query.Get("quiet")
	// If no quiet is specified, be loggy
	if q == "" {
		return false, nil
	}
	return strconv.ParseBool(q)
}

// parseTimeoutParam parses the 'timeout' query parameter.
func parseTimeoutParam(query url.Values) (time.Duration, error) {
	t := query.Get("timeout")
	// If no timeout is specified, use the default
	if t == "" {
		return DefaultTimeout, nil
	}
	return time.ParseDuration(t)
}

// parseGcIntervalParam parses the 'gc' query parameter.
func parseGcIntervalParam(query url.Values) (time.Duration, error) {
	t := query.Get("gc")
	// If no gc interval is specified, use the default
	if t == "" {
		return DefaultGcInterval, nil
	}
	return time.ParseDuration(t)
}

// gc probes the remote endpoint status and closes the listener if it's unresponsive.
func gc(ln *webteleport.Listener, interval time.Duration) {
	endpoint := ln.AsciiURL() + "/.well-known/health"
	for {
		time.Sleep(interval)
		resp, err := http.Get(endpoint)
		if err != nil {
			println("🛸 can't probe remote endpoint status. skipping...")
			continue
		}
		// if response is not 200, close the listener
		if resp.StatusCode != 200 {
			println("🛸 closing the listener because the server is unresponsive")
			ln.Close()
			break
		}
	}
}
