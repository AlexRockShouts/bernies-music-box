FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o api-gateway .

# Runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

WORKDIR /root/

COPY --from=builder /app/api-gateway .

# Expect volumes for:
# - server.crt, server.key
# - users.json
# - .env (or env vars)

EXPOSE 8443

ENTRYPOINT ["./api-gateway"]