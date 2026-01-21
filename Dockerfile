# Multi-stage build for Coto
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=$(git describe --tags 2>/dev/null || echo 'dev')" -o coto .

# Create final minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/coto .

# Create volume for data
VOLUME ["/data"]

# Set entrypoint
ENTRYPOINT ["./coto"]
CMD ["--help"]