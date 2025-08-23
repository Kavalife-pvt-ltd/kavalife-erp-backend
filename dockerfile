# Stage 1: Build the Go binary
FROM golang:1.24 AS builder

WORKDIR /app

# Download Go modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main ./cmd/server

# Stage 2: Use a scratch (empty) image for minimal size
FROM scratch

WORKDIR /app

# Copy binary and .env if needed
COPY --from=builder /app/main .
COPY .env .env

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
