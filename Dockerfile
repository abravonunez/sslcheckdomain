# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo 'dev') -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -o sslcheckdomain \
    ./cmd/sslcheckdomain

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 sslcheck && \
    adduser -D -u 1000 -G sslcheck sslcheck

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/sslcheckdomain /usr/local/bin/sslcheckdomain

# Set ownership
RUN chown -R sslcheck:sslcheck /app

# Switch to non-root user
USER sslcheck

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/sslcheckdomain"]
CMD ["--help"]
