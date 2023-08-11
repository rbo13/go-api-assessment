package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/rbo13/go-api-assessment/src/http/server"
)

func TestServer_StartAndStop(t *testing.T) {
	s := server.New(
		server.WithAddress("127.0.0.1:8080"), // Use a different port to avoid conflicts
	)

	// Start the server in a separate goroutine
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Error starting server: %v", err)
		}
	}()

	// Allow some time for the server to start
	time.Sleep(100 * time.Millisecond)

	// Stop the server gracefully using a context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		t.Errorf("Error stopping server: %v", err)
	}
}
