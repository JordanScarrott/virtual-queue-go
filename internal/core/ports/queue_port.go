package ports

import (
	"context"

	"example.com/virtual-queue/internal/core/domain"
)

type QueueService interface {
	CreateQueue(ctx context.Context, queueID string) error
	JoinQueue(ctx context.Context, queueID, userID string) error
	LeaveQueue(ctx context.Context, queueID, userID string) error
	GetQueueStatus(ctx context.Context, queueID string) (*domain.Queue, error)
}
