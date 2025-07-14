# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create application directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o waterlogger cmd/waterlogger/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S waterlogger && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G waterlogger -g waterlogger waterlogger

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/waterlogger .

# Copy configuration template
COPY config.yaml config.example.yaml

# Copy web assets
COPY --chown=waterlogger:waterlogger web/ web/

# Create directory for database
RUN mkdir -p /app/data && chown -R waterlogger:waterlogger /app/data

# Change ownership of application files
RUN chown -R waterlogger:waterlogger /app

# Switch to non-root user
USER waterlogger

# Expose port
EXPOSE 2341

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:2341/ || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV CONFIG_PATH=/app/config.yaml

# Run the application
CMD ["./waterlogger"]