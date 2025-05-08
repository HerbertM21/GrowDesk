FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy module files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd/server

# Create minimal production image
FROM alpine:latest

# Install dependencies
RUN apk --no-cache add tzdata ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/backend /app/backend

# Copy data files
COPY *.json ./data/

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DATA_DIR=/app/data
ENV MOCK_AUTH=true

# Run the application
CMD ["/app/backend", "-port", "8080", "-data-dir", "/app/data", "-mock-auth", "true"] 