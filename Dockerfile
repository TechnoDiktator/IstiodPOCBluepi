# Use an official lightweight Golang image
FROM golang:1.19-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build Service A from the cmd/serviceA directory
RUN go build -o service-a ./cmd/serviceA

# Expose the service port
EXPOSE 8082

# Command to run the executable
CMD ["/app/service-a"]
