package analytics

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	InsertEvent(ctx context.Context, payload EventPayload) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) InsertEvent(ctx context.Context, payload EventPayload) error {
	query := `
		INSERT INTO analytics_events (event_type, business_id, user_id, timestamp, properties)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, payload.Type, payload.BusinessID, payload.UserID, payload.Timestamp, payload.Properties)
	return err
}
