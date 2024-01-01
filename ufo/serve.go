package ufo

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/webteleport/auth"
	"github.com/webteleport/utils"
)

// DefaultTimeout is the default dialing timeout for the UFO client.
var DefaultTimeout = 10 * time.Second

// DefaultGcInterval is the default garbage collection interval for the UFO client.
var DefaultGcInterval = 5 * time.Second

// DefaultGcRetry is the default garbage collection retry limit.
var DefaultGcRetry int64 = 3

// Serve starts a UFO server on the given station URL.
// GC: Automatically close the server when health check fails for the given times
// - gc: health check interval interval (0 for disable)
// - retry: Retry the health check for the given times
// - timeout: Automatically close the client when dial timeouts
// - quiet: Do not log the status of the server (false for loggy)
// - persist: Automatically restart the server when it becomes unresponsive (0 for disable)
func Serve(stationURL string, handler http.Handler) error {
	// Parse the station URL and inject client info
	u, err := createURLWithQueryParams(stationURL)
	if err != nil {
		return err
	}

	// Parse the 'quiet' query parameter
	quiet, err := parseQuietParam(u.Query())
	if err != nil {
		return err
	}

	// Parse the 'timeout' query parameter
	timeout, err := parseTimeoutParam(u.Query())
	if err != nil {
		return err
	}

	// Parse the 'gc' query parameter
	interval, err := parseGcIntervalParam(u.Query())
	if err != nil {
		return err
	}

	// Parse the 'retry' query parameter
	retry, err := parseGcRetryParam(u.Query())
	if err != nil {
		return err
	}

	// Parse the 'persist' query parameter
	persist, err := parsePersistParam(u.Query())
	if err != nil {
		return err
	}

	if !persist {
		// Serve with the parsed configuration
		return ServeWithConfig(&ServerConfig{
			StationURL: u,
			Handler:    handler,
			Timeout:    timeout,
			GcInterval: interval,
			GcRetry:    retry,
			Quiet:      quiet,
		})
	}

	// Serve with the parsed configuration
	for {
		err = ServeWithConfig(&ServerConfig{
			StationURL: u,
			Handler:    handler,
			Timeout:    timeout,
			GcInterval: interval,
			GcRetry:    retry,
			Quiet:      quiet,
		})
		if err != nil {
			slog.Warn(fmt.Sprintf("serve error: %v", err))
		}
		// retry indefinitely with 1s interval to avoid DoS
		time.Sleep(time.Second)
	}
}

// ServerConfig is the configuration for the UFO server.
type ServerConfig struct {
	StationURL *url.URL
	Handler    http.Handler
	Timeout    time.Duration
	GcInterval time.Duration
	GcRetry    int64
	Quiet      bool
}

// Serve starts a UFO server on the given station URL.
func ServeWithConfig(config *ServerConfig) error {
	u := config.StationURL

	// listen on the station URL with a timeout
	ln, err := listenWithTimeout(u.String(), config.Timeout)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	// log the status of the server
	if !config.Quiet {
		logServerStatus(ln, u)
	}

	// use the default serve mux if nil handler is provided
	if config.Handler == nil {
		config.Handler = http.DefaultServeMux
	}

	// close the listener when the server is unresponsive
	if config.GcInterval > 0 {
		go gc(ln, config.GcInterval, config.GcRetry)
	}

	// attach default middlewares
	var handler = config.Handler
	handler = auth.WithPassword(handler, u.Fragment)
	handler = utils.WellKnownHealthMiddleware(handler)

	err = http.Serve(ln, handler)
	if err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}
