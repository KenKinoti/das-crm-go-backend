# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies for PostgreSQL
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application (CGO_ENABLED=1 for SQLite support)
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy any static files if needed
# Environment variables will be set directly in Cloud Run

# Expose port 8080 (Cloud Run default)
EXPOSE 8080

# Set environment variables for production
ENV GIN_MODE=release
ENV PORT=8080

# Run the application
CMD ["./main"]