package secondary

import (
	"context"

	"example.com/virtual-queue/internal/core/domain"
	"example.com/virtual-queue/internal/core/ports"
	"example.com/virtual-queue/internal/workflows"
	"go.temporal.io/sdk/client"
)

type TemporalQueueClient struct {
	client client.Client
}

// Ensure TemporalQueueClient implements QueueService
var _ ports.QueueService = (*TemporalQueueClient)(nil)

func NewTemporalQueueClient(c client.Client) *TemporalQueueClient {
	return &TemporalQueueClient{client: c}
}

func (c *TemporalQueueClient) CreateQueue(ctx context.Context, queueID string) error {
	options := client.StartWorkflowOptions{
		ID:        queueID,
		TaskQueue: workflows.TaskQueue,
	}
	// We execute the workflow. It's a long running workflow, so we don't wait for result.
	// ExecuteWorkflow returns a Run, which we can ignore if we just want to start it.
	// If it already exists, this might fail depending on ID reuse policy, but that's expected.
	_, err := c.client.ExecuteWorkflow(ctx, options, workflows.QueueWorkflow)
	return err
}

func (c *TemporalQueueClient) JoinQueue(ctx context.Context, queueID, userID string) error {
	return c.client.SignalWorkflow(ctx, queueID, "", workflows.SignalJoinQueue, workflows.JoinQueueSignal{UserID: userID})
}

func (c *TemporalQueueClient) LeaveQueue(ctx context.Context, queueID, userID string) error {
	return c.client.SignalWorkflow(ctx, queueID, "", workflows.SignalLeaveQueue, workflows.LeaveQueueSignal{UserID: userID})
}

func (c *TemporalQueueClient) GetQueueStatus(ctx context.Context, queueID string) (*domain.Queue, error) {
	resp, err := c.client.QueryWorkflow(ctx, queueID, "", workflows.QueryGetState)
	if err != nil {
		return nil, err
	}
	var state domain.Queue
	if err := resp.Get(&state); err != nil {
		return nil, err
	}
	return &state, nil
}
