package tunnel

import (
	"context"
	"net"
)

// Dialer represents anything that can initiate a tunnel session to a remote address.
type Dialer interface {
	// Dial connects to the remote address and returns a session.
	Dial(ctx context.Context, addr string) (Session, error)
}

// Listener represents anything that can listen for incoming tunnel sessions.
type Listener interface {
	// Listen binds to the given address and accepts incoming sessions.
	Listen(ctx context.Context, addr string) (net.Listener, error)
}

// Session represents a bidirectional multiplexed connection,
// allowing multiple logical streams over a single tunnel.
type Session interface {
	// Accept waits for an incoming stream initiated by the remote peer.
	Accept(ctx context.Context) (Stream, error)

	// Open initiates a new stream to the remote peer.
	Open(ctx context.Context) (Stream, error)

	// Close shuts down the entire session and all its streams.
	Close() error

	// Context returns the context associated with this session (e.g., for cancellation).
	Context() context.Context
}

// Stream represents a logical bidirectional channel, similar to a net.Conn,
// used for data transmission within a session.
type Stream interface {
	net.Conn
}

// Transport is a full-duplex abstraction that supports both dialing and listening.
type Transport interface {
	Dialer
	Listener
}
