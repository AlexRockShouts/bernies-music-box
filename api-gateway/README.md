# API Gateway for Bernie's Music Box

## Overview
Go Gin-based API gateway with:
- JWT authentication (username/password from `users.json`)
- WebSocket (WSS) for real-time task updates per user
- slog structured JSON logging
- Kubernetes liveness/readiness probes
- Prometheus `/metrics` endpoint
- Task orchestration to RunPod serverless (async poll)

## Quick Start (Local Dev)

1. `cp .env.example .env` &amp; edit with RunPod creds
2. Generate self-signed TLS certs:
   ```
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 \\
     -keyout server.key -out server.crt -subj &quot;/CN=localhost&quot;
   ```
3. `go mod tidy`
4. `go run main.go`

- HTTPS: https://localhost:8443 (accept self-signed cert)
- WSS: wss://localhost:8443/api/ws

## Environment Variables (.env)
```
RUNPOD_URL=https://api.runpod.io/v2/{endpoint_id}/runsync  # or /run
AUTH_USERNAME=runpod_user
AUTH_PASSWORD=runpod_api_key
JWT_SECRET=your-super-secret-jwt-key-change-in-prod
PORT=8443
```

## Authentication
Users in `users.json` (bcrypt passwords, reload on restart):
```json
[
  {&quot;username&quot;:&quot;admin&quot;,&quot;password&quot;:&quot;$2a$10$...&quot;}
]
```
Default: admin/password

`POST /api/login` -> JWT token

## API Endpoints (HTTPS + `Authorization: Bearer &lt;token&gt;`)

### Tasks (owner-scoped)
- `POST /api/tasks` `{&quot;prompt&quot;:&quot;lyrics here&quot;}` → 202 `{id, status:&quot;pending&quot;}`
- `GET /api/tasks` → `[{id, prompt, status, result_url?, owner, timestamps}]`
- `GET /api/tasks/{id}` → task or 404

### WebSocket
- `GET /api/ws` (wss:// + auth header) → real-time owner task updates

### Health &amp; Monitoring
- `GET /livez` → 200 OK (liveness)
- `GET /readyz` → 200 OK (readiness)
- `GET /metrics` → Prometheus metrics

## Docker Build &amp; Run
```
docker build -t api-gateway .
docker run -p 8443:8443 -v $(pwd)/server.{crt,key}:/root/ -v $(pwd)/users.json:/root/ -v $(pwd)/.env:/root/.env api-gateway
```

## Kubernetes
- Probes: `livenessProbe: httpGet path:/livez`, `readinessProbe: path:/readyz`
- TLS: Secret for server.crt/key
- ConfigMap: users.json, .env vars
- Service: port 8443 HTTPS
- HPA/Metrics via Prometheus scrape /metrics

## Logging
slog JSON to stdout, request/task events.

## Notes
- Tasks: pending → processing (orchestrate RunPod) → completed/failed (poll)
- WS: per-user hub, initial tasks list + updates
- Memory stores, json persist on user update
- godotenv loads .env early