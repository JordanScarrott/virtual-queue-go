package temporal

import (
	"context"

	"red-duck/internal/core/domain"
	"go.temporal.io/sdk/client"
)

type QueueClient struct {
	client    client.Client
	taskQueue string
}

func NewQueueClient(c client.Client, taskQueue string) *QueueClient {
	return &QueueClient{
		client:    c,
		taskQueue: taskQueue,
	}
}

func (c *QueueClient) CreateQueue(ctx context.Context, queueID string) error {
	options := client.StartWorkflowOptions{
		ID:        queueID,
		TaskQueue: c.taskQueue,
	}

	_, err := c.client.ExecuteWorkflow(ctx, options, QueueWorkflow, queueID)
	return err
}

func (c *QueueClient) JoinQueue(ctx context.Context, queueID string, userID string) error {
	signal := JoinQueueSignal{UserID: userID}
	return c.client.SignalWorkflow(ctx, queueID, "", SignalJoinQueue, signal)
}

func (c *QueueClient) LeaveQueue(ctx context.Context, queueID string, userID string) error {
	signal := LeaveQueueSignal{UserID: userID}
	return c.client.SignalWorkflow(ctx, queueID, "", SignalLeaveQueue, signal)
}

func (c *QueueClient) GetQueue(ctx context.Context, queueID string) (*domain.Queue, error) {
	response, err := c.client.QueryWorkflow(ctx, queueID, "", QueryGetQueue)
	if err != nil {
		return nil, err
	}
	var q domain.Queue
	if err := response.Get(&q); err != nil {
		return nil, err
	}
	return &q, nil
}

func (c *QueueClient) GetPosition(ctx context.Context, queueID string, userID string) (int, error) {
	q, err := c.GetQueue(ctx, queueID)
	if err != nil {
		return 0, err
	}
	return q.Position(userID)
}
