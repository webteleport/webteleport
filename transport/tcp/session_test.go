package tcp

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/hashicorp/yamux"
)

func TestTcpTransportWithNetPipe(t *testing.T) {
	// Create a pair of connected pipes
	client, server := net.Pipe()

	// Create client and server sessions using the pipes
	clientConfig := yamux.DefaultConfig()
	serverConfig := yamux.DefaultConfig()

	// Create server session
	serverYamuxSession, err := yamux.Server(server, serverConfig)
	if err != nil {
		t.Fatalf("Failed to create server session: %v", err)
	}
	serverSession := &TcpSession{Session: serverYamuxSession}

	// Create client session
	clientYamuxSession, err := yamux.Client(client, clientConfig)
	if err != nil {
		t.Fatalf("Failed to create client session: %v", err)
	}
	clientSession := &TcpSession{Session: clientYamuxSession}

	// Test opening and accepting streams
	ctx := context.Background()

	// Client opens a stream
	t.Run("OpenStream", func(t *testing.T) {
		clientStream, err := clientSession.Open(ctx)
		if err != nil {
			t.Fatalf("Failed to open client stream: %v", err)
		}

		// Server accepts the stream
		serverStream, err := serverSession.Accept(ctx)
		if err != nil {
			t.Fatalf("Failed to accept server stream: %v", err)
		}

		// Test data transfer
		testData := []byte("hello from client")
		go func() {
			_, err := clientStream.Write(testData)
			if err != nil {
				t.Errorf("Failed to write to client stream: %v", err)
			}
			clientStream.Close()
		}()

		receivedData := make([]byte, len(testData))
		n, err := io.ReadFull(serverStream, receivedData)
		if err != nil {
			t.Fatalf("Failed to read from server stream: %v", err)
		}
		if n != len(testData) {
			t.Fatalf("Expected to read %d bytes, got %d", len(testData), n)
		}
		if string(receivedData) != string(testData) {
			t.Fatalf("Expected %q, got %q", string(testData), string(receivedData))
		}

		// Clean up
		clientStream.Close()
		serverStream.Close()
	})

	// Test bidirectional communication
	t.Run("BidirectionalCommunication", func(t *testing.T) {
		clientStream, err := clientSession.Open(ctx)
		if err != nil {
			t.Fatalf("Failed to open client stream: %v", err)
		}

		serverStream, err := serverSession.Accept(ctx)
		if err != nil {
			t.Fatalf("Failed to accept server stream: %v", err)
		}

		// Client to server
		clientMsg := []byte("hello from client")
		go func() {
			_, err := clientStream.Write(clientMsg)
			if err != nil {
				t.Errorf("Failed to write to client stream: %v", err)
			}
		}()

		clientMsgReceived := make([]byte, len(clientMsg))
		_, err = io.ReadFull(serverStream, clientMsgReceived)
		if err != nil {
			t.Fatalf("Failed to read client message: %v", err)
		}
		if string(clientMsgReceived) != string(clientMsg) {
			t.Fatalf("Expected %q, got %q", string(clientMsg), string(clientMsgReceived))
		}

		// Server to client
		serverMsg := []byte("hello from server")
		go func() {
			_, err := serverStream.Write(serverMsg)
			if err != nil {
				t.Errorf("Failed to write to server stream: %v", err)
			}
		}()

		serverMsgReceived := make([]byte, len(serverMsg))
		_, err = io.ReadFull(clientStream, serverMsgReceived)
		if err != nil {
			t.Fatalf("Failed to read server message: %v", err)
		}
		if string(serverMsgReceived) != string(serverMsg) {
			t.Fatalf("Expected %q, got %q", string(serverMsg), string(serverMsgReceived))
		}

		// Clean up
		clientStream.Close()
		serverStream.Close()
	})

	// Clean up
	clientSession.Close()
	serverSession.Close()
}

func TestTcpSessionCloseReturnsError(t *testing.T) {
	client, server := net.Pipe()

	serverYamuxSession, err := yamux.Server(server, yamux.DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create server session: %v", err)
	}
	clientYamuxSession, err := yamux.Client(client, yamux.DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create client session: %v", err)
	}

	session := &TcpSession{Session: clientYamuxSession}

	// Close should return nil on success, not http.ErrServerClosed
	if err := session.Close(); err != nil {
		t.Fatalf("Expected nil error on close, got: %v", err)
	}

	// After close, Open should fail
	_, err = session.Open(context.Background())
	if err == nil {
		t.Fatal("Expected error when opening stream on closed session")
	}

	serverYamuxSession.Close()
}
