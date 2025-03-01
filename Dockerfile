# Use the official Golang image as the base image
FROM golang:1.23.6-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o simple-messenger .

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/simple-messenger .
COPY --from=builder /app/config.json .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./simple-messenger"]