# Build stage
FROM golang:alpine AS builder

# Git is required for fetching dependencies
RUN apk update && apk add --no-cache git

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory inside the builder
WORKDIR /app

# Copy and download dependencies using modules
COPY go.mod ./
RUN go mod download

# Copy project files to the working directory
COPY . .

# Build the Go application
RUN go build -o main

# Start a new stage from scratch
FROM scratch

# Copy the statically built Go binary from the builder
COPY --from=builder /app/main /main

# Copy the web folder containing HTML, assets, etc.
COPY --from=builder /app/web /web

# Expose port 8080
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["/main"]
