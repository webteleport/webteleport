package endpoint

import "testing"

func TestExtractWebteleport(t *testing.T) {
	type testCase struct {
		name     string
		hostname string
		lines    []string
		want     []Endpoint
	}

	cases := []testCase{
		{
			name:     "empty lines",
			hostname: "example.com",
		},
		{
			name:     "webteleport Altsvc",
			hostname: "example.com",
			lines:    []string{`webteleport=":443"`},
			want:     []Endpoint{{Protocol: "webtransport", Addr: "example.com:443"}},
		},
		{
			name:     "non-webteleport filtered out",
			hostname: "example.com",
			lines:    []string{`h3=":443"`},
		},
		{
			name:     "explicit alt-host",
			hostname: "fallback.com",
			lines:    []string{`webteleport="specific.com:9090"`},
			want:     []Endpoint{{Protocol: "webtransport", Addr: "specific.com:9090"}},
		},
		{
			name:     "multiple entries one line",
			hostname: "example.com",
			lines:    []string{`webteleport=":443", webteleport=":444"`},
			want: []Endpoint{
				{Protocol: "webtransport", Addr: "example.com:443"},
				{Protocol: "webtransport", Addr: "example.com:444"},
			},
		},
		{
			name:     "multiple lines",
			hostname: "example.com",
			lines:    []string{`webteleport=":443"`, `webteleport=":9090"`},
			want: []Endpoint{
				{Protocol: "webtransport", Addr: "example.com:443"},
				{Protocol: "webtransport", Addr: "example.com:9090"},
			},
		},
		{
			name:     "malformed line",
			hostname: "example.com",
			lines:    []string{`garbage!!!`},
		},
		{
			name:     "mixed valid and malformed",
			hostname: "example.com",
			lines:    []string{`garbage!!!`, `webteleport=":443"`},
			want:     []Endpoint{{Protocol: "webtransport", Addr: "example.com:443"}},
		},
		{
			name:     "unquoted port is invalid altsvc",
			hostname: "example.com",
			lines:    []string{`webteleport=:443`},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ExtractWebteleport(c.hostname, c.lines...)
			if len(got) != len(c.want) {
				t.Fatalf("got %d endpoints, want %d: %v", len(got), len(c.want), got)
			}
			for i := range c.want {
				if got[i].Protocol != c.want[i].Protocol {
					t.Errorf("endpoint[%d] Protocol = %q, want %q", i, got[i].Protocol, c.want[i].Protocol)
				}
				if got[i].Addr != c.want[i].Addr {
					t.Errorf("endpoint[%d] Addr = %q, want %q", i, got[i].Addr, c.want[i].Addr)
				}
			}
		})
	}
}
