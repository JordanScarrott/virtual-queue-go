package temporal

import (
	"context"
	"log"

	"red-duck/analytics"
	"red-duck/internal/workflows"
)

type QueueActivities struct {
	Tracker analytics.EventTracker
}

type JoinQueueParams struct {
	BusinessID      string
	UserID          string
	QueueLength     int
	WaitTimeMinutes int
}

func (a *QueueActivities) JoinQueue(ctx context.Context, params JoinQueueParams) error {
	// Simulate Database Transaction
	// In a real application, we would save the JoinRequest to the database here.
	// Since no database code was found, we simulate success.
	log.Printf("Simulating DB transaction for JoinQueue businessID=%s userID=%s", params.BusinessID, params.UserID)

	// Publish Event to NATS via Tracker
	props := map[string]interface{}{
		"queue_length":   params.QueueLength,
		"estimated_wait": params.WaitTimeMinutes,
	}

	// Fire and forget tracking
	a.Tracker.Track("queue.joined", params.BusinessID, params.UserID, props)
	return nil
}

func (a *QueueActivities) LeaveQueue(ctx context.Context, params JoinQueueParams) error {
	// Simulate Database Transaction
	log.Printf("Simulating DB transaction for LeaveQueue businessID=%s userID=%s", params.BusinessID, params.UserID)

	// Publish Event via Tracker
	props := map[string]interface{}{
		"reason":       "user_quit",
		"queue_length": params.QueueLength, // Optional context
	}

	// Fire and forget tracking
	a.Tracker.Track("queue.left", params.BusinessID, params.UserID, props)
	return nil
}

func (a *QueueActivities) CallNext(ctx context.Context, params workflows.CallNextParams) error {
	// Simulate Database Transaction
	log.Printf("Simulating DB transaction for CallNext businessID=%s userID=%s counter=%s", params.BusinessID, params.UserID, params.CounterID)

	// Publish Event via Tracker
	props := map[string]interface{}{
		"status":      params.Status,
		"instruction": "Go to " + params.CounterID,
		"counter_id":  params.CounterID,
	}

	// Fire and forget tracking
	a.Tracker.Track("queue.called", params.BusinessID, params.UserID, props)
	return nil
}
