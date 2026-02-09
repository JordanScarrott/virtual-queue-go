# Red Duck - Virtual Queue Application

Red Duck is a backend service for a Virtual Queue application, built with **Go** and **Temporal.io**. It follows a strict **Hexagonal Architecture (Ports and Adapters)** to ensure separation of concerns and testability.

## Architecture

The project is structured according to the Ports and Adapters pattern:

- **Core (`internal/core`)**: Contains the business logic and domain entities. This layer is **Pure Go** and has zero external dependencies (no Temporal, no Database drivers).
    - `domain`: Entity definitions.
    - `ports`: Interfaces for driving (services) and driven (repositories/adapters) components.
- **Application (`internal/application`)**: Orchestrates the business logic using Use Cases or Command Handlers.
- **Adapters (`internal/adapters`)**: implementations of the ports.
    - `temporal`: Primary (Driving) adapter. Contains Workflows, Activities, and the Worker implementation.
    - `config`: Secondary (Driven) adapter. Handles configuration loading.
- **Cmd (`cmd`)**: Entry points for the application.
    - `worker`: The main executable that starts the Temporal Worker.

## Prerequisites

- **Go**: Latest stable version (e.g., 1.22+).
- **Temporal Server**: You need a running instance of the Temporal Server.
    - The easiest way is using the [Temporal CLI](https://docs.temporal.io/cli): `temporal server start-dev`.
- **NATS Server**: You need a running NATS server for real-time notifications.
    - Docker: `docker run -p 4222:4222 -p 8222:8222 nats:latest`

## Getting Started

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd red-duck
   ```

2. **Configuration**
   The application uses `application.yaml` for configuration. The default settings assume a local Temporal instance:
   ```yaml
   temporal:
     host: "localhost"
     port: 7233
     taskQueue: "red-duck-queue"
   
   nats:
     url: "nats://localhost:4222"
   ```

3. **Install Dependencies**
   ```bash
   go mod tidy
   ```

## Running the Worker

To start the Temporal Worker:

```bash
go run cmd/worker/main.go
```

Or build and run:

```bash
go build -o worker cmd/worker/main.go
./worker
```

You should see logs indicating the worker has started and is polling the task queue `red-duck-queue`.

## Features

### Real-Time Notifications
When a user joins or leaves a queue, the application publishes an event to NATS on the subject `events.queue.{business_id}`. This allows downstream services (e.g., websockets, analytics) to react to queue changes in real-time.

### Join Queue (Synchronous Update)

The application now supports a synchronous "Join Queue" operation using **Temporal Updates**. This ensures that the client receives immediate feedback on whether they successfully joined the queue (and their position) or if the request was rejected (e.g., duplicate user, closed queue).

## Running the HTTP Server

To start the HTTP Server (which exposes the Join Queue endpoint):

```bash
go run cmd/server/main.go
```

The server listens on `localhost:8080`.

### API Endpoints

For detailed API documentation, please refer to [docs/API.md](docs/API.md).

### Example Usage (cURL)

```bash
# 1. Create a queue
curl -X POST "http://localhost:8080/create_queue?business_id=biz1&queue_id=q1"

# 2. Join the queue
curl -X POST "http://localhost:8080/join_queue?business_id=biz1&queue_id=q1" \
  -H "Content-Type: application/json" \
  -d '{"userId": "user-1"}'

# 3. Check Status
curl -X GET "http://localhost:8080/queue_status?business_id=biz1&queue_id=q1"

# 4. Leave the queue
curl -X POST "http://localhost:8080/leave_queue?business_id=biz1&queue_id=q1" \
  -H "Content-Type: application/json" \
  -d '{"userId": "user-1"}'
```

## Triggering a Workflow

You can verify the worker is functioning correctly by starting the placeholder `NoOpWorkflow`.

### Option 1: Using the Temporal Web UI

1. Open your browser to the Temporal UI (default: [http://localhost:8233](http://localhost:8233)).
2. Navigate to the **Workflows** page.
3. Click the **Start Workflow** button (top right).
4. Fill in the following details:
   - **Workflow ID**: `test-noop-1` (or any unique string)
   - **Task Queue**: `red-duck-queue`
   - **Workflow Type**: `NoOpWorkflow`
5. Click **Start Workflow**.

You should see the workflow execution status change to **Completed** almost immediately. You can click on the workflow ID to view the execution history and verify that `NoOpActivity` was executed.

### Option 2: Using the Temporal CLI

```bash
temporal workflow start \
  --task-queue red-duck-queue \
  --type NoOpWorkflow \
  --workflow-id test-noop-cli
```

## Running Tests

Run all tests using:

```bash
go test ./...
```

## Project Structure

```
.
├── application.yaml          # Configuration file
├── cmd
│   └── worker                # Worker entry point
│       └── main.go
├── go.mod
├── internal
│   ├── adapters
│   │   ├── config            # Configuration adapter (Viper)
│   │   └── temporal          # Temporal Workflows & Activities
│   ├── application           # Use Cases
│   └── core
│       ├── domain            # Domain Entities (Pure Go)
│       └── ports             # Interfaces (Pure Go)
└── README.md
```
