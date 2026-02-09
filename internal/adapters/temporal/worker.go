package temporal

import (
	"log"
	"red-duck/internal/analytics"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// StartWorker initializes and starts the Temporal worker.
func StartWorker(c client.Client, taskQueue string, tracker *analytics.Tracker) error {
	w := worker.New(c, taskQueue, worker.Options{})

	w.RegisterWorkflow(NoOpWorkflow)
	w.RegisterWorkflow(BusinessQueueWorkflow)

	// Register Activities
	activities := &QueueActivities{
		Tracker: tracker,
	}
	w.RegisterActivity(activities)
	w.RegisterActivity(NoOpActivity)

	log.Printf("Starting worker on task queue: %s", taskQueue)
	return w.Run(worker.InterruptCh())
}
