# API Contract — Qirs Mezgeb

Canonical request/response shapes for the Go backend. Use this as the frontend source of truth (F-01+).

All routes are prefixed with `/api/v1`. Success responses use:

```json
{ "success": true, "data": { ... }, "message": "..." }
```

Errors use:

```json
{ "success": false, "error": "...", "code": 401 }
```

Paginated lists use `items` (not `records`):

```json
{ "items": [...], "total": 0, "page": 1, "limit": 20, "total_pages": 0 }
```

---

## User object shapes

All user-facing JSON uses **`full_name`** (never `name`).

| Context | Fields |
|---------|--------|
| Login / token user (`UserPublic`) | `id`, `full_name`, `email`, `role`, `language` |
| Profile / user admin (`UserItem`) | above + `is_active`, `created_at`, `updated_at` |

---

## Auth

### POST `/auth/login`

Request: `{ "email", "password" }`

Response `data`:

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "user": {
    "id": "uuid",
    "full_name": "Admin User",
    "email": "admin@qirsmezgeb.gov.et",
    "role": "manager",
    "language": "am"
  }
}
```

### POST `/auth/refresh`

Request: `{ "refresh_token" }` → `data`: `{ "access_token" }`

### POST `/auth/logout`

Request: `{ "refresh_token" }`

---

## Users

Create/update request bodies use `full_name`.  
`GET /users/me` → `data`: `{ "user": UserItem }`

---

## Record detail

`GET /records/immovable/:id` and `GET /records/movable/:id` return:

```json
{
  "record": { ...full record fields... },
  "photos": [ RecordPhoto ],
  "comments": [ RecordComment ],
  "history": [ StatusHistoryEntry ]
}
```

`RecordComment` includes `author_name`.  
`StatusHistoryEntry` includes `changed_by_name`.

Record update (`PUT /records/:type/:id`) returns `data`: `{ "record": { ... } }`.

Submit and resubmit use the same endpoint: `PUT /records/:type/:id/submit` (works for `draft` and `returned`).

---

## Workflow comment bodies

All comment-bearing write endpoints accept **`comment_text`**. Legacy **`comment`** is still accepted on approve/return actions for backward compatibility.

| Endpoint | Body |
|----------|------|
| `PUT .../review-approve` | `{ "comment_text"?: string }` |
| `PUT .../review-return` | `{ "comment_text": string }` (required) |
| `PUT .../final-approve` | `{ "comment_text"?: string }` |
| `PUT .../final-return` | `{ "comment_text": string }` (required) |
| `POST .../comments` | `{ "comment_text": string }` (required) |

Return with empty comment → **422**.

---

## Unified record list

`GET /records` — cross-type search (dashboard).

Query: `type`, `status`, `woreda`, `kebele`, `search`, `page`, `limit`, `date_from`, `date_to`

Response `data` is paginated `RecordSummary` items:

| Field | Type |
|-------|------|
| `id` | uuid |
| `record_type` | `immovable` \| `movable` |
| `record_id` | string |
| `name_amharic` | string |
| `status` | enum |
| `woreda`, `kebele` | string? |
| `registrar_id` | uuid |
| `created_at`, `updated_at` | timestamp |

---

## Dashboard stats

`GET /dashboard/stats` → `data`:

```json
{
  "total_immovable": 0,
  "total_movable": 0,
  "by_status": {
    "draft": 0,
    "pending_review": 0,
    "under_review": 0,
    "returned": 0,
    "approved": 0
  },
  "pending_my_action": 0
}
```

`pending_my_action` is omitted for registrars. Supervisor = pending review count; manager = under review count.

---

## Photos

Upload: `multipart/form-data`, field name **`photo`** (not JSON).

---

## Export

| Endpoint | Auth | Notes |
|----------|------|-------|
| `GET /export/records/csv` | Supervisor, Manager | Same filters as `/records`, max 10k rows |
| `GET /records/:type/:id/pdf` | All roles (service enforces rules) | Registrar: own approved only; draft forbidden |

---

## Movable submit validation

At submit time, movable records require `woreda` and `kebele` (same as immovable). Draft saves may leave them null.
