.PHONY: build run-gateway run-runpod run-ui help install-py lint-py clean-py

# Default target
help:
	@echo "Available commands:"
	@echo "  make build         - Build all projects"
	@echo "  make run-gateway   - Run Go API Gateway"
	@echo "  make run-runpod    - Run RunPod FastAPI (local simulation)"
	@echo "  make run-ui        - Run Svelte UI (dev mode)"
	@echo "  make install-py    - Install Python dependencies"
	@echo "  make lint-py       - Run linting for Python project"
	@echo "  make clean-py      - Clean Python temporary files"

build:
	cd api-gateway && go build -o api-gateway main.go
	cd ui && npm install && npm run build
	cd runpod-serverless && $(MAKE) docker-build

run-gateway:
	cd api-gateway && AUTH_USERNAME=admin AUTH_PASSWORD=password RUNPOD_URL=http://localhost:8000 go run main.go

run-runpod:
	cd runpod-serverless && $(MAKE) run

run-ui:
	cd ui && npm install && npm run dev

install-py:
	cd runpod-serverless && $(MAKE) install

lint-py:
	cd runpod-serverless && $(MAKE) lint

clean-py:
	cd runpod-serverless && $(MAKE) clean
