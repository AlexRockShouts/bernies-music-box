
FULL_IMAGE_NAME = docker.io/digitalfallout/music-box-runpod:latest

.PHONY: help install lint test run run-serverless clean docker-build docker-push

# Default target
help:
	@echo "Available commands:"
	@echo "  make install         - Install Python dependencies"
	@echo "  make lint            - Run syntax check and linting"
	@echo "  make test            - Run tests (placeholder)"
	@echo "  make run             - Run FastAPI app locally"
	@echo "  make run-serverless  - Run as RunPod serverless worker (local simulation)"
	@echo "  make clean           - Remove temporary files and __pycache__"
	@echo "  make docker-build    - Build Docker image ($(FULL_IMAGE_NAME))"
	@echo "  make docker-push     - Push Docker image to registry"

install:
	pip install -r requirements.txt
	pip install --no-deps git+https://github.com/HeartMuLa/heartlib.git

lint:
	python3 -m py_compile *.py
	# Optional: add ruff or flake8 if desired
	# ruff check .

test:
	# Placeholder for future tests
	@echo "No tests defined yet."

run:
	python3 main.py

run-serverless:
	python3 main.py --serverless

clean:
	find . -type d -name "__pycache__" -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete
	rm -rf output/*
	rm -rf .pytest_cache

docker-build:
	docker buildx build --platform linux/amd64 --load -t $(FULL_IMAGE_NAME) .

docker-push:
	docker buildx build --platform linux/amd64 --push -t $(FULL_IMAGE_NAME) .
