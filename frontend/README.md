# Qirs Mezgeb Frontend

React + TanStack Start app for the Harari Heritage Registry System (ቅርስ መዝገብ).

## Prerequisites

- [Bun](https://bun.sh/) (or Node.js 20+ with npm/pnpm)
- Backend API running at `http://localhost:8080` (see root [README.md](../README.md))

## Setup

```bash
cd frontend
cp .env.example .env
bun install
bun run dev
```

Open **http://localhost:5173** (frontend). The API runs separately on **http://localhost:8080** — see root [README.md](../README.md) § “Run locally”.

### Environment

Copy `.env.example` → `.env`:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

Do not change the frontend dev port to 8080 — that port is reserved for the Go API.

## Scripts

| Command | Description |
|---------|-------------|
| `bun run dev` | Start dev server (Vite) |
| `bun run build` | Production build |
| `bun run lint` | ESLint |
| `bun run format` | Prettier |

## Folder structure

```
frontend/src/
  api/              HTTP clients — one file per backend domain (no React)
  components/
    ui/             shadcn/ui primitives (generated; minimal manual edits)
    common/         Shared app UI (StatusBadge, LoadingSpinner, EmptyState, …)
    layout/         AppLayout, Sidebar, Topbar
    forms/          Field primitives, ImmovableForm, MovableForm, *Options.ts
    records/        RecordCard, RecordsList, RecordFilters
    dashboard/      Role-specific dashboard widgets
    workflow/       StatusTimeline, CommentThread, ReturnModal (as added)
  routes/           TanStack Router file-based routes (= pages)
  hooks/            Shared React hooks
  stores/           Zustand (auth, language)
  i18n/             am.json, en.json
  lib/              Utils, server config, error helpers
  types/            TypeScript types aligned with Go backend
```

## Routing

This project uses **TanStack Start file-based routing**. Route files live in `src/routes/`.

- Do **not** create `src/pages/` — that is a Next.js / Remix convention.
- The app shell is `src/routes/__root.tsx`; authenticated layout is `src/routes/_authenticated.tsx`.
- `routeTree.gen.ts` is auto-generated — do not edit by hand.

See [src/routes/README.md](src/routes/README.md) for naming conventions.

## Conventions

| Layer | Responsibility |
|-------|----------------|
| `routes/` | Routing, data loading, page composition |
| `api/` | Axios calls via `api/client.ts` only |
| `components/` | Presentational and domain UI |
| `stores/` | Client-only global state |
| `types/` | Shared API contracts |

- All user-visible strings go through `react-i18next` (`t('key')`).
- Data fetching uses TanStack Query; mutations invalidate related query keys.
- Component folders use **lowercase** names (`common/`, `forms/`). Only `ui/` holds shadcn.

## Tech stack

- React 19 + Vite 7 + TanStack Start / Router
- TanStack Query v5 + Axios
- Zustand, react-i18next, Tailwind CSS 4, shadcn/ui
