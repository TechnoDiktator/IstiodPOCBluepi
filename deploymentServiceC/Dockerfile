# Use an official lightweight Golang image
FROM golang:1.22-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build Service C from the cmd/serviceC directory
RUN go build -o service-c ./cmd/serviceC

# Expose the service port
EXPOSE 8083

# Command to run the executable
CMD ["/app/service-c"]
