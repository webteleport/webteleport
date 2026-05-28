package common

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestReadHandshakeTimeout(t *testing.T) {
	// A reader that never produces data should trigger the timeout.
	r := io.NopCloser(strings.NewReader(""))
	// Override timeout for test to avoid waiting 10s.
	oldTimeout := HandshakeTimeout
	HandshakeTimeout = 50 * time.Millisecond
	defer func() { HandshakeTimeout = oldTimeout }()
	_, err := ReadHandshake(r)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("expected timeout error, got: %v", err)
	}
}

func TestReadHandshakeNoTimeout(t *testing.T) {
	// A reader that produces HOST immediately should NOT hit the timeout.
	oldTimeout := HandshakeTimeout
	HandshakeTimeout = 50 * time.Millisecond
	defer func() { HandshakeTimeout = oldTimeout }()
	r := strings.NewReader("HOST example.com:8080")
	got, err := ReadHandshake(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "example.com:8080" {
		t.Errorf("ReadHandshake() = %q, want %q", got, "example.com:8080")
	}
}

func TestReadHandshake(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    string
		wantErr bool
	}{
		{
			name:  "simple host",
			input: []string{"HOST example.com:8080"},
			want:  "example.com:8080",
		},
		{
			name:  "onion host",
			input: []string{"HOST 6gdgizvqo5hhygxqgeocfq2klew4fy4qd2wmmkip2y2be5dkndtqhuyd"},
			want:  "6gdgizvqo5hhygxqgeocfq2klew4fy4qd2wmmkip2y2be5dkndtqhuyd",
		},
		{
			name:    "err response",
			input:   []string{"ERR invalid auth"},
			wantErr: true,
		},
		{
			name:  "ping before host",
			input: []string{"PING", "HOST 127.0.0.1:8080"},
			want:  "127.0.0.1:8080",
		},
		{
			name:  "empty line before host",
			input: []string{"", "HOST localhost:9090"},
			want:  "localhost:9090",
		},
		{
			name:  "multiple ping before host",
			input: []string{"PING", "PING", "PING", "HOST host:1234"},
			want:  "host:1234",
		},
		{
			name:    "ping before err",
			input:   []string{"PING", "ERR bad token"},
			wantErr: true,
		},
		{
			name:  "unknown command ignored before host",
			input: []string{"FOO", "HOST ok:1"},
			want:  "ok:1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(strings.Join(tt.input, "\n"))
			got, err := ReadHandshake(r)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ReadHandshake() = %q, want %q", got, tt.want)
			}
		})
	}
}
