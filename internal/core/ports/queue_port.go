package ports

import (
	"context"

	"red-duck/internal/core/domain"
)

type QueueService interface {
	CreateQueue(ctx context.Context, businessID, queueID string) error
	JoinQueue(ctx context.Context, businessID, queueID, userID string) error
	LeaveQueue(ctx context.Context, businessID, queueID, userID string) error
	GetQueueStatus(ctx context.Context, businessID, queueID string) (*domain.Queue, error)
}
