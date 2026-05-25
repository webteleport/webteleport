package quic

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
			name: "host:port",
			addr: "example.com:8080",
			relay: &url.URL{
				Host: "relay.com:9090",
			},
			want: "example.com:8080",
		},
		{
			name: "host only",
			addr: "example.com",
			relay: &url.URL{
				Host: "relay.com:9090",
			},
			want: "example.com",
		},
		{
			name: ":port grafts relay host",
			addr: ":8080",
			relay: &url.URL{
				Host: "relay.com:9090",
			},
			want: "relay.com:8080",
		},
		{
			name: ":port with relay host no port",
			addr: ":8080",
			relay: &url.URL{
				Host: "relay.com",
			},
			want: "relay.com:8080",
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
