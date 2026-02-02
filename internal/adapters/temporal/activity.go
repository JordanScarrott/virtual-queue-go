package temporal

import (
	"context"
)

// NoOpActivity is a placeholder activity.
func NoOpActivity(ctx context.Context) (string, error) {
	return "done", nil
}
