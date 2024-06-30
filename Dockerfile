# Use the official Golang image to build the app
FROM golang:1.22.2 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o server ./cmd/userapp/main.go

# Run tests with verbose output
RUN go test -v ./...

# Start a new stage from scratch
FROM debian:latest

# Install SQLite
RUN apt-get update && apt-get install -y sqlite3 && apt-get clean

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/server /app/server

# Copy the .env file to the working directory
COPY .env /app/.env

# Copy the SQLite database file if it exists
COPY user.db /app/user.db

# Copy the migrations directory
COPY migrations /app/migrations

# Set the working directory
WORKDIR /app

# Expose the port the app runs on
EXPOSE 50051

# Command to run the executable
CMD ["./server"]
