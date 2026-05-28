package common

import "expvar"

// StreamMetrics holds stream counters exported via expvar.
type StreamMetrics struct {
	Accepted *expvar.Int
	Opened   *expvar.Int
	Closed   *expvar.Int
}

// NewStreamMetrics creates StreamMetrics with the given prefix.
// Expvar keys are "{prefix}_streams_accepted", "{prefix}_streams_opened", "{prefix}_streams_closed".
func NewStreamMetrics(prefix string) *StreamMetrics {
	return &StreamMetrics{
		Accepted: expvar.NewInt(prefix + "_streams_accepted"),
		Opened:   expvar.NewInt(prefix + "_streams_opened"),
		Closed:   expvar.NewInt(prefix + "_streams_closed"),
	}
}
