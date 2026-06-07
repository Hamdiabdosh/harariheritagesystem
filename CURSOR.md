# CURSOR.md — Qirs Mezgeb (ቅርስ መዝገብ)
# Heritage Registry System — AI Agent Rules File
# Read this before every code generation task in this project.

---

## Project Overview

Qirs Mezgeb is a government heritage registration management system for the Harari Regional
State Culture, Heritage and Tourism Bureau in Ethiopia. It digitizes the paper-based process
of registering immovable heritage assets (Form 02) and movable heritage assets (Form 01),
routing each record through a 3-stage approval chain: Registrar → Supervisor → Manager.

The system runs on an Ubuntu Linux government server. The frontend is a React PWA accessible
from any phone browser. The backend is a Go REST API. All UI text is bilingual (Amharic + English).

---

## Tech Stack

| Layer       | Technology                                         |
|-------------|----------------------------------------------------|
| Frontend    | React 19 + Vite 7 + TanStack Start, TypeScript, Tailwind CSS |
| Routing     | TanStack Router (file-based routes in `frontend/src/routes/`) |
| State       | Zustand                                            |
| Data fetch  | TanStack Query v5 + Axios                          |
| Forms       | React Hook Form + Zod                              |
| i18n        | i18next (am.json + en.json)                        |
| Backend     | Go 1.22 + Gin                                      |
| Database    | PostgreSQL 16                                      |
| Migrations  | golang-migrate (raw SQL, no ORM)                   |
| Auth        | JWT (golang-jwt/jwt/v5) + bcrypt                   |
| IDs         | UUID (google/uuid)                                 |
| Server      | Nginx (reverse proxy) + systemd                    |

---

## Folder Structure

### Backend
```
/cmd/server/main.go
/internal/
  auth/           handler.go, service.go, repository.go
  users/          handler.go, service.go, repository.go
  immovable/      handler.go, service.go, repository.go
  movable/        handler.go, service.go, repository.go
  workflow/       handler.go, service.go
  photos/         handler.go, service.go
  dashboard/      handler.go, service.go
  export/         handler.go, service.go
  middleware/     auth.go, role.go, logger.go, error.go
  db/             connection.go, migrations/
  models/         user.go, immovable.go, movable.go, photo.go,
                  comment.go, history.go
  config/         config.go
/media/           (uploaded photos — gitignored)
.env.example
go.mod
CURSOR.md
```

### Frontend (`frontend/src/`)
```
api/              client.ts, auth.ts, users.ts, records.ts, immovable.ts,
                  movable.ts, workflow.ts, dashboard.ts, export.ts
components/
  ui/             shadcn/ui primitives (generated)
  common/         StatusBadge, LoadingSpinner, EmptyState, StatusTimeline,
                  CommentThread, ReturnModal, …
  layout/         AppLayout.tsx, Sidebar.tsx, Topbar.tsx
  forms/          Field.tsx, ImmovableForm, MovableForm, PhotoUploader, *Options.ts
  records/        RecordCard, RecordsList, RecordFilters
  dashboard/      RegistrarDashboard, StatCard, …
  workflow/       Review actions, modals (as added)
routes/           TanStack Router file-based routes (= pages)
  __root.tsx      app shell
  login.tsx
  _authenticated.tsx          layout + auth guard
  _authenticated/             dashboard, records, pending, reviewed, users, reports
stores/           authStore.ts, languageStore.ts
hooks/            use-mobile.tsx, …
i18n/             index.ts, am.json, en.json
lib/              utils.ts, config.server.ts, error helpers
types/            index.ts
```

Do **not** create `src/pages/` — routes live under `src/routes/`. See `frontend/README.md`.

---

## Architecture Rules

### Backend — STRICT Handler → Service → Repository pattern

```
HTTP Request
    ↓
Handler     — Parse request, call service, return response. NO business logic.
    ↓
Service     — All business logic, validation, rules. Calls repository.
    ↓
Repository  — Database queries ONLY. No business logic.
```

**Handler must only:**
- Parse and bind request body/params
- Call one service method
- Return JSON response

**Service must:**
- Validate inputs
- Enforce business rules
- Orchestrate multiple repository calls if needed
- Return domain errors (not HTTP errors)

**Repository must only:**
- Execute SQL queries
- Return domain models or errors

### Frontend — Pages → Hooks → API

```
Page component
    ↓
Custom hook (useXxx)    — TanStack Query, state, mutations
    ↓
API module (src/api/)   — Axios calls
    ↓
Go backend
```

---

## API Response Format — MANDATORY on every endpoint

```json
// Success
{ "success": true, "data": { ... }, "message": "Record created" }

// Error
{ "success": false, "error": "Validation failed", "code": 422 }

// Paginated list
{
  "success": true,
  "data": {
    "items": [...],
    "total": 150,
    "page": 1,
    "limit": 20,
    "total_pages": 8
  }
}
```

All routes are prefixed: `/api/v1/`

---

## Status Transition Table

| From            | To              | Action                | Who        |
|-----------------|-----------------|----------------------|------------|
| draft           | pending_review  | submit               | Registrar  |
| pending_review  | under_review    | review-approve       | Supervisor |
| pending_review  | returned        | review-return*       | Supervisor |
| under_review    | approved        | final-approve        | Manager    |
| under_review    | pending_review  | final-return*        | Manager    |
| returned        | pending_review  | resubmit             | Registrar  |
| approved        | (locked)        | —                    | Nobody     |

*Return actions ALWAYS require a non-empty comment. Enforce in service layer.
Approved records cannot be modified by anyone under any circumstance.

---

## Role Permission Matrix

| Action                          | Registrar | Supervisor | Manager |
|---------------------------------|-----------|------------|---------|
| Create record                   | ✓ own     | ✗          | ✗       |
| Edit draft / returned record    | ✓ own     | ✗          | ✗       |
| Submit record                   | ✓ own     | ✗          | ✗       |
| View own records                | ✓         | ✓ all      | ✓ all   |
| View all records                | ✗         | ✓          | ✓       |
| Review-approve / review-return  | ✗         | ✓          | ✗       |
| Final-approve / final-return    | ✗         | ✗          | ✓       |
| Add comment                     | ✗         | ✓          | ✓       |
| Create / manage users           | ✗         | ✗          | ✓       |
| Export CSV / PDF                | own appr. | ✓          | ✓       |
| View dashboard stats            | own only  | ✓          | ✓ full  |

---

## Naming Conventions

### Go (Backend)
- Files: `snake_case.go`
- Types/Structs: `PascalCase`
- Functions/Methods: `PascalCase` (exported), `camelCase` (unexported)
- Constants: `PascalCase` or `SCREAMING_SNAKE` for enums
- DB column names in queries: `snake_case`
- Error variables: `ErrXxx` (e.g. `ErrRecordNotFound`)

### TypeScript (Frontend)
- Component files: `PascalCase.tsx`
- Hook files: `useCamelCase.ts`
- API files: `camelCase.ts`
- Store files: `camelCaseStore.ts`
- Type/Interface names: `PascalCase`
- CSS classes: Tailwind utility classes only (no custom CSS unless unavoidable)
- i18n keys: `snake_case` nested by module (e.g. `form.section1.name_amharic`)

---

## What NOT To Do

```
✗ NO business logic in HTTP handlers — move to service layer
✗ NO hard deletes on any heritage records (immovable or movable)
✗ NO hardcoded UI text — every string must be in am.json and en.json
✗ NO mock/fake data in production code
✗ NO skipping validation — validate on both client (Zod) and server (Go)
✗ NO direct DB calls from handlers — always go through service → repository
✗ NO allowing registrar to query another registrar's records
   (always filter by registrar_id WHERE role = 'registrar')
✗ NO allowing status transitions outside the table above
✗ NO allowing "Return" action with empty comment
✗ NO editing of records with status = 'approved'
✗ NO exposing password_hash in any API response
✗ NO storing raw JWT refresh tokens — store bcrypt hash only
✗ NO missing loading/error/empty states in React components
```

---

## How to Add a New API Endpoint (Checklist)

1. [ ] Add SQL query to the relevant `repository.go`
2. [ ] Add business logic method to `service.go` — call the repository
3. [ ] Add handler function to `handler.go` — parse request, call service, return JSON
4. [ ] Register route in the router setup in `main.go` with correct middleware
5. [ ] Add role guard: `middleware.RequireRole("supervisor")` etc.
6. [ ] Add corresponding function to `/src/api/[module].ts`
7. [ ] Add or update TanStack Query hook in `/src/hooks/`
8. [ ] Update TypeScript types in `/src/types/index.ts` if new shapes returned

---

## How to Add a New Form Field (Checklist)

1. [ ] Write migration: `ALTER TABLE xxx ADD COLUMN yyy TYPE;`
2. [ ] Update Go model struct in `/internal/models/`
3. [ ] Update repository INSERT and UPDATE queries
4. [ ] Update TypeScript interface in `/src/types/index.ts`
5. [ ] Add field to the correct form component in `frontend/src/components/forms/`
6. [ ] Add Zod validation rule in the form schema
7. [ ] Add i18n keys in `am.json` and `en.json`
8. [ ] Add field to the read-only record detail view

---

## Environment Variables

### Backend (.env)
```
PORT=8080
DB_URL=postgres://user:password@localhost:5432/qirsmezgeb?sslmode=disable
JWT_SECRET=<long-random-string>
JWT_REFRESH_SECRET=<different-long-random-string>
MEDIA_PATH=/var/qirsmezgeb/media
ALLOWED_ORIGINS=http://localhost:5173,https://heritage.harari.gov.et
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080/api/v1
```

---

## Common Patterns

### Role guard middleware (Go)
```go
// In route registration
authorized := r.Group("/api/v1")
authorized.Use(middleware.AuthRequired())
{
    supervisorOnly := authorized.Group("/")
    supervisorOnly.Use(middleware.RequireRole("supervisor", "manager"))
    supervisorOnly.PUT("/records/:type/:id/review-approve", workflow.Handler.ReviewApprove)
}
```

### Standard paginated query (Go repository)
```go
func (r *Repository) List(ctx context.Context, filters ListFilters) ([]ImmovableRecord, int, error) {
    query := `SELECT *, COUNT(*) OVER() as total FROM immovable_records WHERE 1=1`
    args := []interface{}{}
    argIdx := 1

    if filters.Status != "" {
        query += fmt.Sprintf(" AND status = $%d", argIdx)
        args = append(args, filters.Status)
        argIdx++
    }
    if filters.RegistrarID != uuid.Nil { // enforce for registrar role
        query += fmt.Sprintf(" AND registrar_id = $%d", argIdx)
        args = append(args, filters.RegistrarID)
        argIdx++
    }
    query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
    args = append(args, filters.Limit, (filters.Page-1)*filters.Limit)
    // ... execute and scan
}
```

### Standard error response (Go handler)
```go
func respondError(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{"success": false, "error": message, "code": code})
}

func respondSuccess(c *gin.Context, data interface{}, message string) {
    c.JSON(http.StatusOK, gin.H{"success": true, "data": data, "message": message})
}
```

### Zustand store (TypeScript)
```typescript
// src/stores/authStore.ts
import { create } from 'zustand'

interface AuthState {
  user: User | null
  accessToken: string | null
  isAuthenticated: boolean
  login: (user: User, token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  accessToken: null,
  isAuthenticated: false,
  login: (user, accessToken) => set({ user, accessToken, isAuthenticated: true }),
  logout: () => set({ user: null, accessToken: null, isAuthenticated: false }),
}))
```

### TanStack Query hook (TypeScript)
```typescript
// src/hooks/useImmovableRecords.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getImmovableRecords, submitImmovableRecord } from '@/api/immovable'

export function useImmovableRecords(filters: RecordFilters) {
  return useQuery({
    queryKey: ['immovable-records', filters],
    queryFn: () => getImmovableRecords(filters),
  })
}

export function useSubmitRecord() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => submitImmovableRecord(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['immovable-records'] })
    },
  })
}
```

### Authenticated layout (TanStack Router)
```typescript
// frontend/src/routes/_authenticated.tsx
// Wraps all logged-in routes with AppLayout and redirects to /login when unauthenticated.
// Per-route role checks use <Navigate to="/unauthorized" /> in the route component.
```

---

## Record ID Format

| Type       | Format                    | Example                |
|------------|---------------------------|------------------------|
| Immovable  | ET-HR-AN-I-YYYY-NNNN      | ET-HR-AN-I-2024-0001   |
| Movable    | ET-HR-AN-V-YYYY-NNNN      | ET-HR-AN-V-2024-0042   |

NNNN = sequential count of that type in that calendar year, zero-padded to 4 digits.
Generated on first INSERT (even as draft). Never changes after generation.

---

## Bilingual Text Rule

Every visible string in the UI must exist in both:
- `/src/i18n/am.json` — Amharic (primary)
- `/src/i18n/en.json` — English

Usage in components:
```typescript
import { useTranslation } from 'react-i18next'
const { t } = useTranslation()
// <label>{t('form.section1.name_amharic')}</label>
```

Key structure:
```json
{
  "nav": { "dashboard": "...", "records": "...", "users": "..." },
  "auth": { "login": "...", "logout": "...", "email": "..." },
  "form": {
    "section1": { "name_amharic": "...", "name_local": "...", "category": "..." },
    "section2": { ... }
  },
  "status": { "draft": "...", "pending_review": "...", "approved": "..." },
  "actions": { "save_draft": "...", "submit": "...", "approve": "...", "return": "..." }
}
```
