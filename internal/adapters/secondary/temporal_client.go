package secondary

import (
	"context"
	"fmt"

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

func (c *TemporalQueueClient) getWorkflowID(businessID, queueID string) string {
	return fmt.Sprintf("%s:%s", businessID, queueID)
}

func (c *TemporalQueueClient) CreateQueue(ctx context.Context, businessID, queueID string) error {
	options := client.StartWorkflowOptions{
		ID:        c.getWorkflowID(businessID, queueID),
		TaskQueue: workflows.TaskQueue,
	}

	input := workflows.QueueWorkflowInput{
		BusinessID: businessID,
		QueueID:    queueID,
	}

	// We execute the workflow. It's a long running workflow, so we don't wait for result.
	// ExecuteWorkflow returns a Run, which we can ignore if we just want to start it.
	// If it already exists, this might fail depending on ID reuse policy, but that's expected.
	_, err := c.client.ExecuteWorkflow(ctx, options, workflows.QueueWorkflow, input)
	return err
}

func (c *TemporalQueueClient) JoinQueue(ctx context.Context, businessID, queueID, userID string) error {
	wfID := c.getWorkflowID(businessID, queueID)
	return c.client.SignalWorkflow(ctx, wfID, "", workflows.SignalJoinQueue, workflows.JoinQueueSignal{UserID: userID})
}

func (c *TemporalQueueClient) LeaveQueue(ctx context.Context, businessID, queueID, userID string) error {
	wfID := c.getWorkflowID(businessID, queueID)
	return c.client.SignalWorkflow(ctx, wfID, "", workflows.SignalLeaveQueue, workflows.LeaveQueueSignal{UserID: userID})
}

func (c *TemporalQueueClient) GetQueueStatus(ctx context.Context, businessID, queueID string) (*domain.Queue, error) {
	wfID := c.getWorkflowID(businessID, queueID)
	resp, err := c.client.QueryWorkflow(ctx, wfID, "", workflows.QueryGetState)
	if err != nil {
		return nil, err
	}
	var state domain.Queue
	if err := resp.Get(&state); err != nil {
		return nil, err
	}
	return &state, nil
}
