package main

import (
	"log"

	"github.com/studyguides-com/study-guides-api/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // Ignore errors; default to environment

	// Initialize store
	appStore, err := store.NewStore()
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	// Create and start server manager
	manager := NewServerManager()
	manager.Start(appStore)
}
