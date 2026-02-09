package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/client"

	"red-duck/internal/adapters/config"
	"red-duck/internal/adapters/temporal"
	"red-duck/internal/analytics"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Temporal Client
	// Note: In a real production setup, you might want to configure TLS and other options here.
	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%d", cfg.Temporal.Host, cfg.Temporal.Port),
	})
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}
	defer c.Close()

	// 3. Connect to NATS
	// 3. Connect to NATS
	natsURL := cfg.Nats.URL
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		// Log error but maybe continue without NATS?
		// The requirement implies NATS is critical enough to be set up here,
		// but failure within activity should be non-fatal.
		// If NATS is down at startup, maybe we fail or just log?
		// "Connect to NATS on startup... Initialize the Tracker... Constraint: If NATS is down, log error but do NOT fail the Activity."
		// Startup failure usually means fatal for the worker heavily relying on it, but here it's notifications.
		// Let's log and proceed with nil connection if desired?
		// But Tracker implementation handles nil connection gracefully (I added that check).
		// So let's log error and continue.
		log.Printf("Failed to connect to NATS: %v", err)
	} else {
		defer nc.Close()
	}

	// 4. Initialize Tracker
	tracker := analytics.NewTracker(nc)

	// 5. Start Worker
	// This will block until the worker is stopped.
	err = temporal.StartWorker(c, cfg.Temporal.TaskQueue, tracker)
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}
}
