package tcp

import (
	"net/url"
	"testing"
)

func TestResolveAddr(t *testing.T) {
	tests := []struct {
		name      string
		addr      string
		relayHost string
		want      string
		wantErr   bool
	}{
		{
			name:      "host and port",
			addr:      "example.com:8080",
			relayHost: "relay.com:9090",
			want:      "example.com:8080",
		},
		{
			name:      "host only",
			addr:      "example.com",
			relayHost: "relay.com:9090",
			want:      "example.com",
		},
		{
			name:      "ipv4 and port",
			addr:      "192.168.1.1:3000",
			relayHost: "relay.com:9090",
			want:      "192.168.1.1:3000",
		},
		{
			name:      "port only - grafts relay host",
			addr:      ":8080",
			relayHost: "relay.com:9090",
			want:      "relay.com:8080",
		},
		{
			name:      "port only with relay host no port",
			addr:      ":8080",
			relayHost: "relay.com",
			want:      "relay.com:8080",
		},
		{
			name:      "url with http scheme",
			addr:      "http://example.com:8080/path",
			relayHost: "relay.com:9090",
			want:      "example.com:8080",
		},
		{
			name:      "url with https scheme and query",
			addr:      "https://example.com:8080/path?q=1",
			relayHost: "relay.com:9090",
			want:      "example.com:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relayURL := &url.URL{Host: tt.relayHost}
			got, err := ResolveAddr(tt.addr, relayURL)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ResolveAddr() expected error, got none")
				}
				return
			}
			if err != nil {
				t.Errorf("ResolveAddr() unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveAddr() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveAddrRelayURLPreserved(t *testing.T) {
	relayURL := &url.URL{
		Scheme:   "https",
		Host:     "relay.com:9090",
		Path:     "/path",
		RawQuery: "a=1",
	}
	got, err := ResolveAddr("endpoint.com:1234", relayURL)
	if err != nil {
		t.Fatalf("ResolveAddr() unexpected error: %v", err)
	}
	if got != "endpoint.com:1234" {
		t.Errorf("ResolveAddr() = %q, want %q", got, "endpoint.com:1234")
	}
}

func TestResolveAddrGraftPort(t *testing.T) {
	relayURL := &url.URL{Host: "relay.com:9090"}
	got, err := ResolveAddr(":8080", relayURL)
	if err != nil {
		t.Fatalf("ResolveAddr() unexpected error: %v", err)
	}
	if got != "relay.com:8080" {
		t.Errorf("ResolveAddr() = %q, want %q", got, "relay.com:8080")
	}
}

func TestResolveAddrSchemeIgnored(t *testing.T) {
	relayURL := &url.URL{Host: "relay.com:9090"}
	tests := []struct {
		addr string
		want string
	}{
		{"http://example.com:8080", "example.com:8080"},
		{"https://example.com:8080", "example.com:8080"},
		{"ws://example.com:8080", "example.com:8080"},
		{"tcp://example.com:8080", "example.com:8080"},
		{"unix://example.com:8080", "example.com:8080"},
	}
	for _, tt := range tests {
		t.Run(tt.addr, func(t *testing.T) {
			got, err := ResolveAddr(tt.addr, relayURL)
			if err != nil {
				t.Errorf("ResolveAddr() unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ResolveAddr() = %q, want %q", got, tt.want)
			}
		})
	}
}
