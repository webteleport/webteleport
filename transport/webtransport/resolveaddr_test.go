package webtransport

import (
	"net/url"
	"testing"
)

func TestResolveAddr(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		relay   *url.URL
		want    string
		wantErr bool
	}{
		{
			name: "host:port with path",
			addr: "example.com:8080",
			relay: &url.URL{
				Host: "relay.com:9090",
				Path: "/pub",
			},
			want: "https://example.com:8080/pub?x-webtransport-upgrade=1",
		},
		{
			name: "host:port no path",
			addr: "example.com:8080",
			relay: &url.URL{
				Host: "relay.com:9090",
			},
			want: "https://example.com:8080?x-webtransport-upgrade=1",
		},
		{
			name: ":port grafts localhost",
			addr: ":8080",
			relay: &url.URL{
				Host: "relay.com:9090",
				Path: "/pub",
			},
			want: "https://relay.com:8080/pub?x-webtransport-upgrade=1",
		},
		{
			name: "host only",
			addr: "example.com",
			relay: &url.URL{
				Host: "relay.com:9090",
				Path: "/pub",
			},
			want: "https://example.com/pub?x-webtransport-upgrade=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveAddr(tt.addr, tt.relay)
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
