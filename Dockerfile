# Step 1: Build the Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server main.go

# Step 2: Create the final lightweight image
FROM alpine:latest

# Install necessary libraries for SQLite
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/server .
# Copy the public assets
COPY --from=builder /app/public ./public

# Create a directory for the database (to be mounted as a volume)
RUN mkdir -p /app/data

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
# We tell the app to look for the DB in the /app/data folder
CMD ["./server"]
