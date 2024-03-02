# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and dependencies file
COPY go.mod go.sum ./

# Download Go modules and dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o websocket-server .

# Expose the port the WebSocket server listens on
EXPOSE 8080

# Command to run the WebSocket server executable
CMD ["./websocket-server"]
