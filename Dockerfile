FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget

WORKDIR /app

COPY --from=builder /server .

RUN adduser -D appuser && \
    mkdir -p /app/media && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --start-period=15s --retries=5 \
  CMD wget -q --spider http://localhost:8080/health || exit 1

CMD ["./server"]
