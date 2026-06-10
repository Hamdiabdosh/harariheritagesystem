# Deploy on Coolify

Single-domain deployment: one URL serves the React app, API (`/api/v1/*`), and uploaded photos (`/media/*`).

## Architecture

```
Internet → Coolify proxy → caddy:80
                              ├─ /api/*     → api:8080
                              ├─ /health*   → api:8080
                              ├─ /media/*   → api:8080
                              └─ /*         → frontend:3000
```

| Service    | Role                          | Public? |
|------------|-------------------------------|---------|
| `caddy`    | Reverse proxy (entry point)   | **Yes** — assign your domain here |
| `frontend` | TanStack Start SSR app        | No      |
| `api`      | Go REST API                   | No      |
| `postgres` | PostgreSQL 16                 | No      |

## 1. Create the resource

1. In Coolify, add a new **Docker Compose** resource.
2. Point it at this repository.
3. Use `docker-compose.yml` as the compose file (no dev overlay).

## 2. Environment variables

Set these in Coolify → **Environment Variables** (see `.env.coolify.example`):

| Variable | Required | Example |
|----------|----------|---------|
| `JWT_SECRET` | Yes | 64+ char random string |
| `JWT_REFRESH_SECRET` | Yes | different 64+ char random string |
| `POSTGRES_USER` | Yes | `qirsmezgeb` |
| `POSTGRES_PASSWORD` | Yes | strong password |
| `POSTGRES_DB` | Yes | `qirsmezgeb` |
| `VITE_API_URL` | Yes | `https://your-domain.com/api/v1` |
| `ALLOWED_ORIGINS` | Yes | `https://your-domain.com` |

**Remove** if Coolify auto-injected them from local dev:

- `DB_URL`
- `PORT`
- `MEDIA_PATH`
- `POSTGRES_HOST` / `POSTGRES_PORT`
- `SERVICE_URL_*` (except what Coolify manages)

`VITE_API_URL` is a **build argument** for the frontend image. Changing it requires a **redeploy/rebuild**, not just a container restart.

## 3. Assign the domain

In Coolify → **Domains**, attach your URL **only to the `caddy` service** (port 80).

Do **not** assign domains to `api`, `frontend`, or `postgres`.

## 4. Deploy

Deploy / Redeploy. Coolify will:

1. Build `api` from `backend/Dockerfile`
2. Build `frontend` from `frontend/Dockerfile` (with `VITE_API_URL`)
3. Build `caddy` from `deploy/Dockerfile`
4. Start `postgres` with a persistent volume
5. Run DB migrations on API startup

## 5. Verify

```bash
# Proxy health (through Caddy)
curl https://your-domain.com/healthz

# API health (through Caddy)
curl https://your-domain.com/api/v1/health

# Frontend
open https://your-domain.com/login
```

Expected `healthz` response:

```json
{"status":"ok"}
```

## 6. Production checklist

- [ ] Change default seed user passwords (see root README)
- [ ] Use strong `JWT_SECRET` and `JWT_REFRESH_SECRET`
- [ ] Use a strong `POSTGRES_PASSWORD`
- [ ] Confirm `VITE_API_URL` matches your real domain
- [ ] Confirm `ALLOWED_ORIGINS` matches your real domain (no trailing slash)
- [ ] Domain is on `caddy` only

## Volumes

| Volume | Purpose |
|--------|---------|
| `qirsmezgeb_pgdata` | PostgreSQL data |
| `qirsmezgeb_media` | Uploaded record photos |

Both persist across redeploys on the same Coolify server.

## Local full-stack test (before Coolify)

```bash
cp .env.docker.example .env
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
```

- App via Caddy: http://localhost:8081/login
- Frontend direct: http://localhost:3000/login
- API direct: http://localhost:8080/healthz
