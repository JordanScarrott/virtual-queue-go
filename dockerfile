# Build Stage
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Install build dependencies if needed (e.g. for CGO, though we usually disable it for scratch/alpine)
# RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build Server
RUN go build -o /bin/server ./cmd/server

# Build Worker
# RUN go build -o /bin/worker ./cmd/worker

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy binaries from builder
COPY --from=builder /bin/server /bin/server
# COPY --from=builder /bin/worker /bin/worker

# Expose ports
EXPOSE 8081 8082