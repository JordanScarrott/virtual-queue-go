package temporal

import (
	"red-duck/internal/core/domain"
	"red-duck/internal/pkg/update"

	"go.temporal.io/sdk/workflow"
)

// BusinessQueueWorkflow manages a business queue.
func BusinessQueueWorkflow(ctx workflow.Context, businessID, queueID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("BusinessQueueWorkflow started", "BusinessID", businessID, "QueueID", queueID)

	state := domain.NewQueue(businessID, queueID)

	// Define the Update type
	joinQueueUpdate := update.New[domain.JoinRequest, int]("JoinQueue")

	// Set Update Handler with Options (Validator)
	err := workflow.SetUpdateHandlerWithOptions(ctx, joinQueueUpdate.Name(),
		func(ctx workflow.Context, req domain.JoinRequest) (int, error) {
			// Handler logic: Add user to state and return position
			position := state.AddUser(req.UserID)
			logger.Info("User joined queue", "UserID", req.UserID, "Position", position)
			return position, nil
		},
		workflow.UpdateHandlerOptions{
			Validator: func(ctx workflow.Context, req domain.JoinRequest) error {
				// Validator logic: Check if queue is closed or user already exists
				return state.CanJoin(req.UserID)
			},
		},
	)
	if err != nil {
		logger.Error("Failed to set update handler", "Error", err)
		return err
	}

	// Keep the workflow running until "Exit" signal is received
	var exitSignal string
	workflow.GetSignalChannel(ctx, "Exit").Receive(ctx, &exitSignal)

	return nil
}
