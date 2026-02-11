package analytics

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

// StartIngest listens to NATS and writes to Postgres
func StartIngest(nc *nats.Conn, db *pgxpool.Pool) {
	_, err := nc.Subscribe("events.>", func(msg *nats.Msg) {
		var payload EventPayload
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			log.Printf("Error unmarshaling event payload: %v", err)
			return
		}

		// Marshal properties map to JSON byte slice
		propsJSON, err := json.Marshal(payload.Properties)
		if err != nil {
			log.Printf("Error marshaling properties: %v", err)
			return
		}

		// Insert into Postgres
		// Note: we assume the table creation handles the specifics,
		// but typically we'd ensure idempotency if possible.
		// Here we just insert a new row for every event.
		query := `
			INSERT INTO analytics_events (event_type, business_id, user_id, timestamp, properties)
			VALUES ($1, $2, $3, $4, $5)
		`

		_, err = db.Exec(context.Background(), query,
			payload.Type,
			payload.BusinessID,
			payload.UserID,
			payload.Timestamp,
			propsJSON,
		)

		if err != nil {
			log.Printf("Error inserting event into database: %v", err)
			return
		}
	})

	if err != nil {
		log.Printf("Error subscribing to NATS: %v", err)
		return
	}

	log.Println("Analytics consumer started, listening on 'events.>'")
}
