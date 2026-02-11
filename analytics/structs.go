package analytics

import (
	"time"
)

type EventPayload struct {
	Type       string                 `json:"type"`
	BusinessID string                 `json:"business_id"`
	UserID     string                 `json:"user_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
}
