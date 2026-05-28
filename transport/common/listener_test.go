package common

import (
	"net"
	"net/url"
	"testing"
)

func TestListenerAddr(t *testing.T) {
	onionAddr := "6gdgizvqo5hhygxqgeocfq2klew4fy4qd2wmmkip2y2be5dkndtqhuyd"
	relayURL, _ := url.Parse("https://relay.example.com:443/path")
	type want struct {
		Network string
		Addr    string // Addr().String()
		String  string
	}
	tests := []struct {
		name  string
		ln    net.Addr
		want  want
		relay *url.URL
	}{
		{
			name: "relay nil",
			ln: &Listener{
				Scheme:  "https",
				Address: onionAddr,
			},
			want: want{
				Network: "https",
				String:  onionAddr,
				Addr:    onionAddr,
			},
			relay: nil,
		},
		{
			name: "relay set",
			ln: &Listener{
				Scheme:  "https",
				Address: onionAddr,
				Relay:   relayURL,
			},
			want: want{
				Network: "https",
				String:  onionAddr,
				Addr:    onionAddr + "." + relayURL.Host,
			},
			relay: relayURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.ln.(*Listener)
			if l.Network() != tt.want.Network {
				t.Errorf("Network() = %q, want %q", l.Network(), tt.want.Network)
			}
			if l.String() != tt.want.String {
				t.Errorf("String() = %q, want %q", l.String(), tt.want.String)
			}
			if l.Addr().String() != tt.want.Addr {
				t.Errorf("Addr().String() = %q, want %q", l.Addr().String(), tt.want.Addr)
			}
			if l.Relay != tt.relay {
				t.Errorf("Relay = %v, want %v", l.Relay, tt.relay)
			}
		})
	}
}
