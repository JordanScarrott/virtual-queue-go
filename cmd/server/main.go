package main

import (
	"fmt"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"

	"red-duck/internal/adapters/config"
	httpAdapter "red-duck/internal/adapters/http"
	"red-duck/internal/adapters/temporal"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Temporal Client
	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%d", cfg.Temporal.Host, cfg.Temporal.Port),
	})
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}
	defer c.Close()

	// 3. Initialize Adapters
	queueClient := temporal.NewQueueClient(c, cfg.Temporal.TaskQueue)
	handler := httpAdapter.NewQueueHandler(queueClient)

	// 4. Start HTTP Server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
