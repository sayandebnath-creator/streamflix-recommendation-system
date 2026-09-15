# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API
RUN go build -o api ./cmd/api


# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Copy compiled Go binary
COPY --from=builder /app/api .

# Copy SQL migrations
COPY migrations ./migrations

# API port
EXPOSE 8080

# Start API
CMD ["./api"]