package ufo

import (
	"os"
	"testing"
)

func TestListenWithTimeout(t *testing.T) {
	// TODO
}

func TestCreateURLWithQueryParams(t *testing.T) {
	// Save original os.Args and restore after the test
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Set os.Args for testing
	os.Args = []string{"bash", "launch.sh"}

	// Test case: valid URL
	stationURL := "http://localhost:8080"
	expectedURL := "http://localhost:8080?args=bash&args=launch.sh&client=ufo"
	u, err := createURLWithQueryParams(stationURL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if u.String() != expectedURL {
		t.Errorf("Expected %s, got %s", expectedURL, u.String())
	}

	// Test case: invalid URL
	stationURL = "://invalid_url"
	_, err = createURLWithQueryParams(stationURL)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
