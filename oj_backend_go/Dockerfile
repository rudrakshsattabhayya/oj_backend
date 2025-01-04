# Use the official Golang base image
FROM golang:1.23.3

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files (go.mod and go.sum)
COPY go.mod go.sum ./

# Download and cache Go modules dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the app will run on
EXPOSE 8001

# Run the Go application
CMD ["./main"]
