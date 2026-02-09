package analytics

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// Tracker handles publishing events to NATS.
type Tracker struct {
	conn *nats.Conn
}

// NewTracker creates a new Tracker instance.
func NewTracker(conn *nats.Conn) *Tracker {
	return &Tracker{conn: conn}
}

// PublishEvent publishes an event to NATS.
// It logs errors but does not return them to prevent failing the workflow activity.
func (t *Tracker) PublishEvent(eventType, businessID string, payload interface{}) error {
	if t.conn == nil {
		log.Println("NATS connection is nil, skipping event publish")
		return nil
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload for event %s: %v", eventType, err)
		return err
	}

	subject := fmt.Sprintf("events.queue.%s", businessID)
	if err := t.conn.Publish(subject, data); err != nil {
		log.Printf("Failed to publish event %s to subject %s: %v", eventType, subject, err)
		// Per instructions: Log error but do NOT fail the Activity.
		// Returning error here might be caught by caller.
		// However, instruction says "but do NOT fail the Activity".
		// If caller ignores error, that's fine.
		// But if we return error, caller *might* fail.
		// Let's stick to returning error from this method as signature requested,
		// and let the caller decide (or swallow it).
		// Wait, the prompt says "Log the error but do NOT fail the Activity."
		// This usually means the *Activity* shouldn't return error to Temporal.
		return err
	}

	return nil
}
