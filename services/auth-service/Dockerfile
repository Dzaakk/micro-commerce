FROM golang:1.21-alpine

# Install dependencies
RUN apk add --no-cache git make gcc musl-dev postgresql-client

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Install protoc and plugins for proto generation
RUN apk add --no-cache protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/micro/micro/v3/cmd/protoc-gen-micro@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy air config
COPY .air.toml ./

# Expose gRPC port
EXPOSE 8081

# Run with air for hot reload
CMD ["air", "-c", ".air.toml"]