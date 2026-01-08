package main

import (
	"log"
	"net/http"

	"budhapp.com/internal/config"
	"budhapp.com/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil || cfg == nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := server.New(*cfg)

	// Start server
	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
