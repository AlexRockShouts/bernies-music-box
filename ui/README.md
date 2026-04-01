# Bernie's Music Box UI

SvelteKit frontend for Bernie's Music Box.

Provides login, historical prompts view, and new prompt form with style/duration selectors.

## Quickstart (Development)

1. `cd ui`
2. `make dev`

- Dev server: http://localhost:5173
- Backend: http://localhost:8080 (api-gateway `make dev-run`)

Demo login: `demo` / `demo`

## Makefile Targets

- `make help` - This help
- `make install` - Install deps (`npm ci`)
- `make dev` - Dev server
- `make build` - Production build
- `make preview` - Preview prod build
- `make clean` - Clean artifacts
- `make docker-build` - Docker image (`ui:latest`)

## Docker

```bash
make docker-build
docker run -p 3000:80 ui:latest
```

View: http://localhost:3000

## API Integration

- `POST /tasks` {prompt, style, duration}
- `GET /history`
- Dummy auth (localStorage token)