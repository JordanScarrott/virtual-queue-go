package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// StartIngest listens to NATS and writes to Postgres
func StartIngest(nc *nats.Conn, repo EventRepository) {
	_, err := nc.Subscribe("events.>", func(msg *nats.Msg) {
		if err := ProcessMessage(msg.Data, repo); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	})

	if err != nil {
		log.Printf("Error subscribing to NATS: %v", err)
		return
	}

	log.Println("Analytics consumer started, listening on 'events.>'")
}

// ProcessMessage handles a single message payload and inserts it into the repository
func ProcessMessage(data []byte, repo EventRepository) error {
	var payload EventPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("error unmarshaling event payload: %w", err)
	}

	// Insert into Repository
	if err := repo.InsertEvent(context.Background(), payload); err != nil {
		return fmt.Errorf("error inserting event into database: %w", err)
	}

	return nil
}
