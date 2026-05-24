package webtransport

import (
	"net/url"
	"testing"
)

func TestDialAddrPreservesRelayPathAndUpgrade(t *testing.T) {
	relayURL := mustURL(t, "https://relay.example.com/demo?token=abc")

	addr, err := DialAddr("edge.example.com:4443", relayURL)
	if err != nil {
		t.Fatalf("DialAddr returned error: %v", err)
	}
	u := mustURL(t, addr)
	if u.Scheme != "https" || u.Host != "edge.example.com:4443" || u.Path != "/demo" {
		t.Fatalf("unexpected dial addr: %s", addr)
	}
	if got := u.Query().Get("token"); got != "abc" {
		t.Fatalf("unexpected token query: %q", got)
	}
	if got := u.Query().Get(UpgradeQuery); got != "1" {
		t.Fatalf("unexpected upgrade query: %q", got)
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
