package temporal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/activity"

	"red-duck/internal/core/domain"
)

type Activities struct {
	NatsClient *nats.Conn
}

// JoinQueue handles a user joining a queue and publishes the event.
func (a *Activities) JoinQueue(ctx context.Context, businessID, queueID, userID string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("JoinQueue activity started", "businessID", businessID, "queueID", queueID, "userID", userID)

	// 1. Load/Create Queue (Simulated)
	// In a real app, this would come from a repository.
	// We create a fresh one here to demonstrate the logic.
	q := domain.NewQueue(businessID, queueID)

	// 2. Domain Logic
	if err := q.CanJoin(userID); err != nil {
		logger.Error("User cannot join queue", "error", err)
		return "", err
	}
	q.AddUser(userID)

	// 3. Publish to NATS
	// Construct the payload representing the current state.
	status := q.ToStatus()
	payload, err := json.Marshal(status)
	if err != nil {
		logger.Error("Failed to marshal queue status", "error", err)
		// We proceed as per "Non-Blocking" constraint
	} else {
		subject := fmt.Sprintf("queues.%s", businessID)
		if a.NatsClient != nil {
			err = a.NatsClient.Publish(subject, payload)
			if err != nil {
				logger.Error("Failed to publish to NATS", "subject", subject, "error", err)
				// We prefer the user joins successfully even if the real-time update drops.
			} else {
				logger.Info("Published queue update to NATS", "subject", subject)
			}
		} else {
			logger.Warn("NATS client is nil, skipping publish")
		}
	}

	return "joined", nil
}

// LeaveQueue handles a user leaving a queue and publishes the event.
func (a *Activities) LeaveQueue(ctx context.Context, businessID, queueID, userID string) (string, error) {
	// Implementation would be similar to JoinQueue:
	// 1. Load Queue
	// 2. Remove User
	// 3. Publish Status
	// For now, this is a stub to demonstrate the pattern.
	return "left", nil
}

// NoOpActivity is a placeholder activity.
func NoOpActivity(ctx context.Context) (string, error) {
	return "done", nil
}
