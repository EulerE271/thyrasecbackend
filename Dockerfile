# Start from a base image with Go installed
FROM golang:1.21.2 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Change to the directory containing main.go
WORKDIR /app/cmd/webserver

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:latest  

# Install dependencies required for runtime
RUN apk add --no-cache libc6-compat

# Copy Goose binary from builder stage
COPY --from=builder /go/bin/goose /usr/local/bin/goose

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/cmd/webserver/main .

# Copy migrations
COPY /data/migrations /data/migrations

# Set environment variables for Goose
ENV GOOSE_DBSTRING="user=postgres password=root host=localhost dbname=thyrasec sslmode=disable"
ENV GOOSE_DRIVER="postgres"

# Copy the entrypoint script
COPY /scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./main"]
