package analytics

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type Tracker struct {
	nc *nats.Conn
}

func NewTracker(nc *nats.Conn) *Tracker {
	return &Tracker{nc: nc}
}

func (t *Tracker) Track(eventType, businessID, userID string, props map[string]interface{}) error {
	payload := EventPayload{
		Type:       eventType,
		BusinessID: businessID,
		UserID:     userID,
		Timestamp:  time.Now().UTC(),
		Properties: props,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return t.nc.Publish(fmt.Sprintf("events.%s", eventType), data)
}
