.PHONY: help install dev build preview clean docker-build

help:
	@echo "Makefile for UI SvelteKit project"
	@echo "  make help          - Show this help"
	@echo "  make install       - Install dependencies (npm ci)"
	@echo "  make dev           - Run dev server"
	@echo "  make build         - Build for production"
	@echo "  make preview       - Preview production build"
	@echo "  make clean         - Clean node_modules, .svelte-kit, build"
	@echo "  make docker-build  - Build Docker image (ui:latest)"

install:
	npm ci

dev:
	npm ci && npm run dev

build:
	npm ci && npm run build

preview:
	npm run preview

clean:
	rm -rf node_modules .svelte-kit build .vite

docker-build:
	docker buildx build --platform linux/amd64 --load -t ui:latest .
