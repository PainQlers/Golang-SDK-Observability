# ==========================
# Stage 1 : Build
# ==========================
FROM golang:1.26-alpine AS builder

# Install certificates and git
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependency files first (better Docker cache)
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o server .

# ==========================
# Stage 2 : Runtime
# ==========================
FROM alpine:3.22

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Copy binary
COPY --from=builder /app/server .

# Copy required runtime files
COPY --from=builder /app/skills.md .

# Change ownership
RUN chown -R appuser:appgroup /app

USER appuser

# Application API
EXPOSE 8080

# Prometheus Metrics
EXPOSE 2112

ENTRYPOINT ["./server"]