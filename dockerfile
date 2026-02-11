# Stage 1: Build the Go binaries
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binaries
RUN go build -o server ./cmd/server
RUN go build -o worker ./cmd/worker

# Stage 2: Create the production image
FROM alpine:latest

WORKDIR /root/

# Copy the binaries from the builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/worker .
# Copy application.yaml if it exists, though environment variables are preferred
COPY --from=builder /app/application.yaml .

# Expose ports
EXPOSE 8080 8081

# Default command (can be overridden in docker-compose)
CMD ["./server"]