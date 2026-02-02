package ports

import (
	"context"

	"red-duck/internal/core/domain"
)

// DuckService defines the interface for duck-related operations.
// This is a placeholder port.
type DuckService interface {
	GetDuck(ctx context.Context, id string) (*domain.Duck, error)
}
