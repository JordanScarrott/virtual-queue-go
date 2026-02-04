package temporal

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// StartWorker initializes and starts the Temporal worker.
func StartWorker(c client.Client, taskQueue string) error {
	w := worker.New(c, taskQueue, worker.Options{})

	w.RegisterWorkflow(NoOpWorkflow)
	w.RegisterActivity(NoOpActivity)

	log.Printf("Starting worker on task queue: %s", taskQueue)
	return w.Run(worker.InterruptCh())
}
