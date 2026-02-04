package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/client"

	"red-duck/internal/adapters/config"
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

	// 3. Initialize NATS Client
	nc, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		log.Printf("Warning: Failed to connect to NATS: %v. Publishing will be disabled.", err)
	} else {
		defer nc.Close()
	}

	// 4. Start Worker
	// This will block until the worker is stopped.
	err = temporal.StartWorker(c, cfg.Temporal.TaskQueue, nc)
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}
}
