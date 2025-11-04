# -------------------------
# Stage 1: Build Go binary
# -------------------------
FROM golang:1.24 AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build statically-linked Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# Optional: copy local .env for local testing
COPY .env /app/.env


# -------------------------
# Stage 2a: Local image with RIE
# -------------------------
FROM public.ecr.aws/lambda/provided:al2 AS local

WORKDIR /var/task

# Copy Go binary
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

# Copy RIE for local Lambda emulation
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/local/bin/aws-lambda-rie
RUN chmod +x /usr/local/bin/aws-lambda-rie

# ENTRYPOINT uses RIE for local testing
ENTRYPOINT ["/usr/local/bin/aws-lambda-rie", "/var/task/main"]
CMD []


# -------------------------
# Stage 2b: Production image (no RIE)
# -------------------------
FROM gcr.io/distroless/base-debian11 AS prod

# Copy Go binary
COPY --from=builder /app/main /var/task/main

# Optional: copy .env if needed
# COPY --from=builder /app/.env /var/task/.env

# ENTRYPOINT is just the Lambda binary for production
ENTRYPOINT ["/var/task/main"]
CMD []
