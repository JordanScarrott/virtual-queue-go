package ports

import (
	"context"

	"red-duck/internal/core/domain"
)

type QueueService interface {
	CreateQueue(ctx context.Context, queueID string) error
	JoinQueue(ctx context.Context, queueID string, userID string) error
	LeaveQueue(ctx context.Context, queueID string, userID string) error
	GetQueue(ctx context.Context, queueID string) (*domain.Queue, error)
	GetPosition(ctx context.Context, queueID string, userID string) (int, error)
}
