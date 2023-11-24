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

# Start a new stage from a Debian base image
FROM debian:slim  

# Update the package list and install Git and other necessary packages
RUN apt-get update && apt-get install -y git g++

# Install Go
RUN apt-get install -y golang-go

# Install Goose directly in the final image
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy the main application binary from the builder stage
COPY --from=builder /app/cmd/webserver/main /root/main

# Copy migrations
COPY /data/migrations /data/migrations

# Set environment variables for Goose
ENV GOOSE_DBSTRING="user=postgres password=root host=localhost dbname=thyrasec sslmode=disable"
ENV GOOSE_DRIVER="postgres"

# Copy the entrypoint script
COPY /scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/root/main"]
