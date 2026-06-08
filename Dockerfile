FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget su-exec

WORKDIR /app

COPY --from=builder /server .
COPY scripts/docker-api-entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh && \
    adduser -D appuser && \
    mkdir -p /app/media && \
    chown -R appuser:appuser /app

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --start-period=60s --retries=6 \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/entrypoint.sh"]
