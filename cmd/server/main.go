package main

import (
	"fmt"
	"log"
	"net/http"

	"go.temporal.io/sdk/client"

	"red-duck/auth"
	"red-duck/internal/adapters/config"
	httpAdapter "red-duck/internal/adapters/http"
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

	// 3. Initialize HTTP Handler
	queueHandler := &httpAdapter.QueueHandler{
		Client:    c,
		TaskQueue: cfg.Temporal.TaskQueue,
	}

	// 4. Setup Routes
	http.HandleFunc("/create_queue", queueHandler.CreateQueue)
	http.HandleFunc("/join_queue", queueHandler.JoinQueue)
	http.HandleFunc("/leave_queue", queueHandler.LeaveQueue)
	http.HandleFunc("/queue_status", queueHandler.GetQueueStatus)

	// Admin/Staff Route with Auth
	http.HandleFunc("POST /queues/{id}/call-next", auth.WithAuth(queueHandler.CallNext))

	// 5. Start Server
	port := 8080
	log.Printf("Starting HTTP server on port %d...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
