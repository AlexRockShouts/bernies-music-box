.PHONY: build run-gateway run-runpod run-ui help

# Default target
help:
	@echo "Available commands:"
	@echo "  make build         - Build all projects"
	@echo "  make run-gateway   - Run Go API Gateway"
	@echo "  make run-runpod    - Run RunPod FastAPI (local simulation)"
	@echo "  make run-ui        - Run Svelte UI (dev mode)"

build:
	cd api-gateway && go build -o api-gateway main.go
	cd ui && npm install && npm run build
	cd runpod-serverless && docker build -t music-box-runpod .

run-gateway:
	cd api-gateway && AUTH_USERNAME=admin AUTH_PASSWORD=password RUNPOD_URL=http://localhost:8000 go run main.go

run-runpod:
	cd runpod-serverless && AUTH_USERNAME=admin AUTH_PASSWORD=password python3 main.py

run-ui:
	cd ui && npm install && npm run dev
