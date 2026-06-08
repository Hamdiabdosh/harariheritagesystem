#!/bin/sh
set -e

# Always build DB_URL from Postgres service vars so a stale Coolify DB_URL
# (e.g. localhost:5434 from local dev) cannot override the in-cluster hostname.
export DB_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST:-postgres}:${POSTGRES_PORT:-5432}/${POSTGRES_DB}?sslmode=disable"

mkdir -p /app/media
chown -R appuser:appuser /app/media

exec su-exec appuser ./server
