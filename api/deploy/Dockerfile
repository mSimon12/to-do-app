# syntax=docker/dockerfile:1

# Use the official Go image as the base image
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY ../src/go.* ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY ../src .

# Build the Go application
RUN go build -o to-do-api .

# Set the entry point command to run the built binary
CMD ["./to-do-api"]