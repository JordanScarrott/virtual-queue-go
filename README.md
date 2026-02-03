# Red Duck - Virtual Queue System

A Temporal-based Virtual Queue system implemented in Go using Hexagonal Architecture.

## Architecture

- **Domain:** Pure Go entities (`Queue`) representing the business logic.
- **Workflow:** `QueueWorkflow` acts as an Actor, maintaining state in memory and handling Signals (`Join`, `Leave`) and Queries (`GetQueue`).
- **Adapter:** `QueueClient` bridges the Domain Port to the Temporal SDK.
- **HTTP:** Simple REST API to interact with the queues.

## Prerequisites

- Go 1.23+
- Temporal Server running locally (e.g., via `temporal server start-dev`)

## Getting Started

### 1. Start the Worker

The worker executes the Queue Workflows.

```bash
go run cmd/worker/main.go
```

### 2. Start the HTTP Server

The server provides a REST API to interact with the queues.

```bash
go run cmd/server/main.go
```

## API Usage (Curl Examples)

Here are the common user stories you can execute using `curl`.

### Create a Queue

Starts a new Queue Workflow.

```bash
# Create a queue named "barbershop-1"
curl -X POST http://localhost:8080/queues/barbershop-1
```

### Join a Queue

Adds a user to the queue.

```bash
# User "alice" joins "barbershop-1"
curl -X POST http://localhost:8080/queues/barbershop-1/join \
     -H "Content-Type: application/json" \
     -d '{"userID": "alice"}'

# User "bob" joins "barbershop-1"
curl -X POST http://localhost:8080/queues/barbershop-1/join \
     -H "Content-Type: application/json" \
     -d '{"userID": "bob"}'
```

### Get Queue Status

Retrieves the current state of the queue (list of users).

```bash
curl http://localhost:8080/queues/barbershop-1
```

*Response:*
```json
{
  "ID": "barbershop-1",
  "Items": [
    {"UserID": "alice", "JoinedAt": "2023-10-27T10:00:00Z"},
    {"UserID": "bob", "JoinedAt": "2023-10-27T10:05:00Z"}
  ]
}
```

### Leave a Queue

Removes a user from the queue.

```bash
# User "alice" leaves "barbershop-1"
curl -X POST http://localhost:8080/queues/barbershop-1/leave \
     -H "Content-Type: application/json" \
     -d '{"userID": "alice"}'
```

### View in Temporal UI

Open the Temporal Web UI (usually at http://localhost:8233) to see the `QueueWorkflow` executions. You can view the Event History to see `Signal` and `Query` events.
