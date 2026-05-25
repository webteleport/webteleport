package websocket

import (
	"net/url"
	"testing"
)

func TestResolveAddr(t *testing.T) {
	tests := []struct {
		name    string
		relay   *url.URL
		want    string
		wantErr bool
	}{
		{
			name: "http relay",
			relay: &url.URL{
				Scheme: "http",
				Host:   "relay.com:9090",
				Path:   "/pub",
			},
			want: "http://relay.com:9090/pub?x-websocket-upgrade=1",
		},
		{
			name: "https relay",
			relay: &url.URL{
				Scheme: "https",
				Host:   "relay.com:9090",
				Path:   "/pub",
			},
			want: "https://relay.com:9090/pub?x-websocket-upgrade=1",
		},
		{
			name: "with query params",
			relay: &url.URL{
				Scheme:   "https",
				Host:     "relay.com:9090",
				Path:     "/pub",
				RawQuery: "token=abc",
			},
			want: "https://relay.com:9090/pub?token=abc&x-websocket-upgrade=1",
		},
		{
			name: "with userinfo",
			relay: &url.URL{
				Scheme: "http",
				Host:   "relay.com:9090",
				Path:   "/pub",
				User:   url.UserPassword("user", "pass"),
			},
			want: "http://user:pass@relay.com:9090/pub?x-websocket-upgrade=1",
		},
		{
			name: "no path",
			relay: &url.URL{
				Scheme: "http",
				Host:   "relay.com:9090",
			},
			want: "http://relay.com:9090?x-websocket-upgrade=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveAddr("ignored", tt.relay)
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
