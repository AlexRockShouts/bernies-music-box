.PHONY: help build test run clean docker-build tidy lint dev-run

help:
	@echo "Makefile for api-gateway Go project"
	@echo "  make help          - Show this help"
	@echo "  make build         - Build api-gateway binary"
	@echo "  make test          - Run all tests"
	@echo "  make tidy          - Tidy go modules"
	@echo "  make lint          - Lint code (fmt, vet)"
	@echo "  make run           - Run with default env (prod-like)"
	@echo "  make dev-run       - Run for local dev"
	@echo "  make clean         - Remove binary"
	@echo "  make docker-build  - Build Docker image"

build:
	go build -o api-gateway main.go

test:
	go test ./... -v

tidy:
	go mod tidy

lint:
	go fmt ./...
	go vet ./...

run:
	go run main.go

dev-run:
	@echo "Using dev env: adjust as needed"
	AUTH_USERNAME=demo AUTH_PASSWORD=demo RUNPOD_URL=http://localhost:8000 JWT_SECRET=super-secret-key-change-in-prod PORT=8443 go run main.go

clean:
	rm -f api-gateway

docker-build:
	docker buildx build --platform linux/amd64 --load -t api-gateway:latest .
