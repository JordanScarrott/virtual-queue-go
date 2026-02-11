package workflows

import (
	"strings"
	"time"

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

	callNextCh := workflow.GetSignalChannel(ctx, SignalCallNext)
	selector.AddReceive(callNextCh, func(c workflow.ReceiveChannel, more bool) {
		var signal CallNextSignal
		c.Receive(ctx, &signal)

		ticket, err := state.ServeNext(signal.CounterID)
		if err != nil {
			logger.Info("Queue Empty", "CounterID", signal.CounterID)
		} else {
			logger.Info("Calling next user", "UserID", ticket.UserID, "CounterID", signal.CounterID)

			// Execute Activity to notify external systems (NATS)
			ao := workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			}
			ctx = workflow.WithActivityOptions(ctx, ao)

			params := CallNextParams{
				BusinessID: state.BusinessID,
				UserID:     ticket.UserID,
				CounterID:  signal.CounterID,
				Status:     string(ticket.Status),
			}

			// We use string name "CallNext" because we can't easily import activities struct due to potential cycles
			// or just simple string binding.
			// Ideally we would use a shared constant or dependency injection, but string is standard for loose coupling.
			err := workflow.ExecuteActivity(ctx, "CallNext", params).Get(ctx, nil)
			if err != nil {
				logger.Error("Failed to execute CallNext activity", "Error", err)
			}
		}
	})

	// Main Loop
	for {
		selector.Select(ctx)
		if exit {
			return nil
		}
	}
}
