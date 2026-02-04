package workflows

import (
	"strings"

	"red-duck/internal/core/domain"
	"go.temporal.io/sdk/workflow"
)

type QueueWorkflowInput struct {
	BusinessID string
	QueueID    string
}

func QueueWorkflow(ctx workflow.Context, input QueueWorkflowInput) error {
	logger := workflow.GetLogger(ctx)

	// Fallback for missing input (if upgrading from old workflow version without input)
	// In a real scenario, we'd handle versioning. Here we assume new executions.
	// If BusinessID is missing, we might try to parse it from WorkflowID if we used a naming convention there,
	// but passing it as input is cleaner.

	// Initialize domain state
	// Note: The Workflow ID is likely "BusinessID:QueueID", but the domain entity
	// stores the raw IDs.
	state := domain.NewQueue(input.QueueID, input.BusinessID)

	// If input was empty (e.g. migration), we might default.
	if state.ID == "" {
		// Attempt to parse from Workflow ID: "bizID:queueID"
		info := workflow.GetInfo(ctx)
		parts := strings.SplitN(info.WorkflowExecution.ID, ":", 2)
		if len(parts) == 2 {
			state.BusinessID = parts[0]
			state.ID = parts[1]
		} else {
			state.ID = info.WorkflowExecution.ID // Fallback
		}
	}

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
