package temporal

import (
	"context"
	"log"
	"time"

	"red-duck/internal/analytics"
)

type QueueActivities struct {
	Tracker *analytics.Tracker
}

type JoinQueueParams struct {
	BusinessID      string
	QueueLength     int
	WaitTimeMinutes int
}

func (a *QueueActivities) JoinQueue(ctx context.Context, params JoinQueueParams) error {
	// Simulate Database Transaction
	// In a real application, we would save the JoinRequest to the database here.
	// Since no database code was found, we simulate success.
	log.Printf("Simulating DB transaction for JoinQueue businessID=%s", params.BusinessID)

	// Publish Event to NATS
	payload := map[string]interface{}{
		"type":              "QUEUE_UPDATED",
		"business_id":       params.BusinessID,
		"queue_length":      params.QueueLength,
		"wait_time_minutes": params.WaitTimeMinutes,
		"last_updated":      time.Now().Format(time.RFC3339),
	}

	err := a.Tracker.PublishEvent("joined", params.BusinessID, payload)
	if err != nil {
		log.Printf("Error publishing event to NATS: %v", err)
		// Constraint: Do NOT fail the Activity
		return nil
	}
	return nil
}

func (a *QueueActivities) LeaveQueue(ctx context.Context, params JoinQueueParams) error {
	// Simulate Database Transaction
	log.Printf("Simulating DB transaction for LeaveQueue businessID=%s", params.BusinessID)

	// Publish Event
	payload := map[string]interface{}{
		"type":              "QUEUE_UPDATED",
		"business_id":       params.BusinessID,
		"queue_length":      params.QueueLength,
		"wait_time_minutes": params.WaitTimeMinutes,
		"last_updated":      time.Now().Format(time.RFC3339),
	}

	err := a.Tracker.PublishEvent("left", params.BusinessID, payload)
	if err != nil {
		log.Printf("Error publishing event to NATS: %v", err)
		return nil
	}
	return nil
}
