# ቅርስ መዝገብ — Qirs Mezgeb
## Heritage Registry System — Full System Analysis & Design
### Harari Region Culture, Heritage & Tourism Bureau

---

## PHASE 1 — PROJECT VISION

```
PROJECT VISION
==============
Project Name:    Qirs Mezgeb (ቅርስ መዝገብ) — Heritage Registry System
Type:            Government Management System (Web + Mobile PWA)
Region:          Harari Regional State, Ethiopia
Core Problem:    Heritage registration is done manually on paper forms,
                 making records impossible to search, easy to lose,
                 and impossible to track through an approval chain.
Who Suffers:     Field registrars, supervisors, managers, and the
                 institution's long-term archival integrity.
Why Current
Solutions Fail:  Paper forms have no approval workflow, no status
                 tracking, no photo attachment, no duplicate prevention,
                 no search, and no protection against physical loss.
What Improves:   A structured digital workflow from field entry to final
                 approval, with a searchable, permanent, bilingual record
                 for every heritage asset.
Success Looks
Like:            Registrar fills form on phone in the field → supervisor
                 reviews same day → manager approves with one click →
                 record locked with unique ID, searchable by all.
```

---

## PHASE 2 — SYSTEM ACTORS & PERMISSIONS

```
SYSTEM ACTORS & PERMISSIONS
============================

Actor: Registrar
  Can View:    Own drafts, own submitted records, own returned records + comments
  Can Create:  Immovable heritage record (Form 02), Movable heritage record (Form 01)
  Can Update:  Own drafts (anytime), own returned records (only after returned by supervisor)
  Can Delete:  Own drafts only — never submitted or approved records
  Depends On:  Supervisor (review), Manager (final approval)
  Workflow:    Login → choose form type → fill form → save draft or submit
               → get notification if returned → fix and resubmit

Actor: Supervisor
  Can View:    All submitted records (any registrar), all reviewed records, own comments
  Can Create:  Review comment on any record
  Can Update:  Record status → "Reviewed" (pass to manager) or "Returned" (back to registrar)
  Can Delete:  Nothing
  Depends On:  Registrar submits, Manager decides final
  Workflow:    Login → see pending queue → open record → review all fields
               → approve to manager or return with written comments

Actor: Manager
  Can View:    All records at any status, full dashboard, all users, reports
  Can Create:  New user accounts with role assignment
  Can Update:  Final record status → "Approved" (locked), user roles
  Can Delete:  User accounts (soft delete / deactivate only)
  Depends On:  Supervisor must pre-review before records reach manager
  Workflow:    Login → see supervisor-reviewed records → give final approval
               → manage users → view reports → export records
```

---

## PHASE 3 — MODULE BREAKDOWN

```
MODULE BREAKDOWN
================

Module 1: Auth & User Management
  Primary Actor(s): Manager (admin), all actors (login)
  Core Responsibility: Secure login, role assignment, JWT session
  Sub-features:
    - Login / logout
    - JWT token issuance and refresh
    - Role-based route guards (frontend + backend)
    - Manager: create / edit / deactivate users
    - Language preference per user (Amharic / English)
    - Password reset (admin-initiated for now)
  Depends On: Nothing — built first
  Priority: 1 — must-have

Module 2: Immovable Heritage Registration (Form 02)
  Primary Actor(s): Registrar
  Core Responsibility: Digital Form 02 for immovable heritage assets
  Sub-features:
    - All 8 form sections (identification, ownership, GPS, history,
      measurements, Harari house classification, condition, conservation)
    - Photo upload (multiple images per record)
    - GPS coordinate capture from device
    - Save as draft (auto-save every 2 minutes)
    - Submit for review
    - Auto-generate unique ID: ET-HR-AN-I-[YYYY]-[NNNN]
    - View own record list
  Depends On: Auth module
  Priority: 1 — must-have

Module 3: Movable Heritage Registration (Form 01)
  Primary Actor(s): Registrar
  Core Responsibility: Digital Form 01 for movable heritage assets
  Sub-features:
    - All form sections (identification, ownership, storage location,
      acquisition method, measurements, material, condition, conservation)
    - Photo upload (multiple images per record)
    - Save as draft
    - Submit for review
    - Auto-generate unique ID: ET-HR-AN-V-[YYYY]-[NNNN]
    - View own record list
  Depends On: Auth module
  Priority: 1 — must-have

Module 4: Approval Workflow
  Primary Actor(s): Supervisor (review), Manager (final approval)
  Core Responsibility: 3-stage status pipeline with comments and audit trail
  Sub-features:
    - Pending queue per role (supervisor sees submitted, manager sees reviewed)
    - Approve (pass to next stage) or Return (back to registrar)
    - Comment thread per record (with author + timestamp)
    - Full status history / audit log per record
    - In-app notification badge (new items in queue)
  Depends On: Both registration modules
  Priority: 1 — must-have

Module 5: Dashboard & Search
  Primary Actor(s): All actors (filtered by role)
  Core Responsibility: Overview statistics, record search, filter
  Sub-features:
    - Record counts by status / type / date (manager sees all, others see own)
    - Search by name, ID, location, date range
    - Filter by form type / status / woreda / kebele
    - Paginated record list
  Depends On: Registration modules + Approval workflow
  Priority: 2 — important

Module 6: Reports & Export
  Primary Actor(s): Manager, Supervisor
  Core Responsibility: Printable records and data exports
  Sub-features:
    - Print single record as formatted PDF (matches original paper form layout)
    - Export filtered record list to CSV
    - Summary statistics report (total records by type, status, region)
  Depends On: Dashboard module
  Priority: 2 — important

BUILD ORDER: Auth → Immovable Form → Movable Form → Approval Workflow → Dashboard → Reports
```

---

## PHASE 4 — DEEP MODULE ANALYSIS

---

### MODULE SPEC: Auth & User Management

```
MODULE SPEC: Auth & User Management
=====================================
Purpose: Securely authenticate all users and enforce role-based access
         across every module. The manager is the only one who can create
         and manage user accounts.

FEATURES:
  1. Login: Email + password, returns JWT access token + refresh token
  2. Logout: Invalidates refresh token server-side
  3. Token Refresh: Silent refresh before expiry (15-min access, 7-day refresh)
  4. User Management: Manager creates users, assigns role, sets initial password
  5. Language Preference: Stored per user, applied on login (AM / EN)
  6. Deactivate User: Manager can deactivate account (soft delete — data kept)

USER JOURNEY (Registrar login):
  1. Open system URL on phone/desktop
  2. Enter email + password
  3. System validates credentials, issues JWT
  4. Redirected to role-specific dashboard
  5. Language loads based on saved preference
  6. Session auto-refreshes while active
  7. After 7 days of inactivity — must log in again

BUSINESS RULES:
  - Only Manager can create, edit, or deactivate users
  - A user cannot change their own role
  - A deactivated user cannot log in but their records remain intact
  - Registrar cannot see other registrars' records
  - Supervisor can see all submitted/reviewed records regardless of registrar
  - Manager can see all records at all statuses
  - Every API route must check: (1) valid JWT, (2) correct role
  - Failed login attempts: lock account after 5 consecutive failures

STATUSES:
  User can be: Active | Deactivated

EDGE CASES:
  - Expired access token while submitting form: silently refresh, resubmit
  - Deactivated user mid-session: next API call returns 401, redirect to login
  - Manager accidentally deactivates themselves: system must prevent this
  - Two sessions on same account: both valid (no single-session enforcement for now)

VALIDATION RULES:
  Field email:    Required, valid email format, unique in system
  Field password: Min 8 chars, at least one number
  Field role:     Must be one of: registrar | supervisor | manager
  Field name:     Required, Amharic or Latin characters, max 100 chars
```

---

### MODULE SPEC: Immovable Heritage Registration (Form 02)

```
MODULE SPEC: Immovable Heritage Registration
=============================================
Purpose: Allow registrars to digitally fill, save, and submit Form 02
         for immovable heritage assets with full field coverage,
         photo uploads, and GPS data.

FEATURES:
  1. Multi-section form: 8 sections, all fields from paper Form 02
  2. Auto-save draft: every 2 minutes, saved to server
  3. Photo upload: up to 10 photos per record, max 5MB each
  4. GPS capture: one click to fill GPS East/North/Elevation from device
  5. Submit for review: changes status to "Pending Review"
  6. Auto-ID generation: ET-HR-AN-I-[YYYY]-[NNNN] on first save
  7. My Records list: registrar sees own records with status badges

USER JOURNEY:
  1. Registrar logs in, taps "New Immovable Record"
  2. Section 1 loads — fills identity fields (name AM, name local, category)
  3. Fills location (woreda, kebele, house no, street, gate)
  4. Taps GPS button → device location fills coordinates automatically
  5. Fills ownership, history, measurements
  6. Selects Harari house grade (checkboxes)
  7. Fills condition section — taps condition buttons per element
  8. Fills heritage values, conservation history
  9. Uploads photos (camera or gallery)
  10. Taps "Save Draft" at any point — progress saved
  11. When complete, taps "Submit for Review"
  12. Status changes to "Pending Review" — supervisor is notified

BUSINESS RULES:
  - A draft can be edited freely until submitted
  - Once submitted, registrar CANNOT edit unless supervisor returns it
  - A returned record is editable again — registrar must resubmit
  - Record ID is generated on first save (even as draft)
  - Region field is always "ሀረሪ" — locked, not editable
  - Photos are stored linked to the record — deleting record deletes photos
  - GPS is optional (not all sites may allow device GPS)
  - Category must be selected (I, II, VII, or VIII)
  - One record = one heritage asset (no bulk submission)

STATUSES:
  Record can be: Draft | Pending Review | Under Review | Returned | Approved

EDGE CASES:
  - Network lost mid-form: auto-save draft locally (IndexedDB), sync when back online
  - GPS unavailable: fields left empty, registrar fills manually
  - Photo upload fails: show error per photo, allow retry without re-filling form
  - Duplicate location (same GPS): warn registrar, allow proceed with confirmation
  - Registrar submits incomplete form: block submission, highlight missing required fields
  - Session expires mid-form: save draft first, then redirect to login

VALIDATION RULES:
  Section 1:
    name_amharic:     Required, Amharic text, max 200 chars
    name_local:       Optional, max 200 chars
    category:         Required, at least one of [I, II, VII, VIII]
    woreda:           Required
    kebele:           Required
    current_use:      Required, at least one selection
  Section 2:
    owner_type:       Required, one of [public, government, religion, private, waqf]
  Section 3:
    construction_period: Required
    age_method:       Required, one of [estimated, exact, relative]
    height/length/width: Optional, numeric only
    doors/windows/rooms: Optional, integer only
  Photos:
    count:            0–10 files
    size:             Max 5MB per file
    format:           JPG, PNG only
```

---

### MODULE SPEC: Movable Heritage Registration (Form 01)

```
MODULE SPEC: Movable Heritage Registration
==========================================
Purpose: Allow registrars to digitally fill, save, and submit Form 01
         for movable heritage assets (objects, manuscripts, artifacts).

FEATURES:
  1. Multi-section form: all fields from paper Form 01
  2. Auto-save draft: every 2 minutes
  3. Photo upload: up to 10 photos, max 5MB each
  4. Submit for review
  5. Auto-ID: ET-HR-AN-V-[YYYY]-[NNNN]
  6. My Records list with status badges

USER JOURNEY:
  1. Registrar taps "New Movable Record"
  2. Fills identity (name AM, local name, category III–VI)
  3. Fills location (where object currently is)
  4. Selects owner type and storage location
  5. Fills history (maker, period, how acquired)
  6. Fills measurements (height, width, length, diameter, weight)
  7. Selects material(s) from checklist (gold, silver, bronze, etc.)
  8. Notes decoration, color, description
  9. Fills condition, threats, conservation history
  10. Uploads photos
  11. Saves draft or submits

BUSINESS RULES:
  - Same draft/submit/return rules as Immovable module
  - Category must be one of: III, IV, V, VI
  - Storage location is required (museum, store, church, private home, other)
  - At least one material must be selected
  - Manuscript-type records: page count and chapter count fields become required

STATUSES:
  Record can be: Draft | Pending Review | Under Review | Returned | Approved

EDGE CASES:
  - Same as Immovable module (offline, GPS, photos, duplicates)
  - Manuscript sub-type: conditionally show page/chapter fields

VALIDATION RULES:
  name_amharic:     Required
  category:         Required, one of [III, IV, V, VI]
  owner_type:       Required
  storage_location: Required
  material:         Required, at least one selected
  height/width etc: Optional, numeric
  photos:           0–10, max 5MB, JPG/PNG
```

---

### MODULE SPEC: Approval Workflow

```
MODULE SPEC: Approval Workflow
==============================
Purpose: Route submitted records through a two-stage review process
         (Supervisor → Manager) with full comment history and audit trail.

FEATURES:
  1. Pending queue: each role sees only records at their stage
  2. Approve action: moves record to next stage, notifies next actor
  3. Return action: sends record back to registrar with required comment
  4. Comment thread: all comments on a record in chronological order
  5. Status audit log: every status change recorded (who, when, from, to)
  6. Notification badge: count of items in queue shown in nav

USER JOURNEY (Supervisor):
  1. Login → notification badge shows pending count
  2. Opens pending queue → sorted by submission date (oldest first)
  3. Opens a record → sees all fields + photos
  4. Reviews data quality and completeness
  5. Option A — Approve: status becomes "Under Review", manager notified
  6. Option B — Return: must type comment, status becomes "Returned",
     registrar notified with comment visible

USER JOURNEY (Manager):
  1. Login → sees records at "Under Review" status
  2. Opens record → sees all fields + supervisor's comment history
  3. Option A — Final Approve: status becomes "Approved", record locked
  4. Option B — Return to Supervisor: sends back with comment (rare case)

BUSINESS RULES:
  - Supervisor can ONLY act on records with status "Pending Review"
  - Manager can ONLY act on records with status "Under Review"
  - Manager cannot approve a record that hasn't been reviewed by supervisor first
  - A "Return" action ALWAYS requires a written comment (cannot be empty)
  - "Approved" records are fully locked — no edits by anyone
  - Audit log is immutable — no one can delete status history
  - Notifications are in-app only (no email for now)
  - Supervisor can add a comment without changing status (for internal notes)

STATUSES:
  Draft → Pending Review → Under Review → Approved
                       ↘ Returned (to Registrar)
                                      ↘ Returned (to Supervisor, rare)

EDGE CASES:
  - Supervisor tries to approve already-approved record: blocked with message
  - Manager returns to supervisor (not registrar): supervisor gets notification,
    must re-review and re-approve to manager
  - Registrar resubmits returned record: goes back to "Pending Review" queue,
    supervisor sees it as a resubmission (flag shown)
  - Two supervisors try to act on same record simultaneously: first write wins,
    second gets "record already actioned" error
```

---

### MODULE SPEC: Dashboard & Search

```
MODULE SPEC: Dashboard & Search
================================
Purpose: Give each actor a role-appropriate overview of the system state
         with the ability to search and filter records.

FEATURES:
  1. Stats cards: total records, by status, by form type
  2. Recent activity feed: latest submissions and approvals
  3. Search bar: search by name, ID, location
  4. Filter panel: form type, status, woreda, kebele, date range
  5. Paginated record list (20 per page)
  6. Quick actions from list (open, view status, view comments)

ROLE-FILTERED VIEWS:
  Registrar:  Sees only own records. Stats show own submission counts.
  Supervisor: Sees all records. Stats show pending count prominently.
  Manager:    Sees all records. Full stats including approval rates.

BUSINESS RULES:
  - Registrar cannot see other registrars' records in any list or search
  - Search results respect role-based filtering
  - Stats are calculated in real-time (no caching needed at this scale)
```

---

### MODULE SPEC: Reports & Export

```
MODULE SPEC: Reports & Export
==============================
Purpose: Allow managers and supervisors to print individual records
         and export data for archival or reporting.

FEATURES:
  1. Print single record: generates PDF matching paper form layout,
     includes all fields and first photo, shows unique ID and approval stamp
  2. Export list to CSV: current filtered/searched list exported
  3. Summary report: total counts by type, status, woreda, date range

BUSINESS RULES:
  - Only Manager and Supervisor can export
  - Registrar can print their own approved records
  - CSV export never includes photo file data — only metadata
  - PDF print shows watermark "APPROVED" on approved records
  - Draft records cannot be printed (no official status yet)
```

---

## PHASE 5 — SMART CHUNKING

```
CHUNK PLAN — All Modules
=========================

── BACKEND (Go + Gin) ──────────────────────────────────────────

Chunk B-01: Project scaffold
  Go module init, Gin setup, folder structure, env config,
  PostgreSQL connection, base middleware (CORS, logging, error handler)

Chunk B-02: Database schema + migrations
  All CREATE TABLE statements, indexes, foreign keys, seed data (roles)

Chunk B-03: Auth — register & login
  POST /auth/login endpoint, JWT issuance (access + refresh tokens),
  password hashing (bcrypt), user lookup

Chunk B-04: Auth — token refresh & middleware
  POST /auth/refresh, auth middleware (JWT validation),
  role-check middleware (requireRole), attach user to context

Chunk B-05: User management APIs (manager only)
  GET/POST/PUT/DELETE /users, role assignment, deactivate user

Chunk B-06: Immovable record CRUD
  POST/GET/PUT /records/immovable, draft save, auto-ID generation,
  status field, owner filtering

Chunk B-07: Movable record CRUD
  POST/GET/PUT /records/movable, same pattern as B-06

Chunk B-08: Photo upload
  POST /records/:id/photos, multipart upload, file validation,
  store to /media/, link to record in DB

Chunk B-09: Approval workflow APIs
  PUT /records/:id/submit, PUT /records/:id/review,
  PUT /records/:id/approve, PUT /records/:id/return,
  POST /records/:id/comments, GET /records/:id/comments

Chunk B-10: Audit log
  Auto-create status_history entry on every status change,
  GET /records/:id/history

Chunk B-11: Dashboard & search APIs
  GET /dashboard/stats, GET /records (with search + filter + pagination)

Chunk B-12: Export APIs
  GET /export/csv (filtered list), GET /records/:id/pdf

── FRONTEND (React — built in Lovable) ─────────────────────────

Chunk F-01: Project scaffold + routing
  Vite + React, React Router, Tailwind, folder structure,
  protected routes by role, language context (AM/EN)

Chunk F-02: Auth pages
  Login page, logout, JWT storage (httpOnly cookie or memory),
  silent refresh hook, role-based redirect on login

Chunk F-03: Layout shell
  Sidebar navigation (role-filtered links), top navbar,
  notification badge, language switcher, mobile-responsive layout

Chunk F-04: Registrar dashboard + record list
  Stats cards, own record list, status badges, search bar

Chunk F-05: Immovable Form — Sections 1–4
  Identity, location, GPS button, ownership fields,
  section progress indicator, auto-save

Chunk F-06: Immovable Form — Sections 5–8
  History, condition matrix, heritage values, conservation,
  photo uploader, submit button

Chunk F-07: Movable Form — all sections
  Same pattern as F-05/F-06 for Form 01 fields

Chunk F-08: Supervisor view
  Pending queue list, record detail view (read-only),
  approve / return with comment modal

Chunk F-09: Manager view
  Reviewed queue, final approval, user management page,
  create/edit/deactivate user

Chunk F-10: Dashboard (supervisor + manager)
  Full stats, search + filter panel, export buttons

CHUNK RULES:
  - Each chunk = one PR, one focused task
  - Backend chunks must precede their frontend counterparts
  - B-01 and B-02 must be completed before any other backend chunk
  - F-01 and F-02 must be completed before any other frontend chunk
```

---

## PHASE 6 — DATABASE ARCHITECTURE

```
DATABASE ARCHITECTURE
=====================
Database: PostgreSQL
All tables use UUID primary keys.
All tables include: created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ.
Soft delete (deleted_at) on: users.

── ENTITIES ────────────────────────────────────────────────────

users
  id              UUID PK
  full_name       VARCHAR(100) NOT NULL
  email           VARCHAR(150) UNIQUE NOT NULL
  password_hash   VARCHAR(255) NOT NULL
  role            ENUM('registrar','supervisor','manager') NOT NULL
  language        ENUM('am','en') DEFAULT 'am'
  is_active       BOOLEAN DEFAULT TRUE
  deleted_at      TIMESTAMPTZ NULL
  created_at, updated_at

refresh_tokens
  id              UUID PK
  user_id         UUID FK → users.id ON DELETE CASCADE
  token_hash      VARCHAR(255) NOT NULL
  expires_at      TIMESTAMPTZ NOT NULL
  created_at

immovable_records
  id                    UUID PK
  record_id             VARCHAR(30) UNIQUE (ET-HR-AN-I-YYYY-NNNN)
  registrar_id          UUID FK → users.id
  status                ENUM('draft','pending_review','under_review','returned','approved')
  -- Section 1: Identity
  name_amharic          VARCHAR(200) NOT NULL
  name_local            VARCHAR(200)
  category              VARCHAR(10)[]  -- array: [I, II, VII, VIII]
  current_use           VARCHAR(50)[]  -- array of uses
  current_use_other     VARCHAR(200)
  previous_id           VARCHAR(100)
  -- Section 1: Location
  woreda                VARCHAR(100) NOT NULL
  kebele                VARCHAR(100) NOT NULL
  house_number          VARCHAR(50)
  street_number         VARCHAR(50)
  gate                  VARCHAR(50)
  -- Section 2: Ownership
  owner_type            ENUM('public','government','religion','private','waqf')
  owner_name            VARCHAR(200)
  map_reference         VARCHAR(200)
  gps_east              DECIMAL(10,6)
  gps_north             DECIMAL(10,6)
  elevation_m           DECIMAL(8,2)
  -- Section 3: History
  built_by              VARCHAR(200)
  construction_period   VARCHAR(100)
  age_method            ENUM('estimated','exact','relative')
  height_m              DECIMAL(8,2)
  length_m              DECIMAL(8,2)
  width_m               DECIMAL(8,2)
  num_doors             INTEGER
  num_windows           INTEGER
  num_rooms             INTEGER
  material              TEXT
  description           TEXT
  harari_house_grades   VARCHAR(50)[]
  neighborhood_type     VARCHAR(50)
  -- Section 3.8: Condition
  overall_condition     ENUM('very_good','good','damaged','severely_damaged')
  damage_roof           ENUM('minor','moderate','medium','severe')
  damage_cornice        ENUM('minor','moderate','medium','severe')
  damage_wall           ENUM('minor','moderate','medium','severe')
  damage_floor          ENUM('minor','moderate','medium','severe')
  damage_door           ENUM('minor','moderate','medium','severe')
  damage_cupboard       ENUM('minor','moderate','medium','severe')
  damage_upper_floor    ENUM('minor','moderate','medium','severe')
  damage_dera           ENUM('minor','moderate','medium','severe')
  damage_pillar         ENUM('minor','moderate','medium','severe')
  -- Section 3.8: Heritage Values
  value_historical      TEXT
  value_craftsmanship   TEXT
  value_artistic        TEXT
  value_scientific      TEXT
  value_cultural        TEXT
  -- Section 3.9–3.14: Conservation
  has_threat            BOOLEAN
  maintenance_done      BOOLEAN
  maintenance_reason    VARCHAR(300)
  maintenance_by        VARCHAR(200)
  maintenance_date      DATE
  maintenance_count     INTEGER
  preventive_level      ENUM('very_good','good','medium','low','very_low')
  accessibility         ENUM('very_good','good','medium','low','very_low','none')
  notes                 TEXT
  -- Section 4: Related records
  related_docs          VARCHAR(30)[]
  has_oral_history      BOOLEAN
  -- Section 5–7: Persons
  caretaker_name        VARCHAR(200)
  caretaker_role        VARCHAR(200)
  informant_name        VARCHAR(200)
  informant_sex         ENUM('male','female')
  informant_age         INTEGER
  registrar_date        DATE
  -- Section 8: Supervisor comment (stored in comments table)
  approved_at           TIMESTAMPTZ
  approved_by           UUID FK → users.id
  created_at, updated_at

movable_records
  id                    UUID PK
  record_id             VARCHAR(30) UNIQUE (ET-HR-AN-V-YYYY-NNNN)
  registrar_id          UUID FK → users.id
  status                ENUM('draft','pending_review','under_review','returned','approved')
  -- Identity
  name_amharic          VARCHAR(200) NOT NULL
  name_local            VARCHAR(200)
  category              VARCHAR(10)   -- III, IV, V, VI
  -- Location
  location_name         VARCHAR(200)
  woreda                VARCHAR(100)
  kebele                VARCHAR(100)
  house_number          VARCHAR(50)
  current_use           VARCHAR(200)
  previous_id           VARCHAR(100)
  -- Ownership
  owner_type            ENUM('public','government','religion','private')
  owner_name            VARCHAR(200)
  storage_location      ENUM('museum','store','church','private_home','other')
  storage_location_other VARCHAR(200)
  -- History
  made_by               VARCHAR(200)
  period_made           VARCHAR(100)
  age_method            ENUM('estimated','exact','relative')
  acquisition_methods   VARCHAR(30)[]
  -- Measurements
  height_cm             DECIMAL(8,2)
  width_cm              DECIMAL(8,2)
  length_cm             DECIMAL(8,2)
  diameter_cm           DECIMAL(8,2)
  thickness_cm          DECIMAL(8,2)
  weight_kg             DECIMAL(8,2)
  num_pages             INTEGER
  num_chapters          INTEGER
  num_illustrations     INTEGER
  -- Characteristics
  color_type            VARCHAR(200)
  has_decoration        BOOLEAN
  materials             VARCHAR(30)[]
  material_other        VARCHAR(200)
  description           TEXT
  notable_because       VARCHAR(30)[]
  notable_other         TEXT
  significance          TEXT
  -- Condition
  condition             ENUM('good','fair','damaged','incomplete')
  has_threat            BOOLEAN
  threat_description    TEXT
  maintenance_done      BOOLEAN
  maintenance_by        VARCHAR(200)
  maintenance_date      DATE
  maintenance_count     INTEGER
  preventive_level      ENUM('very_good','good','medium','low','very_low')
  accessibility         ENUM('very_good','good','medium','low','very_low','none')
  notes                 TEXT
  -- Related docs
  related_docs          VARCHAR(30)[]
  -- Persons
  informant_name        VARCHAR(200)
  informant_sex         ENUM('male','female')
  informant_age         INTEGER
  informant_occupation  VARCHAR(200)
  caretaker_name        VARCHAR(200)
  caretaker_role        VARCHAR(200)
  registrar_date        DATE
  approved_at           TIMESTAMPTZ
  approved_by           UUID FK → users.id
  created_at, updated_at

record_photos
  id              UUID PK
  record_type     ENUM('immovable','movable') NOT NULL
  record_id       UUID NOT NULL  -- references either table
  file_path       VARCHAR(500) NOT NULL
  file_name       VARCHAR(255)
  file_size_bytes INTEGER
  uploaded_by     UUID FK → users.id
  created_at

record_comments
  id              UUID PK
  record_type     ENUM('immovable','movable') NOT NULL
  record_id       UUID NOT NULL
  author_id       UUID FK → users.id
  comment_text    TEXT NOT NULL
  created_at

status_history
  id              UUID PK
  record_type     ENUM('immovable','movable') NOT NULL
  record_id       UUID NOT NULL
  changed_by      UUID FK → users.id
  from_status     VARCHAR(30)
  to_status       VARCHAR(30) NOT NULL
  note            TEXT
  created_at

── RELATIONSHIPS ────────────────────────────────────────────────

users            ||--o{ immovable_records  : registers
users            ||--o{ movable_records    : registers
users            ||--o{ record_comments   : writes
users            ||--o{ status_history    : triggers
immovable_records ||--o{ record_photos   : has
movable_records   ||--o{ record_photos   : has
immovable_records ||--o{ record_comments : has
movable_records   ||--o{ record_comments : has
immovable_records ||--o{ status_history  : has
movable_records   ||--o{ status_history  : has

── INDEXES ──────────────────────────────────────────────────────

CREATE INDEX idx_immovable_registrar  ON immovable_records(registrar_id);
CREATE INDEX idx_immovable_status     ON immovable_records(status);
CREATE INDEX idx_immovable_woreda     ON immovable_records(woreda);
CREATE INDEX idx_immovable_name       ON immovable_records(name_amharic);
CREATE INDEX idx_movable_registrar    ON movable_records(registrar_id);
CREATE INDEX idx_movable_status       ON movable_records(status);
CREATE INDEX idx_photos_record        ON record_photos(record_type, record_id);
CREATE INDEX idx_comments_record      ON record_comments(record_type, record_id);
CREATE INDEX idx_history_record       ON status_history(record_type, record_id);

── CONSTRAINTS ──────────────────────────────────────────────────

- record_id is unique and auto-generated on first INSERT (draft)
- status transitions must follow allowed paths (enforced in service layer)
- A "Returned" action must have a comment (enforced in service layer)
- Approved records: approved_at and approved_by must be set
```

---

## PHASE 7 — API CONTRACT

```
API CONTRACT
============
Base URL:     /api/v1
Auth:         Bearer JWT in Authorization header
Content-Type: application/json
Response format (always):
  Success: { "success": true,  "data": {...},  "message": "..." }
  Error:   { "success": false, "error": "...", "code": 400 }

── AUTH ENDPOINTS ───────────────────────────────────────────────

POST /auth/login
  Auth: None
  Request:  { email, password }
  Response: { data: { access_token, refresh_token, user: {id,name,role,language} } }
  Errors:   401 invalid credentials | 403 account deactivated

POST /auth/refresh
  Auth: None (uses refresh_token in body)
  Request:  { refresh_token }
  Response: { data: { access_token } }
  Errors:   401 expired/invalid refresh token

POST /auth/logout
  Auth: Any role
  Request:  { refresh_token }
  Response: { message: "Logged out" }

── USER MANAGEMENT ──────────────────────────────────────────────

GET /users
  Auth: Manager only
  Query: ?page=1&limit=20&role=registrar&is_active=true
  Response: { data: { users: [...], total, page } }

POST /users
  Auth: Manager only
  Request:  { full_name, email, password, role, language }
  Response: { data: { user } }
  Errors:   409 email already exists

PUT /users/:id
  Auth: Manager only
  Request:  { full_name, role, language, is_active }
  Response: { data: { user } }
  Errors:   403 cannot deactivate yourself

GET /users/me
  Auth: Any role
  Response: { data: { user } }

PUT /users/me/language
  Auth: Any role
  Request:  { language: "am" | "en" }
  Response: { data: { language } }

── IMMOVABLE RECORDS ────────────────────────────────────────────

POST /records/immovable
  Auth: Registrar
  Request:  { all form fields as JSON object, status: "draft" }
  Response: { data: { record_id, id, status } }

GET /records/immovable
  Auth: All roles (filtered by role)
  Query: ?status=&woreda=&search=&page=1&limit=20&date_from=&date_to=
  Response: { data: { records: [...], total, page } }

GET /records/immovable/:id
  Auth: All roles (registrar: own only)
  Response: { data: { record, photos, comments, history } }

PUT /records/immovable/:id
  Auth: Registrar (own draft or returned records only)
  Request:  { updated fields }
  Response: { data: { record } }
  Errors:   403 not owner | 409 record not in editable status

PUT /records/immovable/:id/submit
  Auth: Registrar (own draft/returned records only)
  Response: { data: { status: "pending_review" } }
  Errors:   422 validation failed (missing required fields)

── MOVABLE RECORDS ──────────────────────────────────────────────

POST /records/movable         (same pattern as immovable)
GET  /records/movable         (same pattern)
GET  /records/movable/:id     (same pattern)
PUT  /records/movable/:id     (same pattern)
PUT  /records/movable/:id/submit  (same pattern)

── APPROVAL WORKFLOW ────────────────────────────────────────────

PUT /records/:type/:id/review-approve
  Auth: Supervisor only
  Request:  { comment? }
  Response: { data: { status: "under_review" } }
  Errors:   403 not supervisor | 409 wrong status

PUT /records/:type/:id/review-return
  Auth: Supervisor only
  Request:  { comment: required }
  Response: { data: { status: "returned" } }
  Errors:   422 comment required

PUT /records/:type/:id/final-approve
  Auth: Manager only
  Request:  { comment? }
  Response: { data: { status: "approved", approved_at } }
  Errors:   403 not manager | 409 wrong status

PUT /records/:type/:id/final-return
  Auth: Manager only
  Request:  { comment: required }
  Response: { data: { status: "pending_review" } }

POST /records/:type/:id/comments
  Auth: Supervisor, Manager
  Request:  { comment_text }
  Response: { data: { comment } }

GET /records/:type/:id/comments
  Auth: All roles (registrar: own records only)
  Response: { data: { comments: [...] } }

GET /records/:type/:id/history
  Auth: All roles (registrar: own records only)
  Response: { data: { history: [...] } }

── PHOTOS ───────────────────────────────────────────────────────

POST /records/:type/:id/photos
  Auth: Registrar (own draft/returned)
  Content-Type: multipart/form-data
  Request:  { photo: File }
  Response: { data: { photo_id, file_path } }
  Errors:   413 file too large | 415 unsupported type | 400 max 10 photos

DELETE /records/:type/:id/photos/:photo_id
  Auth: Registrar (own draft/returned)
  Response: { message: "Photo deleted" }

── DASHBOARD & SEARCH ───────────────────────────────────────────

GET /dashboard/stats
  Auth: All roles (role-filtered)
  Response: { data: {
    total_immovable, total_movable,
    by_status: { draft, pending_review, under_review, returned, approved },
    pending_my_action  (for supervisor/manager)
  } }

── EXPORT ───────────────────────────────────────────────────────

GET /export/records/csv
  Auth: Manager, Supervisor
  Query: same filters as GET /records
  Response: CSV file download

GET /records/:type/:id/pdf
  Auth: Manager, Supervisor, Registrar (own approved)
  Response: PDF file download
```

---

## PHASE 8 — FRONTEND ARCHITECTURE

```
FRONTEND ARCHITECTURE
=====================
Framework:        React + Vite (PWA)
Routing:          React Router v6
Styling:          Tailwind CSS
State:            Zustand (global: auth, language)
Data fetching:    TanStack Query (React Query v5)
API Layer:        Axios with interceptors (auto token refresh)
Forms:            React Hook Form + Zod validation
i18n:             i18next (Amharic + English JSON files)
PWA:              Vite PWA plugin (offline draft saving via IndexedDB)
PDF export:       react-pdf or server-generated

── PAGES BY ROLE ────────────────────────────────────────────────

Shared:
  /login                  LoginPage (public)
  /unauthorized           UnauthorizedPage (public)

Registrar (/registrar/*):
  /registrar/dashboard     RegistrarDashboard
  /registrar/records       MyRecordsList
  /registrar/records/:id   RecordDetail (read-only if submitted)
  /registrar/new/immovable ImmovableForm (new)
  /registrar/new/movable   MovableForm (new)
  /registrar/edit/:type/:id EditForm (draft or returned only)

Supervisor (/supervisor/*):
  /supervisor/dashboard    SupervisorDashboard
  /supervisor/pending      PendingQueue
  /supervisor/records      AllRecordsList
  /supervisor/records/:id  RecordDetail + ReviewActions

Manager (/manager/*):
  /manager/dashboard       ManagerDashboard
  /manager/pending         ReviewedQueue
  /manager/records         AllRecordsList
  /manager/records/:id     RecordDetail + FinalApproveActions
  /manager/users           UserManagementPage
  /manager/users/new       CreateUserPage
  /manager/users/:id/edit  EditUserPage
  /manager/reports         ReportsPage

── SHARED COMPONENTS ────────────────────────────────────────────

Layout/
  AppLayout           Main shell with sidebar + topbar
  Sidebar             Role-filtered navigation links
  Topbar              Language switcher + notification badge + user menu
  ProtectedRoute      Wraps pages — checks auth + role

UI/
  StatusBadge         Colored badge for record status
  RecordCard          Card in list views (title, ID, status, date)
  CommentThread       List of comments with author + date
  StatusTimeline      Visual audit log (history of status changes)
  PhotoUploader       Drag & drop + camera capture, preview grid
  GpsButton           One-tap GPS capture component
  ConfirmModal        Reusable confirm/cancel dialog
  ReturnModal         Modal with required comment textarea
  SectionProgress     Progress indicator across form sections
  LoadingSpinner      Full-page and inline variants
  EmptyState          Empty list illustration + CTA
  ErrorBoundary       Catches render errors per section

Form/
  ImmovableFormSection1   through Section8 (one component per section)
  MovableFormSections     (equivalent components for Form 01)
  AutoSaveIndicator       "Saved X seconds ago" indicator

── STATE MANAGEMENT ─────────────────────────────────────────────

authStore (Zustand):
  user, accessToken, isAuthenticated
  actions: login(), logout(), refreshToken(), updateLanguage()

languageStore (Zustand):
  language ('am' | 'en')
  actions: setLanguage()

── API LAYER ────────────────────────────────────────────────────

/src/api/
  client.ts       Axios instance, base URL, request/response interceptors,
                  auto-refresh on 401
  auth.ts         login, logout, refresh
  users.ts        getUsers, createUser, updateUser
  immovable.ts    getRecords, getRecord, createRecord, updateRecord,
                  submitRecord, approveRecord, returnRecord
  movable.ts      (same pattern)
  photos.ts       uploadPhoto, deletePhoto
  workflow.ts     reviewApprove, reviewReturn, finalApprove, finalReturn,
                  postComment, getComments, getHistory
  dashboard.ts    getStats
  export.ts       downloadCSV, downloadPDF

── COMPONENT TREE ───────────────────────────────────────────────

App
├── Router
│   ├── /login → LoginPage
│   └── ProtectedRoute
│       ├── AppLayout
│       │   ├── Sidebar
│       │   ├── Topbar
│       │   └── Outlet
│       │       ├── [Registrar pages]
│       │       ├── [Supervisor pages]
│       │       └── [Manager pages]
│       └── UnauthorizedPage

── MOBILE REQUIREMENTS ──────────────────────────────────────────

- Sidebar collapses to bottom navigation bar on mobile
- Forms scroll vertically section by section
- Photo upload uses native camera on mobile
- GPS button uses navigator.geolocation API
- All touch targets minimum 44x44px
- Offline: drafts saved to IndexedDB, synced on reconnect
- PWA installable from browser (Add to Home Screen)
```

---

## PHASE 9 — DEVELOPMENT WORKFLOW

```
DEVELOPMENT WORKFLOW
====================

BUILD ORDER (strict):
  1.  [B-01] Go project scaffold + Gin + PostgreSQL connection     [S]
  2.  [B-02] Database migrations — all tables + indexes            [M]
  3.  [B-03] Auth — login + JWT issuance + bcrypt                  [M]
  4.  [B-04] Auth middleware + role guard + token refresh          [M]
  5.  [B-05] User management CRUD (manager only)                   [M]
  6.  [B-06] Immovable record CRUD + auto-ID + status field        [L]
  7.  [B-07] Movable record CRUD                                   [L]
  8.  [B-08] Photo upload (multipart + file storage)               [M]
  9.  [B-09] Approval workflow endpoints                           [M]
  10. [B-10] Status audit log                                      [S]
  11. [B-11] Dashboard stats + search/filter/pagination            [M]
  12. [B-12] CSV + PDF export                                      [M]
  13. [F-01] React scaffold + routing + i18n + Zustand             [M]
  14. [F-02] Login page + auth flow + token handling               [M]
  15. [F-03] App layout shell (sidebar, topbar, mobile nav)        [M]
  16. [F-04] Registrar dashboard + record list                     [M]
  17. [F-05] Immovable form — sections 1–4                         [L]
  18. [F-06] Immovable form — sections 5–8 + photos                [L]
  19. [F-07] Movable form — all sections + photos                  [L]
  20. [F-08] Supervisor review view + approve/return               [M]
  21. [F-09] Manager final approval + user management              [M]
  22. [F-10] Full dashboard + search + export                      [M]

ENVIRONMENT SETUP:
  Backend:
    - Go 1.22+
    - PostgreSQL 16
    - golang-migrate (migrations)
    - Gin web framework
    - golang-jwt/jwt/v5
    - bcrypt (golang.org/x/crypto)
    - godotenv (env vars)
    - uuid (google/uuid)
  Frontend:
    - Node 20+
    - Vite 5 + React 18
    - React Router v6
    - TanStack Query v5
    - Zustand
    - React Hook Form + Zod
    - Axios
    - i18next + react-i18next
    - Tailwind CSS 3
    - Vite PWA plugin

FOLDER STRUCTURE — Backend:
  /cmd/server/main.go
  /internal/
    auth/         handler, service, repository
    users/        handler, service, repository
    immovable/    handler, service, repository
    movable/      handler, service, repository
    workflow/     handler, service
    photos/       handler, service
    dashboard/    handler, service
    export/       handler, service
    middleware/   auth.go, role.go, logger.go, error.go
    db/           connection.go, migrations/
    models/       all struct definitions
    config/       env.go
  /media/         uploaded photos (gitignored)

FOLDER STRUCTURE — Frontend:
  /src/
    api/          client.ts + per-module files
    components/   Layout/, UI/, Form/
    pages/        registrar/, supervisor/, manager/, shared/
    stores/       authStore.ts, languageStore.ts
    hooks/        useAuth.ts, useAutoSave.ts, useGPS.ts
    i18n/         am.json, en.json
    types/        all TypeScript interfaces
    utils/        formatters, validators

TESTING STRATEGY:
  Backend unit tests: auth service, ID generation, status transition guards
  Backend integration tests: full workflow (submit → review → approve)
  Frontend: form validation, role-based route protection
  Manual E2E: full registrar → supervisor → manager flow

DEPLOYMENT (Ubuntu Linux):
  PostgreSQL 16 → systemd service
  Go binary   → built with `go build`, run via systemd service
  React build → `npm run build` → served as static files by Nginx
  Nginx       → reverse proxy /api/* to Go, serve /* from React build
  SSL         → Certbot Let's Encrypt
  Backups     → daily pg_dump cron job + /media/ rsync to backup drive
  Logs        → /var/log/qirs-mezgeb/
```

---

## PHASE 10 — MASTER PROMPT LIBRARY

---

### MASTER CONTEXT BLOCK
*Paste this at the top of EVERY Cursor or Lovable prompt.*

```
# MASTER PROJECT CONTEXT — Qirs Mezgeb (ቅርስ መዝገብ)
# Paste this block at the top of every AI prompt for this project.

Project:      Qirs Mezgeb — Heritage Registry System
Type:         Government management system (Web + Mobile PWA)
Institution:  Harari Region Culture, Heritage & Tourism Bureau, Ethiopia
Purpose:      Digitize the paper-based heritage registration workflow for
              immovable (Form 02) and movable (Form 01) heritage assets,
              with a 3-stage approval chain: Registrar → Supervisor → Manager.

Actors:
  - Registrar:  Creates and submits records. Edits drafts and returned records only.
  - Supervisor: Reviews submitted records. Approves to manager or returns with comment.
  - Manager:    Final approval. Manages users. Full access to all records and reports.

Record statuses:
  draft → pending_review → under_review → approved
                        ↘ returned (to registrar)

Modules:
  1. Auth & User Management
  2. Immovable Heritage Registration (Form 02)
  3. Movable Heritage Registration (Form 01)
  4. Approval Workflow
  5. Dashboard & Search
  6. Reports & Export

Frontend Stack:
  React 18 + Vite (PWA) | React Router v6 | Tailwind CSS
  TanStack Query v5 | Zustand | React Hook Form + Zod
  Axios with auto-refresh interceptor | i18next (Amharic + English)

Backend Stack:
  Go 1.22 + Gin | PostgreSQL 16 | golang-migrate
  golang-jwt/jwt/v5 | bcrypt | google/uuid | godotenv

Architecture Rules:
  - Backend: Handler → Service → Repository pattern (no business logic in handlers)
  - Frontend: Pages call hooks, hooks call API layer, API layer calls Axios client
  - All API responses: { "success": bool, "data": {...}, "message": "..." }
  - All routes protected by JWT middleware + role guard
  - All form inputs validated both client-side (Zod) and server-side
  - Registrar can only see own records — enforced at DB query level
  - Status transitions enforced in service layer, not just frontend
  - Soft delete on users (deleted_at), never hard delete records
  - All tables have created_at, updated_at
  - UUID primary keys everywhere
  - Amharic text support: PostgreSQL UTF-8, Go handles natively, React with i18next

Coding Standards:
  - Go: follow standard Go conventions, exported types in models package
  - React: functional components only, TypeScript strict mode
  - No hardcoded strings — all UI text in am.json / en.json
  - Loading, error, and empty states required on every data-fetching component
  - Mobile-first CSS — test all UI at 375px width
  - No mock data in production code

Current Task: [REPLACE THIS with the specific chunk you are working on]
```

---

### PROMPT LIBRARY

**PROMPT 1 — Backend Scaffold (Chunk B-01)**
```
[PASTE MASTER CONTEXT ABOVE]

Current Task: Chunk B-01 — Go project scaffold

Set up the complete Go backend project for Qirs Mezgeb.

DELIVER:
1. Go module: github.com/qirs-mezgeb/api
2. Folder structure as specified in the architecture rules
3. Gin server with: CORS middleware, request logger, centralized error handler
4. PostgreSQL connection pool (pgx driver) with retry on startup
5. Config loader from .env file (godotenv) with a Config struct containing:
   DB_URL, PORT, JWT_SECRET, JWT_REFRESH_SECRET, MEDIA_PATH, ALLOWED_ORIGINS
6. Health check endpoint: GET /health → { "status": "ok", "db": "connected" }
7. Graceful shutdown on SIGTERM/SIGINT
8. README with setup instructions

Folder structure required:
  /cmd/server/main.go
  /internal/middleware/
  /internal/config/
  /internal/db/
  /internal/models/
  .env.example
  go.mod
```

**PROMPT 2 — Database Migrations (Chunk B-02)**
```
[PASTE MASTER CONTEXT ABOVE]

Current Task: Chunk B-02 — Database schema and migrations

Using golang-migrate, create all migration files for Qirs Mezgeb.

DELIVER:
1. Migration files in /internal/db/migrations/
   - 001_create_users.up.sql / .down.sql
   - 002_create_refresh_tokens.up.sql / .down.sql
   - 003_create_immovable_records.up.sql / .down.sql
   - 004_create_movable_records.up.sql / .down.sql
   - 005_create_record_photos.up.sql / .down.sql
   - 006_create_record_comments.up.sql / .down.sql
   - 007_create_status_history.up.sql / .down.sql
   - 008_create_indexes.up.sql / .down.sql
2. Seed file: 009_seed_admin_user.up.sql
   (creates one manager account: admin@qirsmezgeb.gov.et / Admin1234)
3. Migration runner function called on app startup
4. All ENUM types defined as PostgreSQL ENUMs
5. All schemas exactly match the DATABASE ARCHITECTURE section of the design doc

Include the full SQL for every table. Do not use an ORM — raw SQL migrations only.
```

**PROMPT 3 — Auth Endpoints (Chunk B-03 + B-04)**
```
[PASTE MASTER CONTEXT ABOVE]

Current Task: Chunks B-03 + B-04 — Authentication

Build the complete authentication system for Qirs Mezgeb backend.

DELIVER:
1. /internal/models/user.go — User struct + Role constants
2. /internal/auth/repository.go — GetUserByEmail, CreateRefreshToken,
   GetRefreshToken, DeleteRefreshToken
3. /internal/auth/service.go — Login (validate password, issue JWT pair),
   Refresh (validate refresh token, issue new access token), Logout
4. /internal/auth/handler.go — POST /auth/login, POST /auth/refresh, POST /auth/logout
5. /internal/middleware/auth.go — JWT validation middleware, attaches user to Gin context
6. /internal/middleware/role.go — RequireRole(...roles) middleware factory
7. JWT: access token 15 min expiry, refresh token 7 days, signed with separate secrets
8. Refresh tokens stored hashed in DB (bcrypt hash of the token)
9. On logout: delete refresh token from DB

Test cases to include (table-driven Go tests):
  - Login with valid credentials → returns token pair
  - Login with wrong password → 401
  - Login with deactivated account → 403
  - Refresh with valid token → new access token
  - Refresh with expired token → 401
```

**PROMPT 4 — Immovable Record CRUD (Chunk B-06)**
```
[PASTE MASTER CONTEXT ABOVE]

Current Task: Chunk B-06 — Immovable heritage record CRUD

Build the complete CRUD API for immovable heritage records.

DELIVER:
1. /internal/models/immovable.go — ImmovableRecord struct (all fields from schema)
2. /internal/immovable/repository.go
   - Create(record) → inserts, generates record_id on first save
   - GetByID(id, userID, role) → enforces registrar sees own only
   - List(filters, userID, role) → paginated, role-filtered, searchable
   - Update(id, fields, userID, role) → only draft/returned, only owner
   - UpdateStatus(id, newStatus, changedByID) → validates transition
3. /internal/immovable/service.go — business logic + validation
4. /internal/immovable/handler.go
   - POST   /api/v1/records/immovable
   - GET    /api/v1/records/immovable
   - GET    /api/v1/records/immovable/:id
   - PUT    /api/v1/records/immovable/:id
   - PUT    /api/v1/records/immovable/:id/submit
5. Auto-ID generation: ET-HR-AN-I-[YYYY]-[NNNN]
   NNNN is the sequential count of immovable records in that year, zero-padded
6. Query params for GET list: status, woreda, search, page, limit, date_from, date_to

Status transition rules (enforce in service):
  draft          → pending_review  (submit action, registrar only)
  pending_review → under_review    (supervisor approve)
  pending_review → returned        (supervisor return)
  under_review   → approved        (manager approve)
  under_review   → returned        (manager return — goes to pending_review)
  returned       → pending_review  (registrar resubmits)
  approved       → [LOCKED — no transitions]
```

**PROMPT 5 — Approval Workflow (Chunk B-09)**
```
[PASTE MASTER CONTEXT ABOVE]

Current Task: Chunk B-09 — Approval workflow endpoints

Build the approval workflow API covering both record types.

DELIVER:
1. /internal/workflow/service.go
   - ReviewApprove(recordType, recordID, supervisorID, comment?)
   - ReviewReturn(recordType, recordID, supervisorID, comment REQUIRED)
   - FinalApprove(recordType, recordID, managerID, comment?)
   - FinalReturn(recordType, recordID, managerID, comment REQUIRED)
   - AddComment(recordType, recordID, authorID, text)
   - GetComments(recordType, recordID, requesterID, role)
   - GetHistory(recordType, recordID, requesterID, role)
2. /internal/workflow/handler.go
   - PUT  /api/v1/records/:type/:id/review-approve  [supervisor]
   - PUT  /api/v1/records/:type/:id/review-return   [supervisor]
   - PUT  /api/v1/records/:type/:id/final-approve   [manager]
   - PUT  /api/v1/records/:type/:id/final-return    [manager]
   - POST /api/v1/records/:type/:id/comments        [supervisor, manager]
   - GET  /api/v1/records/:type/:id/comments        [all, registrar: own only]
   - GET  /api/v1/records/:type/:id/history         [all, registrar: own only]
3. On every status change: auto-insert row into status_history
4. "Return" actions: if comment is empty string → 422 error
5. :type param must be "immovable" or "movable" → 400 if anything else
```

**PROMPT 6 — React Scaffold (Chunk F-01 + F-02)**
```
[PASTE MASTER CONTEXT ABOVE]

Current Task: Chunks F-01 + F-02 — React frontend scaffold + auth

Set up the complete React frontend project and authentication flow for Qirs Mezgeb.
This frontend will be built in Lovable and integrated with the Go backend.

DELIVER:
1. Vite + React 18 + TypeScript project
2. Folder structure as specified in architecture:
   /src/api/, /src/components/, /src/pages/, /src/stores/, /src/hooks/,
   /src/i18n/, /src/types/, /src/utils/
3. /src/api/client.ts
   - Axios instance with baseURL from env var VITE_API_URL
   - Request interceptor: attach access token from authStore
   - Response interceptor: on 401, attempt silent refresh, retry request once,
     on second 401 logout and redirect to /login
4. /src/stores/authStore.ts (Zustand)
   - State: user, accessToken, isAuthenticated
   - Actions: login(), logout(), setTokens()
5. /src/stores/languageStore.ts (Zustand)
   - State: language ('am' | 'en')
   - Actions: setLanguage(), persisted to localStorage
6. /src/i18n/index.ts — i18next setup with am.json and en.json
7. /src/pages/shared/LoginPage.tsx — email/password form, calls POST /auth/login,
   redirects to role-appropriate dashboard on success
8. /src/components/Layout/ProtectedRoute.tsx
   - Checks isAuthenticated + role
   - Redirects to /login if not authenticated
   - Redirects to /unauthorized if wrong role
9. /src/App.tsx — full route tree for all roles with ProtectedRoute wrappers
10. /src/types/index.ts — TypeScript interfaces for User, ImmovableRecord,
    MovableRecord, Comment, StatusHistory, PaginatedResponse, ApiResponse

All UI text must use i18next t() function — no hardcoded strings.
Mobile-first layout.
```

**PROMPT 7 — Immovable Form UI (Chunks F-05 + F-06)**
```
[PASTE MASTER CONTEXT ABOVE]

Current Task: Chunks F-05 + F-06 — Immovable heritage registration form UI

Build the complete Form 02 digital form for immovable heritage registration.
This is the most important UI in the system.

DELIVER:
1. /src/pages/registrar/ImmovableFormPage.tsx — main page component
   - Loads existing draft if :id param present
   - Shows 8 section tabs with progress indicator
   - "Save Draft" button always visible (calls PUT /records/immovable/:id)
   - "Submit for Review" button only enabled when required fields are filled
   - Auto-save every 2 minutes (useAutoSave hook)

2. Form sections as separate components in /src/components/Form/:
   Section1_Identity.tsx    — name fields, category checkboxes, current use
   Section2_Location.tsx    — woreda, kebele, house/street/gate, GpsButton
   Section3_Ownership.tsx   — owner type radio cards, owner name, map ref
   Section4_History.tsx     — built by, period, age method, measurements, rooms
   Section5_Classification.tsx — Harari house grade checkboxes, neighborhood
   Section6_Condition.tsx   — overall condition selector, damage matrix table
   Section7_Values.tsx      — 5 heritage value text areas
   Section8_Conservation.tsx — threat, maintenance, accessibility, notes + informant fields

3. /src/components/UI/GpsButton.tsx
   - Calls navigator.geolocation.getCurrentPosition
   - Fills gps_east, gps_north, elevation fields on success
   - Shows error if GPS unavailable

4. /src/components/UI/PhotoUploader.tsx
   - Accepts up to 10 photos (JPG/PNG, max 5MB each)
   - Shows preview grid with delete button per photo
   - On mobile: opens camera or gallery via input[type=file][accept="image/*"][capture]
   - Uploads immediately via POST /records/:type/:id/photos

5. /src/hooks/useAutoSave.ts
   - Debounced save every 2 minutes
   - Shows "Saved X seconds ago" indicator
   - On network error: saves to IndexedDB for offline recovery

6. All field labels bilingual: shown in Amharic with English subtitle
7. All required fields marked and validated with Zod before submit
8. Mobile-first: each section full width, large touch targets, sticky save buttons
```

**PROMPT 8 — Supervisor Review UI (Chunk F-08)**
```
[PASTE MASTER CONTEXT ABOVE]

Current Task: Chunk F-08 — Supervisor review interface

Build the supervisor's review interface.

DELIVER:
1. /src/pages/supervisor/PendingQueuePage.tsx
   - Lists all records with status "pending_review"
   - Shows: record name, ID, type (immovable/movable), registrar name, submitted date
   - Sorted by submitted date (oldest first)
   - Click row → go to RecordDetailPage

2. /src/pages/supervisor/RecordDetailPage.tsx
   - Read-only display of ALL form fields (same layout as registrar form, view mode)
   - Shows all uploaded photos in a gallery
   - Shows StatusTimeline component (full history)
   - Shows CommentThread component (existing comments)
   - Two action buttons: "Approve to Manager" and "Return to Registrar"

3. /src/components/UI/ReturnModal.tsx
   - Opens when "Return to Registrar" clicked
   - Textarea for comment (required — cannot submit empty)
   - Confirm / Cancel buttons
   - Calls PUT /records/:type/:id/review-return

4. /src/components/UI/StatusTimeline.tsx
   - Visual list of all status_history entries
   - Shows: status badge, changed by (name + role), date, optional note

5. /src/components/UI/CommentThread.tsx
   - Lists all comments with author name, role badge, timestamp
   - Text input + post button for adding new comment (supervisor/manager only)

All text bilingual via i18next.
```

**PROMPT 9 — CURSOR.md Rules File**
```
This is NOT a code prompt. Generate the CURSOR.md file for the Qirs Mezgeb project.
This file will be placed at the root of the repository and read by Cursor AI
before every code generation task.

[PASTE MASTER CONTEXT ABOVE]

CURSOR.md must include:
1. Project overview (2 paragraphs)
2. Folder structure (full tree, backend + frontend)
3. Architecture rules (Handler → Service → Repository, etc.)
4. API response format standard
5. Status transition table (all valid transitions)
6. Role permission matrix
7. Naming conventions (Go: snake_case files, PascalCase types; React: PascalCase components, camelCase hooks)
8. What NOT to do (no business logic in handlers, no hardcoded strings, no mock data, no hard deletes)
9. How to add a new API endpoint (step-by-step checklist)
10. How to add a new form field (step-by-step: DB migration → model → handler → frontend type → form component → i18n keys)
11. Environment variables list with descriptions
12. Common patterns used (with short code examples):
    - Role guard middleware usage
    - Standard paginated query pattern
    - Standard error response pattern
    - How to write a Zustand store
    - How to write a TanStack Query hook
```

---

## PROJECT RULES SUMMARY

```
PROJECT RULES — Qirs Mezgeb
============================

NEVER:
  ✗ Put business logic in HTTP handlers
  ✗ Hard-delete any heritage record
  ✗ Allow a registrar to see another registrar's records
  ✗ Allow status transitions outside the defined flow
  ✗ Allow "Return" without a written comment
  ✗ Allow editing of an Approved record
  ✗ Hardcode UI text — always use i18next
  ✗ Use mock/fake data in production code
  ✗ Skip validation on any API input

ALWAYS:
  ✓ Validate inputs on both client (Zod) and server (Go)
  ✓ Write a status_history row on every status change
  ✓ Return { success, data, message } on every API response
  ✓ Filter records by registrar_id when role = registrar
  ✓ Show loading + error + empty states on every data-fetching component
  ✓ Use UUID for all primary keys
  ✓ Include created_at and updated_at on all tables
  ✓ Test all UI at 375px (mobile) width
  ✓ Prefix all API routes with /api/v1/
```
