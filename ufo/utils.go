package ufo

import (
	"context"
	"fmt"
	"log/slog"
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
	slog.Info(fmt.Sprintf("🛸 listening on %s", ln.ClickableURL()))

	if u.Fragment == "" {
		slog.Info("🔓 publicly accessible without a password")
	} else {
		slog.Info("🔒 secured by password authentication")
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

// parseGcRetryParam parses the 'retry' query parameter.
func parseGcRetryParam(query url.Values) (int64, error) {
	r := query.Get("retry")
	// If no retry limit is specified, use the default
	if r == "" {
		return DefaultGcRetry, nil
	}
	return strconv.ParseInt(r, 10, 64)
}

// parsePersistParam parses the 'persist' query parameter.
func parsePersistParam(query url.Values) (bool, error) {
	p := query.Get("persist")
	// If no persist is specified, be ephemeral
	if p == "" {
		return false, nil
	}
	return strconv.ParseBool(p)
}

// gc probes the remote endpoint status and closes the listener if it's unresponsive.
func gc(ln *webteleport.Listener, interval time.Duration, limit int64) {
	endpoint := ln.AsciiURL() + "/.well-known/health"
	client := &http.Client{
		Timeout: interval,
	}
	for retry := limit; retry >= 0; {
		resp, err := client.Get(endpoint)
		// if request isn't successful, decrease retry
		if err != nil {
			retry -= 1
			werr := fmt.Errorf("🛸 failed to reach healthcheck endpoint (retry = %d): %v", retry, err)
			slog.Warn(werr.Error())
			time.Sleep(interval)
			continue
		}
		// if response stats code is not 200, decrease retry
		if resp.StatusCode != 200 {
			retry -= 1
			werr := fmt.Errorf("🛸 healthcheck endpoint returns status %d (retry = %d): %v", resp.StatusCode, retry, err)
			slog.Warn(werr.Error())
		} else {
			// otherwise reset retry to limit
			retry = limit
		}

		resp.Body.Close()
		time.Sleep(interval)
	}
	slog.Error("🛸 closing the listener")
	ln.Close()
}
