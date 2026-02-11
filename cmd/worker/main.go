package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"red-duck/analytics"
	"red-duck/auth"
	"red-duck/db"
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

	// 4. Initialize Database & Migrations
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Fallback for local dev
		dbURL = "postgres://user:password@localhost:5432/virtual_queue?sslmode=disable"
	}

	// Run Migrations
	migrationConn, err := db.ConnectAndMigrate(dbURL)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	migrationConn.Close(context.Background())

	// Create Connection Pool for Consumer
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer dbPool.Close()

	// 5. Initialize Analytics
	tracker := analytics.NewTracker(nc)

	// Start Ingest Consumer
	repo := analytics.NewPostgresRepository(dbPool)
	go analytics.StartIngest(nc, repo)

	// 6. Initialize Worker
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

		// Public Route: Join Queue (Guest Mode Support)
		http.HandleFunc("/queues/join", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Define request struct locally or use map
			var req struct {
				BusinessID string `json:"business_id"`
				UserID     string `json:"user_id"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// 1. Check for User ID (Guest Mode)
			if req.UserID == "" {
				req.UserID = uuid.New().String()
			}

			// 2. Start Workflow
			// preserving logic: workflow_id = queue-<business_id>-<user_id>
			workflowID := fmt.Sprintf("queue-%s-%s", req.BusinessID, req.UserID)
			workflowOptions := client.StartWorkflowOptions{
				ID:        workflowID,
				TaskQueue: cfg.Temporal.TaskQueue,
			}

			// BusinessQueueWorkflow signature: (ctx, businessID, queueID)
			// Assuming queueID is usually derived or same as businessID for simplicity in this context
			queueID := req.BusinessID

			run, err := c.ExecuteWorkflow(context.Background(), workflowOptions, temporal.BusinessQueueWorkflow, req.BusinessID, queueID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to start workflow: %v", err), http.StatusInternalServerError)
				return
			}

			// 3. Update Response
			json.NewEncoder(w).Encode(map[string]string{
				"workflow_id": run.GetID(),
				"run_id":      run.GetRunID(),
				"user_id":     req.UserID,
			})
		})

		// Protected Route Example
		http.HandleFunc("/req/me", auth.WithAuth(func(w http.ResponseWriter, r *http.Request) {
			// Extract user info from context (optional)
			// userID := r.Context().Value(auth.UserKey).(string)
			// role := r.Context().Value(auth.RoleKey).(string)

			json.NewEncoder(w).Encode(map[string]string{
				"message": "You are authenticated!",
				// "user_id": userID,
				// "role":    role,
			})
		}))

		log.Println("Starting HTTP server on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
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
