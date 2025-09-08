FROM golang:1.23-alpine

# Install dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o auth-service \
    ./cmd/main.go

EXPOSE 8081

CMD ["./auth-service"]