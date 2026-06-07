# CURSOR PROMPT 5 — Qirs Mezgeb
# Supervisor Flow: Dashboard + Pending Queue + All Records + Record Detail with Review Actions
# Paste this entire file into Cursor Chat with the project open.

---

## CONTEXT

Continuing from Prompt 4. The project now has:
- `src/api/workflow.ts` with `reviewApprove`, `reviewReturn`, `getComments`, `addComment`, `getHistory`
- `src/components/UI/ReturnModal.tsx` — dialog with required textarea
- `src/components/UI/StatusTimeline.tsx` and `CommentThread.tsx`
- `src/components/Forms/MovableForm.tsx` fully built
- `src/routes/_authenticated/records.$id.edit.tsx` — shows form + tabbed detail panel
- `src/routes/_authenticated/pending.tsx` — currently a placeholder EmptyState
- `src/routes/_authenticated/dashboard.tsx` — currently renders `<RegistrarDashboard />` for all roles

The `_authenticated` layout guard at `src/routes/_authenticated.tsx` already handles auth. Role-specific redirects happen inside each page component.

Existing shared components to reuse:
- `RecordsList` + `RecordCard` — accept a `fetcher` prop and `queryKey`
- `StatusBadge`, `LoadingSpinner`, `EmptyState`
- `ReturnModal`, `StatusTimeline`, `CommentThread`
- `Field`, `FormSection` primitives in `src/components/Forms/Field.tsx`
- shadcn/ui: `Dialog`, `Tabs`, `Card`, `Button`, `Textarea` from `src/components/ui/`

---

## WHAT TO BUILD

### 1. `/src/api/users.ts`

Create this file — needed for supervisor's "view registrar name" in record detail:

```typescript
import { apiClient } from './client'
import type { ApiResponse, UserItem, PaginatedUsers, UserListParams, UpdateLanguageBody } from '@/types'

// GET /users/me  [all roles]
export async function getMe(): Promise<UserItem> {
  const { data } = await apiClient.get<ApiResponse<{ user: UserItem }>>('/users/me')
  return data.data.user
}

// PUT /users/me/language  [all roles]
export async function updateMyLanguage(language: string): Promise<void> {
  await apiClient.put('/users/me/language', { language })
}

// GET /users  [manager only]
export async function listUsers(params: UserListParams = {}): Promise<PaginatedUsers> {
  const { data } = await apiClient.get<ApiResponse<PaginatedUsers>>('/users', { params })
  return data.data
}

// POST /users  [manager only]
export async function createUser(body: import('@/types').CreateUserBody): Promise<UserItem> {
  const { data } = await apiClient.post<ApiResponse<{ user: UserItem }>>('/users', body)
  return data.data.user
}

// PUT /users/:id  [manager only]
export async function updateUser(id: string, body: import('@/types').UpdateUserBody): Promise<UserItem> {
  const { data } = await apiClient.put<ApiResponse<{ user: UserItem }>>(`/users/${id}`, body)
  return data.data.user
}

// DELETE /users/:id  [manager only — soft deactivate]
export async function deactivateUser(id: string): Promise<void> {
  await apiClient.delete(`/users/${id}`)
}
```

---

### 2. `/src/components/Dashboard/SupervisorDashboard.tsx`

Create this component. It fetches from `GET /dashboard/stats` (same endpoint as registrar — server scopes it).

Layout:

```
┌──────────────────────────────────────────────────────┐
│  ዳሽቦርድ / Dashboard                                   │
│  ተቆጣጣሪ / Supervisor                                   │
├──────────────────────────────────────────────────────┤
│  [Stat card]        [Stat card]        [Stat card]   │
│  Pending Review     Under Review       Approved       │
│  stats.by_status    stats.by_status    stats.by_status│
│  .pending_review    .under_review      .approved      │
├──────────────────────────────────────────────────────┤
│  [Stat card]        [Stat card]                       │
│  Total Immovable    Total Movable                     │
│  stats.total_imm    stats.total_mov                   │
├──────────────────────────────────────────────────────┤
│  [Big CTA button — full width]                        │
│  "ለግምገማ የቀረቡ / Go to Pending Queue →"               │
│  → navigates to /pending                              │
└──────────────────────────────────────────────────────┘
```

- Use the existing `StatCard` component from `src/components/Dashboard/StatCard.tsx`
- Pending Review card: amber color, `ListChecks` icon
- Under Review card: blue color, `Eye` icon  
- Approved card: emerald color, `CheckCircle2` icon
- Total Immovable: `Building2` icon, default color
- Total Movable: `Package` icon, default color
- CTA button: full-width, primary color, `ArrowRight` icon
- Loading state: show skeleton cards (use `LoadingSpinner` centered in a card-shaped div)
- Error state: show error message with retry button
- Query key: `['dashboard', 'stats']`

---

### 3. `/src/components/Dashboard/ManagerDashboard.tsx`

Create a placeholder for now (Prompt 6 will flesh it out):

```typescript
// Same stat cards as SupervisorDashboard, but shows ALL statuses including draft and returned.
// For now, render the same layout as SupervisorDashboard but with additional cards:
// - Draft (gray, FileText icon)
// - Returned (rose, RotateCcw icon)
// Then a CTA: "ለመጨረሻ ፈቃድ የቀረቡ / Go to Final Review Queue →" → /reviewed
```

Build it fully — don't stub it. It uses the same `GET /dashboard/stats` endpoint.

---

### 4. Update `/src/routes/_authenticated/dashboard.tsx`

Replace the current "all roles see RegistrarDashboard" with role-specific dashboards:

```typescript
import { RegistrarDashboard } from '@/components/Dashboard/RegistrarDashboard'
import { SupervisorDashboard } from '@/components/Dashboard/SupervisorDashboard'
import { ManagerDashboard } from '@/components/Dashboard/ManagerDashboard'

function DashboardPage() {
  const user = useAuthStore((s) => s.user)
  if (!user) return null
  if (user.role === 'supervisor') return <SupervisorDashboard />
  if (user.role === 'manager') return <ManagerDashboard />
  return <RegistrarDashboard />
}
```

---

### 5. `/src/components/Records/SupervisorRecordDetail.tsx`

This is the heart of Prompt 5. Create a full record detail component for supervisors.

**Props:**
```typescript
interface SupervisorRecordDetailProps {
  recordId: string
  recordType: RecordType
}
```

**Data fetching:**
- If `recordType === 'immovable'`: query `['immovable', recordId]` via `getImmovable(recordId)`
- If `recordType === 'movable'`: query `['movable', recordId]` via `getMovable(recordId)`
- Both return `{ record, photos, comments, history }`

**Layout:**

```
┌──────────────────────────────────────────────────────┐
│  [← Back]  ET-HR-AN-I-2024-0001  [StatusBadge]       │
│  Name: [name_amharic]                                 │
│  Location: woreda / kebele      Type: Immovable       │
├──────────────────────────────────────────────────────┤
│  ACTION BAR (only shown when status = pending_review):│
│                                                       │
│  [✓ Approve]              [↩ Return to Registrar]    │
│  primary button           destructive/outline button  │
└──────────────────────────────────────────────────────┘
│  TABS: [Details] [History] [Comments] [Photos]        │
├──────────────────────────────────────────────────────┤
│  Details tab: read-only key-value grid of all record  │
│               fields that have values (skip nulls)    │
│  History tab: <StatusTimeline history={detail.history}│
│  Comments tab: <CommentThread canComment={true} ...>  │
│  Photos tab: photo grid (read-only)                   │
└──────────────────────────────────────────────────────┘
```

**Action bar logic:**
- Only render when `record.status === 'pending_review'`
- **Approve button**: 
  - Label: `am="አፅድቅ"` `en="Approve"`
  - Calls `reviewApprove(recordType, recordId)` from `src/api/workflow.ts`
  - On success: invalidate `['immovable'/'movable', recordId]` + `['dashboard', 'stats']` + `['records']`, show success toast using `sonner`, navigate to `/pending`
  - Shows `Loader2` spinner while pending
- **Return button**:
  - Label: `am="ወደ መዝጋቢ መለስ"` `en="Return to Registrar"`
  - Opens `ReturnModal`
  - On modal confirm: calls `reviewReturn(recordType, recordId, comment)` from `src/api/workflow.ts`
  - On success: same invalidations + toast + navigate to `/pending`
  - The comment in the modal body MUST be non-empty (enforced inside `ReturnModal`)

**Details tab — read-only field display:**
Create a simple `ReadField` sub-component:
```typescript
function ReadField({ labelAm, labelEn, value }: { labelAm: string; labelEn: string; value: React.ReactNode }) {
  if (value === null || value === undefined || value === '' || (Array.isArray(value) && value.length === 0)) return null
  return (
    <div>
      <div className="text-xs text-muted-foreground">
        <span className="font-amharic">{labelAm}</span> / {labelEn}
      </div>
      <div className="font-amharic mt-0.5 text-sm text-foreground">{
        Array.isArray(value) ? value.join(', ') : value
      }</div>
    </div>
  )
}
```

Show these field groups in the details tab (skip any field where the value is null/undefined/empty):
- Identity: name_amharic, name_local, category, current_use, previous_id
- Location: woreda, kebele, house_number, street_number, gate
- GPS: gps_east, gps_north, elevation_m
- Ownership: owner_type, owner_name, map_reference
- History: built_by, construction_period, age_method
- Dimensions: height_m, length_m, width_m, num_doors, num_windows, num_rooms, material
- Classification: harari_house_grades, neighborhood_type, description
- Condition: overall_condition, damage_roof, damage_cornice, damage_wall, damage_floor, damage_door, damage_cupboard, damage_upper_floor, damage_dera, damage_pillar
- Values: value_historical, value_craftsmanship, value_artistic, value_scientific, value_cultural
- Conservation: has_threat, maintenance_done, maintenance_reason, maintenance_by, maintenance_count, preventive_level, accessibility, notes
- Informant: caretaker_name, caretaker_role, informant_name, informant_sex, informant_age, registrar_date

(For movable records, map the equivalent movable fields.)

**Photos tab:** Show photos in a 2-3 column responsive grid, read-only. Photo src: if `file_path` starts with `http` use as-is, otherwise prepend `import.meta.env.VITE_API_URL + '/' + file_path.replace(/^\//, '')`.

---

### 6. Update `/src/routes/_authenticated/pending.tsx`

Replace the EmptyState placeholder with the real pending queue:

```typescript
import { createFileRoute, Navigate, Link, useNavigate } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/authStore'
import { useTranslation } from 'react-i18next'
import { RecordsList } from '@/components/Records/RecordsList'
import { listRecords } from '@/api/records'
import type { UnifiedListParams } from '@/types'

export const Route = createFileRoute('/_authenticated/pending')({
  component: PendingPage,
})

function PendingPage() {
  const { t } = useTranslation()
  const role = useAuthStore((s) => s.user?.role)
  if (role && role !== 'supervisor') return <Navigate to="/unauthorized" replace />

  // Fetch records with status=pending_review
  const fetcher = (params: UnifiedListParams) =>
    listRecords({ ...params, status: 'pending_review' })

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t('nav.pendingReview')}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t('supervisor.pendingSubtitle')}
        </p>
      </div>
      <RecordsList
        queryKey={['records', 'pending_review']}
        fetcher={fetcher}
        // Override card link to go to supervisor detail, not edit
        detailPath="/supervisor/records"
      />
    </div>
  )
}
```

**Note:** `RecordsList` currently links cards to `/records/$id/edit`. Supervisor should link to a read-only detail view instead. Update `RecordsList` (or `RecordCard`) to accept an optional `detailPath` prop. When provided, card links to `${detailPath}/${record.record_type}/${record.id}` instead of the default edit path.

---

### 7. New route: `/src/routes/_authenticated/supervisor.records.$type.$id.tsx`

This is the supervisor read-only detail page:

```typescript
import { createFileRoute, Navigate, Link } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/authStore'
import { useTranslation } from 'react-i18next'
import { ArrowLeft } from 'lucide-react'
import { SupervisorRecordDetail } from '@/components/Records/SupervisorRecordDetail'
import type { RecordType } from '@/types'

export const Route = createFileRoute('/_authenticated/supervisor/records/$type/$id')({
  component: SupervisorRecordDetailPage,
})

function SupervisorRecordDetailPage() {
  const { type, id } = Route.useParams()
  const { t } = useTranslation()
  const role = useAuthStore((s) => s.user?.role)

  if (role && role !== 'supervisor') return <Navigate to="/unauthorized" replace />
  if (type !== 'immovable' && type !== 'movable') return <Navigate to="/pending" replace />

  return (
    <div className="space-y-4">
      <Link
        to="/pending"
        className="font-amharic inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        {t('nav.pendingReview')}
      </Link>
      <SupervisorRecordDetail recordId={id} recordType={type as RecordType} />
    </div>
  )
}
```

---

### 8. New route: `/src/routes/_authenticated/supervisor.records.index.tsx`

All records list for supervisor (can see all statuses, not just pending):

```typescript
import { createFileRoute, Navigate } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/authStore'
import { useTranslation } from 'react-i18next'
import { RecordsList } from '@/components/Records/RecordsList'
import { listRecords } from '@/api/records'

export const Route = createFileRoute('/_authenticated/supervisor/records/')({
  component: SupervisorAllRecordsPage,
})

function SupervisorAllRecordsPage() {
  const { t } = useTranslation()
  const role = useAuthStore((s) => s.user?.role)
  if (role && role !== 'supervisor') return <Navigate to="/unauthorized" replace />

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t('nav.allRecords')}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t('supervisor.allSubtitle')}
        </p>
      </div>
      <RecordsList
        queryKey={['records', 'all']}
        fetcher={listRecords}
        detailPath="/supervisor/records"
      />
    </div>
  )
}
```

---

### 9. Update `src/components/Layout/Sidebar.tsx`

Update the NAV array so supervisor "All Records" links to `/supervisor/records` and "Pending Review" links to `/pending`:

```typescript
const NAV: NavItem[] = [
  { to: '/dashboard',           labelKey: 'nav.dashboard',     icon: LayoutDashboard, roles: ['registrar', 'supervisor', 'manager'] },
  { to: '/records',             labelKey: 'nav.myRecords',     icon: FileText,        roles: ['registrar'] },
  { to: '/supervisor/records',  labelKey: 'nav.allRecords',    icon: FileText,        roles: ['supervisor'] },
  { to: '/pending',             labelKey: 'nav.pendingReview', icon: ListChecks,      roles: ['supervisor'] },
  { to: '/manager/records',     labelKey: 'nav.allRecords',    icon: FileText,        roles: ['manager'] },   // placeholder for prompt 6
  { to: '/reviewed',            labelKey: 'nav.reviewed',      icon: CheckCircle2,    roles: ['manager'] },
  { to: '/users',               labelKey: 'nav.users',         icon: Users,           roles: ['manager'] },
  { to: '/reports',             labelKey: 'nav.reports',       icon: BarChart3,       roles: ['manager'] },
]
```

Also add a **pending count badge** next to "Pending Review" in the sidebar for supervisors. Fetch it from `useAuthStore`'s user or from the dashboard stats query if already loaded. Use `pending_my_action` from `DashboardStats` when available. Show it as a small amber pill: `<span className="ml-auto rounded-full bg-amber-500 px-1.5 py-0.5 text-[10px] font-bold text-white">{count}</span>`.

---

### 10. i18n additions

Add to **both** `en.json` and `am.json`:

**en.json:**
```json
{
  "supervisor": {
    "pendingSubtitle": "Records waiting for your review",
    "allSubtitle": "All heritage records across all registrars",
    "approveSuccess": "Record approved and moved to final review",
    "returnSuccess": "Record returned to registrar",
    "actionBar": "Review Actions"
  },
  "actions": {
    "approve": "Approve",
    "returnToRegistrar": "Return to Registrar",
    "returnToSupervisor": "Return to Supervisor",
    "finalApprove": "Final Approve",
    "viewDetail": "View Detail"
  },
  "toast": {
    "approveSuccess": "Record approved successfully",
    "returnSuccess": "Record returned successfully",
    "error": "Action failed — please try again"
  }
}
```

**am.json:**
```json
{
  "supervisor": {
    "pendingSubtitle": "ለግምገማዎ የሚጠብቁ መዛግብት",
    "allSubtitle": "ከሁሉም መዝጋቢዎች ያሉ ሁሉም መዛግብት",
    "approveSuccess": "መዝገቡ ፀድቆ ወደ ሥራ አስኪያጅ ተልኳል",
    "returnSuccess": "መዝገቡ ወደ መዝጋቢ ተመልሷል",
    "actionBar": "የግምገማ እርምጃዎች"
  },
  "actions": {
    "approve": "አፅድቅ",
    "returnToRegistrar": "ወደ መዝጋቢ መለስ",
    "returnToSupervisor": "ወደ ተቆጣጣሪ መለስ",
    "finalApprove": "የመጨረሻ ፈቃድ",
    "viewDetail": "ዝርዝር ይመልከቱ"
  },
  "toast": {
    "approveSuccess": "መዝገቡ በተሳካ ሁኔታ ፀድቋል",
    "returnSuccess": "መዝገቡ በተሳካ ሁኔታ ተመልሷል",
    "error": "እርምጃው አልተሳካም — እንደገና ይሞክሩ"
  }
}
```

---

### 11. Toast setup

The project has `sonner` installed. Make sure it is wired in. Check `src/routes/__root.tsx` — if `<Toaster />` from `sonner` is not already rendered inside `RootComponent`, add it:

```typescript
import { Toaster } from 'sonner'

function RootComponent() {
  const { queryClient } = Route.useRouteContext()
  return (
    <QueryClientProvider client={queryClient}>
      <LanguageSync />
      <Outlet />
      <Toaster richColors position="top-right" />
    </QueryClientProvider>
  )
}
```

Use `toast.success(t('toast.approveSuccess'))` and `toast.error(t('toast.error'))` from `sonner` in the action handlers.

---

## IMPORTANT RULES — DO NOT VIOLATE

1. **Supervisors NEVER see edit controls.** `SupervisorRecordDetail` is read-only. No Save Draft, no Submit button. Only Approve and Return actions — and only when status is `pending_review`.

2. **`reviewApprove` body** — pass `{}` if no comment, or `{ comment }` if one is provided. The endpoint accepts an optional comment.

3. **`reviewReturn` body** — `{ comment: string }`. The `ReturnModal` enforces non-empty. The `comment` field name here is different from `addComment` which uses `comment_text`. Do not mix them up.

4. **Query invalidation after workflow actions** — always invalidate:
   - `['immovable', recordId]` or `['movable', recordId]`
   - `['dashboard', 'stats']`
   - `['records']` (the unified list)

5. **RecordCard link override** — when `detailPath` is passed to `RecordsList`/`RecordCard`, the card href must be `${detailPath}/${record.record_type}/${record.id}`. This gives URLs like `/supervisor/records/immovable/uuid-here`.

6. **Do not touch** `src/types/index.ts`, `src/stores/authStore.ts`, `src/routes/__root.tsx` (except adding `<Toaster />`), or `src/routes/_authenticated.tsx`.

7. **TanStack Router** — new route files must follow the existing naming convention. The router plugin auto-generates `routeTree.gen.ts` on build/dev — never edit it manually.
