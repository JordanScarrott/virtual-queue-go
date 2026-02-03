package temporal

import (
	"red-duck/internal/core/domain"
	"go.temporal.io/sdk/workflow"
)

// Signal and Query type names
const (
	SignalJoinQueue  = "JoinQueue"
	SignalLeaveQueue = "LeaveQueue"
	QueryGetQueue    = "GetQueue"
)

type JoinQueueSignal struct {
	UserID string
}

type LeaveQueueSignal struct {
	UserID string
}

func QueueWorkflow(ctx workflow.Context, queueID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("QueueWorkflow started", "queueID", queueID)

	// Initialize domain entity
	q := domain.NewQueue(queueID)

	// Setup query handler
	err := workflow.SetQueryHandler(ctx, QueryGetQueue, func() (*domain.Queue, error) {
		return q, nil
	})
	if err != nil {
		logger.Error("SetQueryHandler failed", "error", err)
		return err
	}

	joinChan := workflow.GetSignalChannel(ctx, SignalJoinQueue)
	leaveChan := workflow.GetSignalChannel(ctx, SignalLeaveQueue)

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(joinChan, func(c workflow.ReceiveChannel, more bool) {
			var signal JoinQueueSignal
			c.Receive(ctx, &signal)
			logger.Info("Received JoinQueue signal", "userID", signal.UserID)

			// Use workflow.Now(ctx) for deterministic time
			err := q.Join(signal.UserID, workflow.Now(ctx))
			if err != nil {
				logger.Warn("Failed to join queue", "userID", signal.UserID, "error", err)
			}
		})

		selector.AddReceive(leaveChan, func(c workflow.ReceiveChannel, more bool) {
			var signal LeaveQueueSignal
			c.Receive(ctx, &signal)
			logger.Info("Received LeaveQueue signal", "userID", signal.UserID)

			err := q.Leave(signal.UserID)
			if err != nil {
				logger.Warn("Failed to leave queue", "userID", signal.UserID, "error", err)
			}
		})

		// Wait for signal
		selector.Select(ctx)

		if ctx.Err() != nil {
			logger.Info("Workflow context cancelled", "error", ctx.Err())
			return ctx.Err()
		}
	}
}
