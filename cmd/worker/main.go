package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"red-duck/auth"
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
	hostPort := os.Getenv("TEMPORAL_HOST_URL")
	if hostPort == "" {
		hostPort = client.DefaultHostPort // Fallback for local dev
	}
	c, err := client.Dial(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		log.Fatalf("Unable to create client: %v", err)
	}
	defer c.Close()

	// 3. Connect to NATS
	natsURL := cfg.Nats.URL
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("Failed to connect to NATS: %v", err)
	} else {
		defer nc.Close()
	}

	// 4. Initialize Tracker
	tracker := analytics.NewTracker(nc)

	// 5. Initialize Worker
	w := worker.New(c, cfg.Temporal.TaskQueue, worker.Options{})

	// Register Core Workflows & Activities
	w.RegisterWorkflow(temporal.NoOpWorkflow)
	w.RegisterWorkflow(temporal.BusinessQueueWorkflow)

	queueActivities := &temporal.QueueActivities{
		Tracker: tracker,
	}
	w.RegisterActivity(queueActivities)
	w.RegisterActivity(temporal.NoOpActivity)

	// Register Auth Workflows & Activities
	w.RegisterWorkflow(auth.LoginWorkflow)
	w.RegisterActivity(auth.SendMagicCode)
	w.RegisterActivity(auth.GenerateToken)

	// 6. Start HTTP Server (in a goroutine)
	go func() {
		http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req auth.LoginRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			workflowOptions := client.StartWorkflowOptions{
				ID:        "auth-" + req.Email,
				TaskQueue: cfg.Temporal.TaskQueue,
			}

			run, err := c.ExecuteWorkflow(context.Background(), workflowOptions, auth.LoginWorkflow, req.Email)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to start workflow: %v", err), http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{
				"workflowID": run.GetID(),
				"runID":      run.GetRunID(),
			})
		})

		http.HandleFunc("/auth/verify", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req auth.VerifyRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			workflowID := "auth-" + req.Email
			ctx := context.Background()

			// Signal the workflow
			err := c.SignalWorkflow(ctx, workflowID, "", "SubmitCode", req.Code)
			if err != nil {
				// If signal fails, workflow might not be running or completed
				http.Error(w, fmt.Sprintf("Failed to signal workflow: %v", err), http.StatusInternalServerError)
				return
			}

			// Wait for result
			workflowRun := c.GetWorkflow(ctx, workflowID, "")
			var token string
			err = workflowRun.Get(ctx, &token)
			if err != nil {
				http.Error(w, fmt.Sprintf("Authentication failed: %v", err), http.StatusUnauthorized)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{
				"token": token,
			})
		})

		log.Println("Starting HTTP server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// 7. Start Worker
	log.Printf("Starting worker on task queue: %s", cfg.Temporal.TaskQueue)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Unable to start worker: %v", err)
	}
}
