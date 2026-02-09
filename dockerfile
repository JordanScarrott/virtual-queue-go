# Stage 1: Build the Go binary
# We use a specific version of Go on Alpine Linux (tiny & fast)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files first.
# Docker caches this layer, so we don't re-download modules if only code changes.
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build the application.
# We output a binary named 'virtual-queue-worker'
RUN go build -o virtual-queue-worker ./cmd/worker

# Stage 2: Create the production image
# We start fresh with a tiny Alpine image and only copy the compiled binary.
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/virtual-queue-worker .
COPY --from=builder /app/application.yaml .

# Expose the port your HTTP server listens on
EXPOSE 8080

# The command to run when the container starts
CMD ["./virtual-queue-worker"]