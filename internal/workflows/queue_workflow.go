package workflows

import (
	"example.com/virtual-queue/internal/core/domain"
	"go.temporal.io/sdk/workflow"
)

func QueueWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	info := workflow.GetInfo(ctx)
	queueID := info.WorkflowExecution.ID

	// Initialize domain state
	state := domain.NewQueue(queueID)

	// Set Query Handler
	err := workflow.SetQueryHandler(ctx, QueryGetState, func() (domain.Queue, error) {
		return state.Snapshot(), nil
	})
	if err != nil {
		return err
	}

	// Setup Selector
	selector := workflow.NewSelector(ctx)

	joinCh := workflow.GetSignalChannel(ctx, SignalJoinQueue)
	selector.AddReceive(joinCh, func(c workflow.ReceiveChannel, more bool) {
		var signal JoinQueueSignal
		c.Receive(ctx, &signal)
		err := state.Enqueue(signal.UserID)
		if err != nil {
			logger.Error("Failed to enqueue user", "UserID", signal.UserID, "Error", err)
		} else {
			logger.Info("User joined queue", "UserID", signal.UserID)
		}
	})

	leaveCh := workflow.GetSignalChannel(ctx, SignalLeaveQueue)
	selector.AddReceive(leaveCh, func(c workflow.ReceiveChannel, more bool) {
		var signal LeaveQueueSignal
		c.Receive(ctx, &signal)
		err := state.Dequeue(signal.UserID)
		if err != nil {
			logger.Error("Failed to dequeue user", "UserID", signal.UserID, "Error", err)
		} else {
			logger.Info("User left queue", "UserID", signal.UserID)
		}
	})

	exitCh := workflow.GetSignalChannel(ctx, SignalExit)
	var exit bool
	selector.AddReceive(exitCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		exit = true
	})

	// Main Loop
	for {
		selector.Select(ctx)
		if exit {
			return nil
		}
	}
}
