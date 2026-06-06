# Qirs Mezgeb API

Backend REST API for the Harari Heritage Registry System (ቅርስ መዝገብ).

## Prerequisites

- Go 1.22+
- PostgreSQL 16

## Setup

### Option A — Docker (recommended for local dev)

Port `5432` is often already in use by other projects. This repo includes a dedicated Postgres on **5434**:

```bash
docker compose up -d
cp .env.example .env
go mod tidy
go run ./cmd/server
```

### Option B — Existing PostgreSQL

1. Copy the environment file and set `DB_URL` to your real credentials:

```bash
cp .env.example .env
```

2. Create the database:

```bash
createdb qirsmezgeb
```

3. Download dependencies and run the server:

```bash
go mod tidy
go run ./cmd/server
```

The server listens on `http://localhost:8080` by default.

## Health Check

```bash
curl http://localhost:8080/health
```

Expected response when the database is reachable:

```json
{ "status": "ok", "db": "connected" }
```

Migrations run automatically on startup. The seed migration creates a default manager account:

| Email | Password |
|-------|----------|
| `admin@qirsmezgeb.gov.et` | `Admin1234` |

Change this password before deploying to production.

## Auth Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/login` | None | Email + password → access & refresh tokens |
| POST | `/api/v1/auth/refresh` | None | Refresh token → new access token |
| POST | `/api/v1/auth/logout` | Bearer JWT | Invalidates refresh token |

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@qirsmezgeb.gov.et","password":"Admin1234"}'
```

## User Management (Manager only)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/users` | Manager | List users (paginated, filter by role/is_active) |
| POST | `/api/v1/users` | Manager | Create user |
| PUT | `/api/v1/users/:id` | Manager | Update user |
| DELETE | `/api/v1/users/:id` | Manager | Deactivate user (soft delete) |
| GET | `/api/v1/users/me` | Any role | Get own profile |
| PUT | `/api/v1/users/me/language` | Any role | Update own language (`am` / `en`) |

## Immovable Records

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/records/immovable` | Registrar | Create draft (auto-generates `ET-HR-AN-I-YYYY-NNNN`) |
| GET | `/api/v1/records/immovable` | All roles | List records (role-filtered, paginated) |
| GET | `/api/v1/records/immovable/:id` | All roles | Get record detail |
| PUT | `/api/v1/records/immovable/:id` | Registrar | Update own draft/returned record |
| PUT | `/api/v1/records/immovable/:id/submit` | Registrar | Submit for review |

List query params: `status`, `woreda`, `search`, `page`, `limit`, `date_from`, `date_to`

## Movable Records

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/records/movable` | Registrar | Create draft (auto-generates `ET-HR-AN-V-YYYY-NNNN`) |
| GET | `/api/v1/records/movable` | All roles | List records (role-filtered, paginated) |
| GET | `/api/v1/records/movable/:id` | All roles | Get record detail |
| PUT | `/api/v1/records/movable/:id` | Registrar | Update own draft/returned record |
| PUT | `/api/v1/records/movable/:id/submit` | Registrar | Submit for review |

## Photos

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/records/:type/:id/photos` | Registrar | Upload photo (`photo` field, JPG/PNG, max 5MB, max 10 per record) |
| DELETE | `/api/v1/records/:type/:id/photos/:photo_id` | Registrar | Delete photo from own draft/returned record |

`:type` must be `immovable` or `movable`.

## Approval Workflow

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| PUT | `/api/v1/records/:type/:id/review-approve` | Supervisor | Move `pending_review` → `under_review` (optional comment) |
| PUT | `/api/v1/records/:type/:id/review-return` | Supervisor | Move `pending_review` → `returned` (comment required) |
| PUT | `/api/v1/records/:type/:id/final-approve` | Manager | Move `under_review` → `approved` (optional comment) |
| PUT | `/api/v1/records/:type/:id/final-return` | Manager | Move `under_review` → `pending_review` (comment required) |
| POST | `/api/v1/records/:type/:id/comments` | Supervisor, Manager | Add a comment |
| GET | `/api/v1/records/:type/:id/comments` | All roles | List comments (registrar: own records only) |
| GET | `/api/v1/records/:type/:id/history` | All roles | List status history (registrar: own records only) |

Return actions return **422** if `comment` is empty. Wrong status returns **409**.

## Status Audit Log

Every status change (submit, review, approve, return) writes a row to `status_history` with who changed it, from/to status, optional note, and timestamp.

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/records/:type/:id/history` | All roles | Full audit trail with `changed_by_name` (registrar: own records only) |

Record detail (`GET .../immovable/:id` and `GET .../movable/:id`) also includes a `history` array with the same enriched entries.

## Dashboard & Search

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/dashboard/stats` | All roles | Role-filtered totals, status breakdown, and `pending_my_action` for supervisor/manager |
| GET | `/api/v1/records` | All roles | Unified searchable list across immovable and movable records |

**Stats response:** `total_immovable`, `total_movable`, `by_status` (draft, pending_review, under_review, returned, approved), and `pending_my_action` (supervisor: pending review count, manager: under review count).

**List query params:** `type` (`immovable` \| `movable`, optional), `status`, `woreda`, `kebele`, `search`, `page`, `limit`, `date_from`, `date_to` (YYYY-MM-DD). Registrars only see their own records.

## Export

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/v1/export/records/csv` | Supervisor, Manager | Download filtered record list as CSV (same query params as `/records`, max 10,000 rows) |
| GET | `/api/v1/records/:type/:id/pdf` | Supervisor, Manager, Registrar | Download printable PDF for a single record |

**PDF rules:** Draft records cannot be printed. Registrars may print only their own **approved** records. Approved PDFs include an `APPROVED` watermark and the first photo when available. CSV exports metadata only (no photo file data).

## Project Structure

```
/cmd/server/main.go          Application entry point
/internal/config/            Environment configuration
/internal/db/                PostgreSQL connection + migrations
/internal/auth/               Login, refresh, logout
/internal/users/              User management CRUD
/internal/immovable/          Form 02 immovable record CRUD
/internal/movable/            Form 01 movable record CRUD
/internal/photos/             Photo upload and storage
/internal/workflow/           Approval workflow, comments
/internal/audit/              Status history writes and reads
/internal/dashboard/          Dashboard stats and unified record search
/internal/export/             CSV and PDF export
/internal/models/             Domain structs
/internal/middleware/        CORS, logging, JWT auth, role guards
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `PORT` | No | HTTP port (default: `8080`) |
| `DB_URL` | Yes | PostgreSQL connection string |
| `JWT_SECRET` | Yes | Access token signing secret |
| `JWT_REFRESH_SECRET` | Yes | Refresh token signing secret |
| `MEDIA_PATH` | No | Photo upload directory (default: `./media`) |
| `ALLOWED_ORIGINS` | No | Comma-separated CORS origins |

## Build

```bash
go build -o bin/server ./cmd/server
```

## Next Steps

- **F-01**: React frontend scaffold + routing
- **F-02**: Auth pages and JWT storage

See `CURSOR.md` and `SYSTEM_DESIGN.md` for full architecture rules.
