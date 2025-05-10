FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy sync-server files
COPY ./cmd/sync-server /app/

# Build the application
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o sync-server .

# Create minimal production image
FROM alpine:latest

# Install dependencies
RUN apk --no-cache add tzdata ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/sync-server /app/sync-server

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8000

# Set environment variables
ENV PORT=8000
ENV DATA_DIR=/app/data

# Run the application
CMD ["/app/sync-server"] 