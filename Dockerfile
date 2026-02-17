# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the checkping service
RUN CGO_ENABLED=0 GOOS=linux go build -o checkping-service ./cmd/checkping

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /app/checkping-service .

# Copy endpoints JSON file
COPY --from=builder /app/internal/config/endpoints.json ./internal/config/endpoints.json

# Create directory structure for endpoints.json
RUN mkdir -p internal/config

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./checkping-service"]
