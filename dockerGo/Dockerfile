# Build stage
FROM golang:1.21.0-alpine3.17 AS builder
WORKDIR /app
COPY . .

# Install dependencies and build the application
RUN go mod download
RUN go build -o main .

# Run stage (lightweight)
FROM alpine:3.18.3
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the API port
EXPOSE 9001

# Run the binary
CMD ["/app/main"]