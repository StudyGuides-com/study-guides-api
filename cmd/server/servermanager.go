package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/studyguides-com/study-guides-api/internal/store"
)

// ServerManager manages the server
type ServerManager struct {
	server *Server
}

// NewServerManager creates a new ServerManager instance
func NewServerManager() *ServerManager {
	return &ServerManager{
		server: NewServer(),
	}
}

// Start starts the server that handles both HTTP and gRPC traffic
func (sm *ServerManager) Start(appStore store.Store) {
	// Start server
	go sm.server.Start(appStore)

	// Wait for shutdown signal
	sm.waitForShutdown()
}

// waitForShutdown handles graceful shutdown
func (sm *ServerManager) waitForShutdown() {
	// Create a channel to receive shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	fmt.Printf("\n")
	log.Println("‼️ detected shutdown signal, starting shutdown...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully stop the server
	done := make(chan struct{})
	go func() {
		// Shutdown server
		if err := sm.server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		close(done)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-done:
		log.Println("✅ server gracefully stopped")
	case <-ctx.Done():
		log.Println("❌ graceful shutdown timeout, forcing stop")
		sm.server.ForceStop()
	}

	log.Println("🙌 server shutdown complete")
}
