FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application - NOTE: main.go is in cmd/ directory
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/main.go

# Final stage - smaller image
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/api-gateway .

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./api-gateway"]