# Stage 1: Build the Go binary
FROM golang:1.24 AS builder

WORKDIR /app

# Download Go modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a statically linked binary for Linux/amd64 or linux/arm64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# Stage 2: Use a scratch (empty) image for minimal size
FROM scratch

WORKDIR /app

# Copy the binary and .env file (if needed)
COPY --from=builder /app/main .
COPY .env .env

# Expose port 80 (Render default)
EXPOSE 80

# Run the binary
CMD ["./main"]
