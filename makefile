install:
	go mod tidy
run:
	go run cmd/server/main.go

dev:
	@echo "▶ Starting Air hot reload..."
	air
build:
	go build -tags netgo -ldflags="-s -w" -o app ./cmd/server

buildAWS:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -tags netgo -ldflags="-s -w" -o app ./cmd/server