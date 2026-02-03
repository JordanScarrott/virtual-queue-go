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
