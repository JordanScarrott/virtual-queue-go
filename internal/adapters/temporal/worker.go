package temporal

import (
	"log"

	"github.com/nats-io/nats.go"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// StartWorker initializes and starts the Temporal worker.
func StartWorker(c client.Client, taskQueue string, nc *nats.Conn) error {
	w := worker.New(c, taskQueue, worker.Options{})

	w.RegisterWorkflow(NoOpWorkflow)

	// Register legacy NoOpActivity
	w.RegisterActivity(NoOpActivity)

	// Register new Activities with NATS injection
	activities := &Activities{
		NatsClient: nc,
	}
	w.RegisterActivity(activities)

	log.Printf("Starting worker on task queue: %s", taskQueue)
	return w.Run(worker.InterruptCh())
}
