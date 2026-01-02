# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

# Install git (required for fetching dependencies)
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the binary
# -o main: Output file name
# ./cmd/api: Location of main.go
RUN go build -o main ./cmd/api

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]