# Use the official Golang image as the base image
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go application source code to the working directory
COPY . .

# Download and install any required Go modules
RUN go mod tidy

# Build the Go application
RUN go build -o app ./cmd/

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./app"]