package webtransport

import (
	"net/url"
	"testing"
)

func TestParseControlLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		kind    string
		payload string
		ok      bool
	}{
		{name: "ignore empty", line: "", ok: false},
		{name: "ignore ping", line: "PING", ok: false},
		{name: "host", line: "HOST example.com:443", kind: "HOST", payload: "example.com:443", ok: true},
		{name: "error", line: "ERR denied", kind: "ERR", payload: "denied", ok: true},
		{name: "unknown", line: "NOPE value", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind, payload, ok := parseControlLine(tt.line)
			if ok != tt.ok || kind != tt.kind || payload != tt.payload {
				t.Fatalf("parseControlLine(%q) = (%q, %q, %v)", tt.line, kind, payload, ok)
			}
		})
	}
}

func TestDialAddrPreservesRelayPathAndUpgrade(t *testing.T) {
	relayURL := mustURL(t, "https://relay.example.com/demo?token=abc")

	addr, err := DialAddr("edge.example.com:4443", relayURL)
	if err != nil {
		t.Fatalf("DialAddr returned error: %v", err)
	}
	if addr != "https://edge.example.com:4443/demo?token=abc&x-webtransport-upgrade=1" &&
		addr != "https://edge.example.com:4443/demo?x-webtransport-upgrade=1&token=abc" {
		t.Fatalf("unexpected dial addr: %s", addr)
	}
}

func mustURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("url.Parse(%q): %v", raw, err)
	}
	return u
}
