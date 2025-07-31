install:
	go mod tidy
run:
	go run cmd/server/main.go

dev:
	@echo "▶ Starting Air hot reload..."
	air