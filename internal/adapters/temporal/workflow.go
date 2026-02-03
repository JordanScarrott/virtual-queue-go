package temporal

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// NoOpWorkflow is a placeholder workflow.
func NoOpWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, NoOpActivity).Get(ctx, &result)
	if err != nil {
		return err
	}

	workflow.GetLogger(ctx).Info("NoOpWorkflow completed", "result", result)
	return nil
}
