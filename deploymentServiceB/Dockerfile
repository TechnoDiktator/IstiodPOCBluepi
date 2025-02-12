# Use an official lightweight Golang image
FROM golang:1.22-alpine


# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build Service B from the cmd/serviceB directory
RUN go build -o service-b ./cmd/serviceB

# Expose the service port
EXPOSE 8081

# Command to run the executable
CMD ["/app/service-b"]
