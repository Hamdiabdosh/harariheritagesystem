# Qirs Mezgeb — Backend Audit Report
## Step 1: What the code actually is | Step 2: Deviations from SYSTEM_DESIGN.md

---

## ✅ What Was Built (Confirmed by Code)

The backend is complete, well-structured, and follows the Handler → Service → Repository pattern correctly. Here is what actually exists:

### Modules
- `auth` — login, refresh, logout, JWT (15min access / 7day refresh)
- `users` — CRUD + deactivate (manager only)
- `immovable` — full Form 02 record lifecycle
- `movable` — full Form 01 record lifecycle
- `workflow` — status transitions + comments + history
- `photos` — upload/delete per record
- `dashboard` — stats + unified record list with filters
- `export` — CSV export + PDF export
- `audit` — internal audit trail (not directly exposed as its own API route)

---

## 🔴 Deviations from SYSTEM_DESIGN.md

### 1. `UserPublic` has `name` not `full_name`
**Design said:** `full_name`  
**Actual code (`models/user.go`):**
```go
type UserPublic struct {
    ID       uuid.UUID `json:"id"`
    FullName string    `json:"name"`   // ← serialized as "name", NOT "full_name"
    Email    string    `json:"email"`
    Role     Role      `json:"role"`
    Language Language  `json:"language"`
}
```
**Impact:** Login response returns `user.name`, not `user.full_name`. Frontend must use `name`.

### 2. `UserItem` (from users list) has `full_name`
**Actual code (`users/types.go`):**
```go
type UserItem struct {
    FullName  string    `json:"full_name"`  // ← "full_name" here
    ...
}
```
**Impact:** The users list uses `full_name` but the login response user object uses `name`. Inconsistency within backend — frontend must handle both.

### 3. Record detail response shape
**Design said:** A single record object  
**Actual code (`immovable/types.go`):**
```go
type RecordDetail struct {
    Record   models.ImmovableRecord `json:"record"`
    Photos   []any                  `json:"photos"`
    Comments []any                  `json:"comments"`
    History  []any                  `json:"history"`
}
```
**Impact:** `GET /records/immovable/:id` returns `{ record, photos, comments, history }` — not just the record. Frontend must access `data.record`, `data.photos`, etc.

### 4. Update record response shape
**Design said:** Updated record object  
**Actual handler (`immovable/handler.go`):**
```go
middleware.RespondSuccess(c, gin.H{"record": record}, "Record updated successfully")
```
**Impact:** Update returns `{ data: { record: {...} } }`, not `{ data: { ...record fields... } }`.

### 5. Unified `/api/v1/records` endpoint exists (not in original design)
**Actual:** `GET /api/v1/records` (dashboard module) — returns `RecordSummary` list across both types  
**Design:** Not explicitly specified — the design only specified type-specific list endpoints  
**Impact:** This is actually an improvement. Use this for dashboards. `RecordSummary` has fewer fields than the full record.

### 6. `MovableRecord` — `woreda` and `kebele` are nullable
**Design:** Same as immovable (required fields)  
**Actual (`models/movable.go`):** `Woreda *string`, `Kebele *string` — both nullable pointers  
**Impact:** Movable form doesn't enforce woreda/kebele at model level.

### 7. `MovableRecord` has no `HasThreat` boolean — it has `ThreatDescription`
**Actual fields:**
- `HasThreat *bool` ✓ (exists)
- `ThreatDescription *string` ✓ (added, not in design)
- No `MaintenanceReason` field (immovable has it, movable doesn't)

### 8. `MovableRecord` has extra informant field: `InformantOccupation`
**Actual:** `InformantOccupation *string` — not in design  
**Impact:** Extra field, frontend should include it.

### 9. `ImmovableRecord` has `RelatedDocs []string` and `HasOralHistory *bool`
Not in the original design spec but present in the model.

### 10. Workflow — `ReviewApprove` and `FinalApprove` accept an *optional* comment
**Design said:** Approve never requires a comment  
**Actual:** `optionalCommentRequest` struct — comment sent but not required for approvals. ✅ Matches design.

### 11. Comment endpoint body field is `comment_text` not `comment`
**Actual (`workflow/handler.go`):**
```go
type addCommentRequest struct {
    CommentText string `json:"comment_text" binding:"required"`
}
```
**Impact:** `POST /records/:type/:id/comments` body must be `{ "comment_text": "..." }`.

### 12. Workflow return actions use `comment` field (not `comment_text`)
**Actual:**
```go
type requiredCommentRequest struct {
    Comment string `json:"comment" binding:"required"`
}
```
**Impact:** `PUT /records/:type/:id/review-return` and `final-return` body must be `{ "comment": "..." }`.

### 13. Photo upload uses `multipart/form-data`, field name is `photo`
**Actual (`photos/handler.go`):** `c.FormFile("photo")`  
**Impact:** Photo upload is NOT JSON — it's a multipart form with field name `photo`.

### 14. Export CSV is supervisor+manager only (as designed), but PDF is available to ALL roles
**Actual routes:**
```
supervisorManager.GET("/export/records/csv", ...)   // supervisor + manager only
authorized.GET("/records/:type/:id/pdf", ...)        // ALL authenticated roles
```
**Design said:** "own approved" for registrar, all for supervisor/manager.  
**Impact:** The registrar can call the PDF endpoint — the service layer likely enforces the approved-only rule.

### 15. No `/records/:type/:id/resubmit` endpoint
**Design mentioned resubmit as a workflow step.**  
**Actual:** There is no explicit resubmit endpoint. A returned record is resubmitted via `PUT /records/:type/:id/submit` (the same submit endpoint). The service layer checks if status is `returned` and transitions to `pending_review`.  
**Impact:** Frontend should use the same submit button/endpoint for both first submission and resubmission.

### 16. `dashboard/stats` response shape
**Actual (`dashboard/types.go`):**
```go
type Stats struct {
    TotalImmovable  int          `json:"total_immovable"`
    TotalMovable    int          `json:"total_movable"`
    ByStatus        StatusCounts `json:"by_status"`
    PendingMyAction *int         `json:"pending_my_action,omitempty"`
}
type StatusCounts struct {
    Draft         int `json:"draft"`
    PendingReview int `json:"pending_review"`
    UnderReview   int `json:"under_review"`
    Returned      int `json:"returned"`
    Approved      int `json:"approved"`
}
```

### 17. `RecordSummary` (unified list) shape
```go
type RecordSummary struct {
    ID          uuid.UUID
    RecordType  RecordType      // "immovable" or "movable"
    RecordID    string          // e.g. "ET-HR-AN-I-2024-0001"
    NameAmharic string
    Status      RecordStatus
    Woreda      *string
    Kebele      *string
    RegistrarID uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

---

## ✅ Things That Match Design Exactly
- JWT 15-min access + 7-day refresh
- Status enum: `draft`, `pending_review`, `under_review`, `returned`, `approved`
- Role enum: `registrar`, `supervisor`, `manager`
- Language enum: `am`, `en`
- Record ID format: `ET-HR-AN-I-YYYY-NNNN` (immovable), `ET-HR-AN-V-YYYY-NNNN` (movable)
- All API routes prefixed `/api/v1/`
- Standard response envelope: `{ success, data, message }` / `{ success, error, code }`
- Pagination: `{ items, total, page, limit, total_pages }`
- Return actions enforce non-empty comment (422 if empty)
- Approved records locked (no edit)
- Registrar filtered to own records only
- Soft delete (deactivate) for users
- Photo: max 10 per record, max 5MB each, JPG/PNG only
- `GET /health` endpoint available at both root and `/api/v1/health`

---

## Resolved (Backend Consistency Fix)

The following audit findings from Step 1 have been addressed. See [API.md](API.md) for the canonical contract.

| # | Issue | Resolution |
|---|-------|------------|
| 1–2 | Login `user.name` vs list `full_name` | All user JSON uses `full_name`; `UserItem` embeds `UserPublic` |
| 3–4 | Record detail / update shapes | Documented in API.md; comments now populated in detail |
| 5 | Unified `GET /records` | Documented in SYSTEM_DESIGN + API.md |
| 6 | Movable woreda/kebele not validated on submit | Added to `validateForSubmit` |
| 11–12 | Mixed `comment` / `comment_text` | Standardized on `comment_text`; legacy `comment` shim on workflow writes |
| 14 | PDF route authorization | Service-layer checks + explicit route group; tests extended |
| 15 | No `/resubmit` endpoint | Documented: use same `PUT .../submit` for returned records |

Items 7–10 (extra model fields) and 16–17 (stats/summary shapes) are documented as intentional in API.md, not bugs.
