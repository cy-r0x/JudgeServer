# Build stage
FROM golang:1.24.7-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o judgeserver .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/judgeserver .

# Copy schema files for migrations
COPY --from=builder /build/schema ./schema

# Copy .env file
COPY .env .env

# Expose port
EXPOSE 8080

# Run the application
CMD ["./judgeserver"]
