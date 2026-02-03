# Virtual Queue System

This project implements a "Virtual Queue" system using Temporal workflows and the "Workflow as an Actor" pattern. It follows Hexagonal Architecture principles to separate core domain logic from infrastructure concerns like Temporal.

## Architecture

The system is structured into layers:

- **`internal/core/domain`**: Pure Go business logic. Defines the `Queue` entity and its behavior. **No external dependencies.**
- **`internal/core/ports`**: Interfaces (Ports) that the core domain uses to interact with the outside world.
- **`internal/workflows`**: Temporal Workflow definitions. These act as the orchestration layer, binding the domain state to Temporal signals and queries.
- **`internal/adapters/secondary`**: Implementations of the ports (Adapters). `TemporalQueueClient` implements `QueueService` using the Temporal SDK.
- **`cmd/worker`**: The entry point for the Temporal Worker process.

### Workflow as an Actor

Each specific Queue (e.g., "Queue-123") is modeled as a single, long-running Temporal Workflow Execution.
- **Workflow ID**: The Queue ID.
- **State**: Held in memory within the workflow execution.
- **Signals**: `JoinQueue`, `LeaveQueue` modify the state.
- **Queries**: `GetState` exposes the current state.

## Prerequisites

- Go 1.20+
- [Temporal Server](https://docs.temporal.io/cli#starting-the-temporal-server) running locally (default: `localhost:7233`).

## Configuration

Configuration is loaded from `application.yaml` in the root directory.

```yaml
temporal:
  hostPort: "localhost:7233"
  taskQueue: "QUEUE_TASK_QUEUE"
```

## Running the Project

### 1. Start Temporal Server
Ensure your local Temporal server is running.
```bash
temporal server start-dev
```

### 2. Run Tests
Run unit tests for the domain logic and the workflow logic (using `testsuite`).
```bash
go test ./...
```

### 3. Run the Worker
The worker listens on the configured Task Queue and executes the workflows.
```bash
go run cmd/worker/main.go
```

## Usage (Client)

To interact with the system, you would typically use the `TemporalQueueClient` adapter in your application code (e.g., an HTTP handler).

Example:
```go
// Initialize Client
c, _ := client.Dial(client.Options{})
queueClient := secondary.NewTemporalQueueClient(c)

// Create/Start a Queue
queueClient.CreateQueue(ctx, "my-queue")

// Join Queue
queueClient.JoinQueue(ctx, "my-queue", "user-1")

// Get Status
status, _ := queueClient.GetQueueStatus(ctx, "my-queue")
fmt.Println(status.Users) // ["user-1"]
```

## Directory Structure

```
.
├── cmd/
│   └── worker/          # Worker entry point
├── internal/
│   ├── adapters/        # Adapters (infrastructure implementations)
│   ├── config/          # Configuration loading
│   ├── core/
│   │   ├── domain/      # Pure business logic
│   │   └── ports/       # Interfaces
│   └── workflows/       # Temporal Workflow definitions
├── application.yaml     # Configuration file
└── go.mod
```
