package temporal

import (
	"time"

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
			// Activity Options
			container := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
			})

			// Call JoinQueue Activity (Simulated DB + NATS)
			var a *QueueActivities
			params := JoinQueueParams{
				BusinessID:      businessID,
				UserID:          req.UserID,
				QueueLength:     state.Len() + 1,
				WaitTimeMinutes: (state.Len() + 1) * 5, // Rough estimate
			}
			err := workflow.ExecuteActivity(container, a.JoinQueue, params).Get(container, nil)
			if err != nil {
				logger.Error("JoinQueue activity failed", "Error", err)
				return 0, err
			}

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

	// Define LeaveQueue Update
	leaveQueueUpdate := update.New[domain.JoinRequest, int]("LeaveQueue")
	err = workflow.SetUpdateHandler(ctx, leaveQueueUpdate.Name(),
		func(ctx workflow.Context, req domain.JoinRequest) (int, error) {
			// Check if user exists before calling activity
			if state.GetPosition(req.UserID) == 0 {
				return 0, domain.ErrUserNotFound
			}

			// Activity Options
			container := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: 10 * time.Second,
			})

			// Call LeaveQueue Activity
			var a *QueueActivities
			params := JoinQueueParams{
				BusinessID:      businessID,
				QueueLength:     state.Len() - 1,
				WaitTimeMinutes: (state.Len() - 1) * 5,
			}
			if params.QueueLength < 0 {
				params.QueueLength = 0
				params.WaitTimeMinutes = 0
			}

			err := workflow.ExecuteActivity(container, a.LeaveQueue, params).Get(container, nil)
			if err != nil {
				return 0, err
			}

			err = state.Dequeue(req.UserID)
			if err != nil {
				return 0, err
			}
			logger.Info("User left queue", "UserID", req.UserID)
			return state.Len(), nil
		},
	)
	if err != nil {
		return err
	}

	// Define GetStatus Query
	err = workflow.SetQueryHandler(ctx, "GetStatus", func() (domain.Queue, error) {
		return state.Snapshot(), nil
	})
	if err != nil {
		return err
	}

	// Keep the workflow running until "Exit" signal is received
	var exitSignal string
	workflow.GetSignalChannel(ctx, "Exit").Receive(ctx, &exitSignal)

	return nil
}
