# syntax=docker/dockerfile:1
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to cache deps
COPY ./backend/go.mod .
COPY ./backend/go.sum .
RUN go mod download

# Copy the rest of the source code
COPY ./backend/ .

# Build the Go app
RUN go build -o main .

# Expose the port (make sure your app listens on 8080)
EXPOSE 8080
EXPOSE 1883

# Run the binary
CMD ["./main"]
