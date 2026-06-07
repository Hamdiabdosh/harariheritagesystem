# CURSOR PROMPT 6 — Qirs Mezgeb
# Manager Flow: Final Review Queue + Record Detail + User Management + Reports + CSV Export
# Paste this entire file into Cursor Chat with the project open.

---

## CONTEXT

This is the final prompt. The project now has:
- `src/api/workflow.ts` — `finalApprove`, `finalReturn`, `getComments`, `addComment`, `getHistory`
- `src/api/users.ts` — `listUsers`, `createUser`, `updateUser`, `deactivateUser`
- `src/components/Dashboard/ManagerDashboard.tsx` — built in Prompt 5 (stat cards + CTA)
- `src/components/Records/SupervisorRecordDetail.tsx` — the supervisor detail pattern to follow
- `src/components/UI/ReturnModal.tsx`, `StatusTimeline`, `CommentThread` — all built
- `src/routes/_authenticated/reviewed.tsx` — currently a placeholder EmptyState
- `src/routes/_authenticated/users.tsx` — currently a placeholder EmptyState
- `src/routes/_authenticated/reports.tsx` — currently a placeholder EmptyState

Reuse all existing shared components: `RecordsList`, `RecordCard`, `StatusBadge`, `LoadingSpinner`,
`EmptyState`, `ReturnModal`, `StatusTimeline`, `CommentThread`, `Field`, shadcn/ui `Dialog`,
`Tabs`, `Table`, `Badge`, `Button`, `Input`, `Select`, `Textarea`.

---

## WHAT TO BUILD

### 1. `/src/api/export.ts`

Create this file:

```typescript
import { apiClient } from './client'
import type { UnifiedListParams } from '@/types'

// GET /export/records/csv  [supervisor + manager only]
// Returns a blob — use responseType: 'blob'
export async function exportCSV(params: UnifiedListParams = {}): Promise<Blob> {
  const response = await apiClient.get('/export/records/csv', {
    params,
    responseType: 'blob',
  })
  return response.data as Blob
}

// GET /records/:type/:id/pdf  [all roles]
export async function exportPDF(
  recordType: 'immovable' | 'movable',
  recordId: string,
): Promise<Blob> {
  const response = await apiClient.get(`/records/${recordType}/${recordId}/pdf`, {
    responseType: 'blob',
  })
  return response.data as Blob
}

// Helper: trigger browser file download from a Blob
export function downloadBlob(blob: Blob, filename: string): void {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}
```

---

### 2. `/src/components/Records/ManagerRecordDetail.tsx`

Create this component — it is the manager's read-only detail view, following the exact same
pattern as `SupervisorRecordDetail.tsx` but with final-approve / final-return actions.

**Props:**
```typescript
interface ManagerRecordDetailProps {
  recordId: string
  recordType: RecordType
}
```

**Data fetching:** identical to `SupervisorRecordDetail` — try immovable first, then movable.

**Layout:** identical to `SupervisorRecordDetail` — header row, action bar, tabbed panel.

**Action bar logic — only shown when `record.status === 'under_review'`:**

- **Final Approve button**:
  - Label: `am="የመጨረሻ ፈቃድ"` `en="Final Approve"`
  - Green/primary, `CheckCircle2` icon
  - Calls `finalApprove(recordType, recordId)` from `src/api/workflow.ts`
  - On success: invalidate `['immovable'/'movable', recordId]` + `['dashboard', 'stats']` + `['records']`
  - Show `toast.success(t('toast.finalApproveSuccess'))`, navigate to `/reviewed`

- **Return to Supervisor button**:
  - Label: `am="ወደ ተቆጣጣሪ መለስ"` `en="Return to Supervisor"`
  - Outline/destructive variant
  - Opens `ReturnModal`
  - On confirm: calls `finalReturn(recordType, recordId, comment)` from `src/api/workflow.ts`
  - On success: same invalidations + `toast.success(t('toast.returnSuccess'))` + navigate to `/reviewed`

- **Download PDF button** (shown for ALL statuses except draft):
  - Label: `am="PDF አውርድ"` `en="Download PDF"`
  - Outline variant, `Download` icon
  - Calls `exportPDF(recordType, recordId)` then `downloadBlob(blob, record.record_id + '.pdf')`
  - Shows spinner while loading
  - On error: `toast.error(t('toast.error'))`

**Tabs:** Details / History / Comments / Photos — identical to SupervisorRecordDetail.
`CommentThread` with `canComment={true}` (managers can always comment).

---

### 3. New route: `/src/routes/_authenticated/manager.records.$type.$id.tsx`

```typescript
import { createFileRoute, Navigate, Link } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/authStore'
import { useTranslation } from 'react-i18next'
import { ArrowLeft } from 'lucide-react'
import { ManagerRecordDetail } from '@/components/Records/ManagerRecordDetail'
import type { RecordType } from '@/types'

export const Route = createFileRoute('/_authenticated/manager/records/$type/$id')({
  component: ManagerRecordDetailPage,
})

function ManagerRecordDetailPage() {
  const { type, id } = Route.useParams()
  const { t } = useTranslation()
  const role = useAuthStore((s) => s.user?.role)

  if (role && role !== 'manager') return <Navigate to="/unauthorized" replace />
  if (type !== 'immovable' && type !== 'movable') return <Navigate to="/reviewed" replace />

  return (
    <div className="space-y-4">
      <Link
        to="/reviewed"
        className="font-amharic inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        {t('nav.reviewed')}
      </Link>
      <ManagerRecordDetail recordId={id} recordType={type as RecordType} />
    </div>
  )
}
```

---

### 4. New route: `/src/routes/_authenticated/manager.records.index.tsx`

Manager all-records list (all statuses, all registrars):

```typescript
import { createFileRoute, Navigate } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/authStore'
import { useTranslation } from 'react-i18next'
import { RecordsList } from '@/components/Records/RecordsList'
import { listRecords } from '@/api/records'

export const Route = createFileRoute('/_authenticated/manager/records/')({
  component: ManagerAllRecordsPage,
})

function ManagerAllRecordsPage() {
  const { t } = useTranslation()
  const role = useAuthStore((s) => s.user?.role)
  if (role && role !== 'manager') return <Navigate to="/unauthorized" replace />

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t('nav.allRecords')}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t('manager.allSubtitle')}
        </p>
      </div>
      <RecordsList
        queryKey={['records', 'manager-all']}
        fetcher={listRecords}
        detailPath="/manager/records"
      />
    </div>
  )
}
```

---

### 5. Update `/src/routes/_authenticated/reviewed.tsx`

Replace the placeholder with the real final-review queue:

```typescript
import { createFileRoute, Navigate } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/authStore'
import { useTranslation } from 'react-i18next'
import { RecordsList } from '@/components/Records/RecordsList'
import { listRecords } from '@/api/records'
import type { UnifiedListParams } from '@/types'

export const Route = createFileRoute('/_authenticated/reviewed')({
  component: ReviewedPage,
})

function ReviewedPage() {
  const { t } = useTranslation()
  const role = useAuthStore((s) => s.user?.role)
  if (role && role !== 'manager') return <Navigate to="/unauthorized" replace />

  // Fetch records with status=under_review (ready for manager's final decision)
  const fetcher = (params: UnifiedListParams) =>
    listRecords({ ...params, status: 'under_review' })

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t('nav.reviewed')}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t('manager.reviewedSubtitle')}
        </p>
      </div>
      <RecordsList
        queryKey={['records', 'under_review']}
        fetcher={fetcher}
        detailPath="/manager/records"
      />
    </div>
  )
}
```

---

### 6. `/src/components/Users/UserManagement.tsx`

Create a full user management component for the manager.

**Layout:**

```
┌──────────────────────────────────────────────────────┐
│  ተጠቃሚዎች / Users              [+ New User button]    │
├──────────────────────────────────────────────────────┤
│  [Filter by role dropdown]  [Active/All toggle]       │
├──────────────────────────────────────────────────────┤
│  TABLE:                                               │
│  Full Name | Email | Role | Status | Actions          │
│  ──────────────────────────────────────────────────  │
│  Abdulhamid │ a@b.c │ Manager │ Active │ [Edit][✕]   │
│  ...                                                  │
└──────────────────────────────────────────────────────┘
```

**Data fetching:**
- `useQuery` with key `['users', filters]`, calls `listUsers(filters)` from `src/api/users.ts`
- Filter state: `role` (all / registrar / supervisor / manager), `is_active` (true / false / undefined)
- Pagination: same prev/next pattern as `RecordsList`

**Table:** Use shadcn/ui `Table` from `src/components/ui/table.tsx`.
Columns: Full Name, Email, Role (use `StatusBadge`-style pill with role colors), Status (Active/Inactive badge), Actions.

**Role badge colors:**
- registrar → blue (`bg-blue-100 text-blue-900`)
- supervisor → amber (`bg-amber-100 text-amber-900`)
- manager → purple (`bg-purple-100 text-purple-900`)

**Status badge:**
- Active → emerald (`bg-emerald-100 text-emerald-900`)
- Inactive → muted/gray

**Actions column:**
- **Edit** button (`Pencil` icon, outline small): opens `UserFormModal` in edit mode
- **Deactivate** button (`UserX` icon, destructive small): shows `ConfirmDeactivateModal`
  - Disabled if `user.id === currentUser.id` (cannot deactivate yourself)
  - On confirm: calls `deactivateUser(id)`, invalidates `['users']`, shows toast

**"+ New User" button:** opens `UserFormModal` in create mode.

---

### 7. `/src/components/Users/UserFormModal.tsx`

Create a modal using shadcn/ui `Dialog` for both create and edit:

**Props:**
```typescript
interface UserFormModalProps {
  open: boolean
  onClose: () => void
  user?: UserItem   // undefined = create mode, defined = edit mode
}
```

**Fields (all using native inputs — no React Hook Form needed, keep it simple):**
- `full_name` — TextInput, required
- `email` — email input, required (only shown in create mode; read-only in edit)
- `password` — password input, required in create, hidden in edit
- `role` — Select: registrar / supervisor / manager (bilingual options)
- `language` — Select: am / en
- `is_active` — Switch (only shown in edit mode)

**Validation (client-side):**
- `full_name`: required, max 100 chars
- `email`: required, must contain `@`
- `password`: required in create mode, min 8 chars
- `role`: required

**On submit:**
- Create mode: calls `createUser(body)` → `toast.success` → close modal → invalidate `['users']`
- Edit mode: calls `updateUser(user.id, body)` → `toast.success` → close modal → invalidate `['users']`
- Shows inline error if API returns 409 (email already exists)
- Shows `Loader2` spinner on submit button while pending

**Use the same `Field`, `TextInput`, `Select`, `Switch` primitives from `src/components/Forms/Field.tsx`.**

---

### 8. `/src/components/Users/ConfirmDeactivateModal.tsx`

Simple confirmation dialog:

```typescript
interface ConfirmDeactivateModalProps {
  open: boolean
  onClose: () => void
  onConfirm: () => Promise<void>
  userName: string
}
```

- Uses shadcn/ui `Dialog`
- Body: `am="[userName]ን ለማቦዘን ይፈልጋሉ?"` `en="Deactivate [userName]? They will lose access immediately."`
- Cancel + Confirm (destructive) buttons
- Spinner on confirm while pending

---

### 9. Update `/src/routes/_authenticated/users.tsx`

Replace placeholder with real user management:

```typescript
import { createFileRoute, Navigate } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/authStore'
import { useTranslation } from 'react-i18next'
import { UserManagement } from '@/components/Users/UserManagement'

export const Route = createFileRoute('/_authenticated/users')({
  component: UsersPage,
})

function UsersPage() {
  const { t } = useTranslation()
  const role = useAuthStore((s) => s.user?.role)
  if (role && role !== 'manager') return <Navigate to="/unauthorized" replace />
  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t('nav.users')}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t('manager.usersSubtitle')}
        </p>
      </div>
      <UserManagement />
    </div>
  )
}
```

---

### 10. `/src/components/Reports/ReportsPage.tsx`

Create the reports/export component:

**Layout:**

```
┌──────────────────────────────────────────────────────┐
│  ሪፖርቶች / Reports                                    │
├──────────────────────────────────────────────────────┤
│  EXPORT SECTION                                       │
│  ┌────────────────────────────────────────────────┐  │
│  │  ሁሉንም ቅርሶች ወደ CSV ላክ / Export All to CSV     │  │
│  │                                                │  │
│  │  Filters:                                      │  │
│  │  [Type: All/Immovable/Movable]                 │  │
│  │  [Status: All/Draft/Pending.../Approved]       │  │
│  │  [Woreda: text input]                          │  │
│  │  [Date from: date] [Date to: date]             │  │
│  │                                                │  │
│  │  [Download CSV]  ← calls exportCSV(filters)   │  │
│  └────────────────────────────────────────────────┘  │
├──────────────────────────────────────────────────────┤
│  SUMMARY STATS SECTION                               │
│  Reads from GET /dashboard/stats                     │
│                                                      │
│  Total Immovable: 124    Total Movable: 47           │
│                                                      │
│  Status breakdown bar chart (use simple CSS bars,    │
│  no recharts needed — just colored div widths):      │
│                                                      │
│  Draft        ████░░░░░░░░  42 (25%)                │
│  Pending      ██░░░░░░░░░░  18 (11%)                │
│  Under Review █░░░░░░░░░░░   9  (5%)                │
│  Returned     ██░░░░░░░░░░  12  (7%)                │
│  Approved     ████████░░░░  90 (53%)                │
└──────────────────────────────────────────────────────┘
```

**CSV export behavior:**
- Filter state: `type`, `status`, `woreda`, `date_from`, `date_to`
- Download button: calls `exportCSV(filters)` from `src/api/export.ts`
- Use `downloadBlob(blob, 'qirs-mezgeb-export.csv')` to trigger download
- Shows `Loader2` spinner while pending
- Error: `toast.error(t('toast.error'))`

**Stats chart:**
- Fetch from `GET /dashboard/stats` — query key `['dashboard', 'stats']`
- Total = `total_immovable + total_movable`
- Bar width: `(count / total * 100).toFixed(1) + '%'`
- Colors match `StatusBadge`: draft=gray, pending=amber, under_review=blue, returned=rose, approved=emerald

---

### 11. Update `/src/routes/_authenticated/reports.tsx`

```typescript
import { createFileRoute, Navigate } from '@tanstack/react-router'
import { useAuthStore } from '@/stores/authStore'
import { ReportsPage } from '@/components/Reports/ReportsPage'

export const Route = createFileRoute('/_authenticated/reports')({
  component: Reports,
})

function Reports() {
  const role = useAuthStore((s) => s.user?.role)
  // Supervisors can also export CSV — allow both roles
  if (role && role !== 'manager' && role !== 'supervisor') {
    return <Navigate to="/unauthorized" replace />
  }
  return <ReportsPage />
}
```

Also update `Sidebar.tsx` to include Reports for supervisor too:
```typescript
{ to: '/reports', labelKey: 'nav.reports', icon: BarChart3, roles: ['supervisor', 'manager'] },
```

---

### 12. Update `src/components/Layout/Topbar.tsx`

Add a **PDF download button** to the topbar only when the user is viewing a detail page.
Actually — skip this. PDF download is already in `ManagerRecordDetail`. The topbar stays as-is.

Instead, add a **language toggle** to the topbar if it does not already have one:
- Button showing current language (`አማ` / `EN`)
- On click: calls `useLanguageStore((s) => s.setLanguage)` to toggle between `'am'` and `'en'`
- Also calls `updateMyLanguage(newLang)` from `src/api/users.ts` to persist to the server

---

### 13. Update `src/components/Layout/Sidebar.tsx`

Add the **under_review count badge** next to "Reviewed" for managers, matching the amber badge
already added for supervisors in Prompt 5. Fetch from dashboard stats:

```typescript
// Inside Sidebar, for manager role:
const { data: stats } = useQuery({
  queryKey: ['dashboard', 'stats'],
  queryFn: getDashboardStats,
  enabled: role === 'supervisor' || role === 'manager',
  staleTime: 60_000,
})

// For supervisor: show stats.by_status.pending_review next to /pending
// For manager: show stats.by_status.under_review next to /reviewed
```

Import `getDashboardStats` from `src/api/dashboard.ts`.

---

### 14. i18n additions

Add to **both** `en.json` and `am.json`:

**en.json:**
```json
{
  "manager": {
    "allSubtitle": "All heritage records — full access",
    "reviewedSubtitle": "Records reviewed by supervisor, awaiting your final decision",
    "usersSubtitle": "Manage registrar and supervisor accounts",
    "approveSuccess": "Record fully approved and locked",
    "returnSuccess": "Record returned to supervisor for revision"
  },
  "users": {
    "newUser": "New User",
    "fullName": "Full Name",
    "email": "Email",
    "password": "Password",
    "role": "Role",
    "language": "Language",
    "isActive": "Active",
    "createTitle": "Create New User",
    "editTitle": "Edit User",
    "confirmDeactivate": "Deactivate User",
    "deactivateWarning": "This user will immediately lose access.",
    "emailExists": "A user with this email already exists.",
    "filterRole": "Filter by role",
    "filterAll": "All roles",
    "active": "Active",
    "inactive": "Inactive",
    "allUsers": "All users"
  },
  "reports": {
    "title": "Reports & Export",
    "exportSection": "Export to CSV",
    "exportButton": "Download CSV",
    "exporting": "Exporting...",
    "statsSection": "Summary Statistics",
    "totalRecords": "Total Records",
    "filterType": "Record type",
    "filterStatus": "Status",
    "filterWoreda": "Woreda",
    "filterDateFrom": "From date",
    "filterDateTo": "To date"
  },
  "toast": {
    "finalApproveSuccess": "Record fully approved",
    "createUserSuccess": "User created successfully",
    "updateUserSuccess": "User updated successfully",
    "deactivateSuccess": "User deactivated"
  }
}
```

**am.json:**
```json
{
  "manager": {
    "allSubtitle": "ሁሉም የቅርስ መዛግብት — ሙሉ መዳረሻ",
    "reviewedSubtitle": "በተቆጣጣሪ የተገመገሙ፣ የእርስዎን የመጨረሻ ውሳኔ የሚጠብቁ",
    "usersSubtitle": "የመዝጋቢ እና ተቆጣጣሪ መለያዎችን ያስተዳድሩ",
    "approveSuccess": "መዝገቡ ሙሉ በሙሉ ፀድቆ ተቆልፏል",
    "returnSuccess": "መዝገቡ ወደ ተቆጣጣሪ ለማሻሻያ ተመልሷል"
  },
  "users": {
    "newUser": "አዲስ ተጠቃሚ",
    "fullName": "ሙሉ ስም",
    "email": "ኢሜይል",
    "password": "የይለፍ ቃል",
    "role": "ሚና",
    "language": "ቋንቋ",
    "isActive": "ንቁ",
    "createTitle": "አዲስ ተጠቃሚ ፍጠር",
    "editTitle": "ተጠቃሚ አርትዕ",
    "confirmDeactivate": "ተጠቃሚን አቦዝን",
    "deactivateWarning": "ይህ ተጠቃሚ ወዲያው መዳረሻ ያጣል።",
    "emailExists": "ይህ ኢሜይል ያለው ተጠቃሚ አስቀድሞ ይገኛል።",
    "filterRole": "በሚና አጣራ",
    "filterAll": "ሁሉም ሚናዎች",
    "active": "ንቁ",
    "inactive": "ንቁ ያልሆነ",
    "allUsers": "ሁሉም ተጠቃሚዎች"
  },
  "reports": {
    "title": "ሪፖርቶች እና ላክ",
    "exportSection": "ወደ CSV ላክ",
    "exportButton": "CSV አውርድ",
    "exporting": "በመላክ ላይ...",
    "statsSection": "ጠቅላላ ስታቲስቲክስ",
    "totalRecords": "ጠቅላላ መዛግብት",
    "filterType": "የቅርስ ዓይነት",
    "filterStatus": "ሁኔታ",
    "filterWoreda": "ወረዳ",
    "filterDateFrom": "ከቀን",
    "filterDateTo": "እስከ ቀን"
  },
  "toast": {
    "finalApproveSuccess": "መዝገቡ ሙሉ በሙሉ ፀድቋል",
    "createUserSuccess": "ተጠቃሚ በተሳካ ሁኔታ ተፈጥሯል",
    "updateUserSuccess": "ተጠቃሚ በተሳካ ሁኔታ ተዘምኗል",
    "deactivateSuccess": "ተጠቃሚ ታቦዝኗል"
  }
}
```

---

## FINAL WIRING CHECKLIST — do all of these

### Route fixes
- Confirm `src/routes/_authenticated/reviewed.tsx` now fetches `status=under_review` (not just a stub)
- Confirm `src/routes/_authenticated/users.tsx` renders `<UserManagement />`
- Confirm `src/routes/_authenticated/reports.tsx` renders `<ReportsPage />` (accessible to supervisor + manager)
- Confirm `src/routes/_authenticated/manager.records.$type.$id.tsx` exists and role-guards to manager only
- Confirm `src/routes/_authenticated/manager.records.index.tsx` exists

### Sidebar final state
After this prompt, the complete NAV array in `Sidebar.tsx` should be:

```typescript
const NAV: NavItem[] = [
  // All roles
  { to: '/dashboard',          labelKey: 'nav.dashboard',     icon: LayoutDashboard, roles: ['registrar','supervisor','manager'] },

  // Registrar only
  { to: '/records',            labelKey: 'nav.myRecords',     icon: FileText,        roles: ['registrar'] },

  // Supervisor
  { to: '/supervisor/records', labelKey: 'nav.allRecords',    icon: FileText,        roles: ['supervisor'] },
  { to: '/pending',            labelKey: 'nav.pendingReview', icon: ListChecks,      roles: ['supervisor'], showCount: 'pending_review' },

  // Manager
  { to: '/manager/records',    labelKey: 'nav.allRecords',    icon: FileText,        roles: ['manager'] },
  { to: '/reviewed',           labelKey: 'nav.reviewed',      icon: CheckCircle2,    roles: ['manager'], showCount: 'under_review' },
  { to: '/users',              labelKey: 'nav.users',         icon: Users,           roles: ['manager'] },

  // Supervisor + Manager
  { to: '/reports',            labelKey: 'nav.reports',       icon: BarChart3,       roles: ['supervisor','manager'] },
]
```

Add a `showCount` optional field to the `NavItem` type. When present, look up `stats.by_status[showCount]` and render the count badge (amber pill) next to the label.

### RecordCard link behavior — final summary
| Role       | Clicks a card in list | Navigates to |
|------------|-----------------------|--------------|
| registrar  | any list              | `/records/${id}/edit` (default) |
| supervisor | pending queue / all records | `/supervisor/records/${type}/${id}` |
| manager    | reviewed queue / all records | `/manager/records/${type}/${id}` |

This is controlled by the `detailPath` prop on `RecordsList` — confirm it is correctly threaded through to `RecordCard`'s `Link` href.

---

## IMPORTANT RULES — DO NOT VIOLATE

1. **Managers NEVER see the registrar form edit controls.** `ManagerRecordDetail` is purely
   read-only. The only write actions are Final Approve, Final Return, and Add Comment.

2. **`finalReturn` body** — `{ comment: string }` (same as `reviewReturn`). The `ReturnModal`
   enforces non-empty. Do not change field names.

3. **`finalApprove` body** — `{}` or `{ comment }`. Optional.

4. **PDF export** — use `responseType: 'blob'` in the Axios call. Never try to parse the response
   as JSON. `downloadBlob` handles the browser download trigger.

5. **CSV export** — same blob handling. The server sets `Content-Disposition: attachment` with
   a filename; the frontend can use `'qirs-mezgeb-export.csv'` as the fallback filename.

6. **Cannot deactivate self** — the server returns 403 for this; also disable the button
   client-side when `user.id === useAuthStore.getState().user?.id`.

7. **Query invalidation after workflow actions:**
   Always invalidate: `['immovable'/'movable', recordId]` + `['dashboard','stats']` + `['records']`
   After user mutations: always invalidate `['users']`

8. **Do not touch** `src/types/index.ts`, `src/stores/authStore.ts`,
   `src/stores/languageStore.ts`, `src/routes/_authenticated.tsx`.

9. **TanStack Router** — new route files follow the existing naming pattern.
   Never manually edit `routeTree.gen.ts`.

10. **No hardcoded strings** — every visible string uses `t('key')` from i18next.
