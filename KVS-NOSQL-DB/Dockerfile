# Use the official Golang image
FROM golang:1.21-alpine

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Build the application
RUN go build -o kvs-server cmd/server/main.go

# Expose the ports
EXPOSE 9001 9002 9003

# Run the binary
ENTRYPOINT ["./kvs-server"]