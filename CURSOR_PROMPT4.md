# CURSOR PROMPT 4 — Qirs Mezgeb
# Movable Form + RecordDetail (StatusTimeline + CommentThread)
# Paste this entire file into Cursor Chat (with the project open).

---

## CONTEXT

This is a TanStack Start project. The existing code uses:
- TanStack Router file-based routes under `src/routes/`
- TanStack Query v5 for data fetching
- Zustand for auth (`useAuthStore`) and language (`useLanguageStore`) stores
- Axios client at `src/api/client.ts` with automatic Bearer token + silent refresh
- Bilingual Amharic/English i18n via `react-i18next` (keys in `src/i18n/am.json` + `src/i18n/en.json`)
- Tailwind CSS + shadcn/ui components in `src/components/ui/`
- Custom form primitives: `Field`, `FormSection`, `TextInput`, `NumberInput`, `TextArea`, `Select`, `CheckboxList`, `Switch` — all in `frontend/src/components/forms/Field.tsx`
- `StatusBadge` at `frontend/src/components/common/StatusBadge.tsx`
- `RecordCard`, `RecordsList` at `frontend/src/components/records/`
- Photo upload at `frontend/src/components/forms/PhotoUploader.tsx` — currently wired only to immovable

The `ImmovableForm` in `frontend/src/components/forms/ImmovableForm.tsx` is the reference pattern for all form work — follow it exactly.

The route `src/routes/_authenticated/records.new.$type.tsx` already renders `<ImmovableForm>` for `type === "immovable"` and shows a "coming soon" placeholder for `type === "movable"`. That placeholder must be replaced.

The route `src/routes/_authenticated/records.$id.edit.tsx` currently only handles immovable records. It must be extended to detect the record type and render the right form + full detail view.

---

## WHAT TO BUILD

### 1. `/src/api/movable.ts` — Movable API functions

Create this file with these exact functions, following the same pattern as `src/api/immovable.ts`:

```typescript
import { apiClient } from './client'
import type {
  ApiResponse, MovableCreateResult, MovableRecordDetail,
  MovableRecordInput, PaginatedMovable, RecordListParams,
  RecordPhoto, SubmitResult,
} from '@/types'

// GET /records/movable
export async function listMovable(params: RecordListParams = {}): Promise<PaginatedMovable> {
  const { data } = await apiClient.get<ApiResponse<PaginatedMovable>>('/records/movable', { params })
  return data.data
}

// GET /records/movable/:id  → returns { record, photos, comments, history }
export async function getMovable(id: string): Promise<MovableRecordDetail> {
  const { data } = await apiClient.get<ApiResponse<MovableRecordDetail>>(`/records/movable/${id}`)
  return data.data
}

// POST /records/movable
export async function createMovable(body: MovableRecordInput): Promise<MovableCreateResult> {
  const { data } = await apiClient.post<ApiResponse<MovableCreateResult>>('/records/movable', body)
  return data.data
}

// PUT /records/movable/:id   — partial patch, returns { record: MovableRecord }
export async function updateMovable(id: string, body: MovableRecordInput): Promise<void> {
  await apiClient.put(`/records/movable/${id}`, body)
}

// PUT /records/movable/:id/submit
export async function submitMovable(id: string): Promise<SubmitResult> {
  const { data } = await apiClient.put<ApiResponse<SubmitResult>>(`/records/movable/${id}/submit`)
  return data.data
}

// POST /records/:type/:id/photos  — multipart, field name "photo"
export async function uploadMovablePhoto(id: string, file: File): Promise<RecordPhoto> {
  const fd = new FormData()
  fd.append('photo', file)
  const { data } = await apiClient.post<ApiResponse<RecordPhoto>>(
    `/records/movable/${id}/photos`, fd,
    { headers: { 'Content-Type': 'multipart/form-data' } }
  )
  return data.data
}

// DELETE /records/:type/:id/photos/:photo_id
export async function deleteMovablePhoto(recordId: string, photoId: string): Promise<void> {
  await apiClient.delete(`/records/movable/${recordId}/photos/${photoId}`)
}
```

---

### 2. `/src/api/workflow.ts` — Workflow API functions

Create this file:

```typescript
import { apiClient } from './client'
import type {
  ApiResponse, RecordType, RecordComment, StatusHistoryEntry,
  AddCommentBody, ApproveBody, ReturnBody,
} from '@/types'

// PUT /records/:type/:id/review-approve  [supervisor only]
export async function reviewApprove(type: RecordType, id: string, comment?: string): Promise<void> {
  const body: ApproveBody = comment ? { comment } : {}
  await apiClient.put(`/records/${type}/${id}/review-approve`, body)
}

// PUT /records/:type/:id/review-return  [supervisor only — comment REQUIRED]
export async function reviewReturn(type: RecordType, id: string, comment: string): Promise<void> {
  const body: ReturnBody = { comment }
  await apiClient.put(`/records/${type}/${id}/review-return`, body)
}

// PUT /records/:type/:id/final-approve  [manager only]
export async function finalApprove(type: RecordType, id: string, comment?: string): Promise<void> {
  const body: ApproveBody = comment ? { comment } : {}
  await apiClient.put(`/records/${type}/${id}/final-approve`, body)
}

// PUT /records/:type/:id/final-return  [manager only — comment REQUIRED]
export async function finalReturn(type: RecordType, id: string, comment: string): Promise<void> {
  const body: ReturnBody = { comment }
  await apiClient.put(`/records/${type}/${id}/final-return`, body)
}

// GET /records/:type/:id/comments  [all roles]
export async function getComments(type: RecordType, id: string): Promise<RecordComment[]> {
  const { data } = await apiClient.get<ApiResponse<{ comments: RecordComment[] }>>(
    `/records/${type}/${id}/comments`
  )
  return data.data.comments
}

// POST /records/:type/:id/comments  [supervisor + manager only]
// IMPORTANT: body field is "comment_text" NOT "comment"
export async function addComment(type: RecordType, id: string, commentText: string): Promise<RecordComment> {
  const body: AddCommentBody = { comment_text: commentText }
  const { data } = await apiClient.post<ApiResponse<{ comment: RecordComment }>>(
    `/records/${type}/${id}/comments`, body
  )
  return data.data.comment
}

// GET /records/:type/:id/history  [all roles]
export async function getHistory(type: RecordType, id: string): Promise<StatusHistoryEntry[]> {
  const { data } = await apiClient.get<ApiResponse<{ history: StatusHistoryEntry[] }>>(
    `/records/${type}/${id}/history`
  )
  return data.data.history
}
```

---

### 3. `frontend/src/components/forms/movableOptions.ts`

Create this options file following the same pattern as `immovableOptions.ts`:

```typescript
import type { Opt } from './immovableOptions'

export const MOVABLE_OWNER_TYPES: Opt[] = [
  { value: 'government', labelAm: 'መንግሥት', labelEn: 'Government' },
  { value: 'private', labelAm: 'የግል', labelEn: 'Private' },
  { value: 'religious', labelAm: 'ሃይማኖታዊ', labelEn: 'Religious' },
  { value: 'community', labelAm: 'ማኅበረሰብ', labelEn: 'Community' },
]

export const STORAGE_LOCATIONS: Opt[] = [
  { value: 'museum', labelAm: 'ሙዚየም', labelEn: 'Museum' },
  { value: 'church', labelAm: 'ቤተ ክርስቲያን', labelEn: 'Church' },
  { value: 'mosque', labelAm: 'መስጂድ', labelEn: 'Mosque' },
  { value: 'private_home', labelAm: 'የግል ቤት', labelEn: 'Private Home' },
  { value: 'office', labelAm: 'ቢሮ', labelEn: 'Office' },
  { value: 'other', labelAm: 'ሌላ', labelEn: 'Other' },
]

export const MOVABLE_AGE_METHODS: Opt[] = [
  { value: 'document', labelAm: 'ሰነድ', labelEn: 'Document' },
  { value: 'oral', labelAm: 'የቃል ማስረጃ', labelEn: 'Oral' },
  { value: 'estimate', labelAm: 'ግምት', labelEn: 'Estimate' },
]

export const MOVABLE_CONDITIONS: Opt[] = [
  { value: 'good', labelAm: 'ጥሩ', labelEn: 'Good' },
  { value: 'fair', labelAm: 'መካከለኛ', labelEn: 'Fair' },
  { value: 'poor', labelAm: 'ደካማ', labelEn: 'Poor' },
  { value: 'critical', labelAm: 'አስቸኳይ', labelEn: 'Critical' },
]

export const ACQUISITION_METHODS: Opt[] = [
  { value: 'purchase', labelAm: 'ግዢ', labelEn: 'Purchase' },
  { value: 'donation', labelAm: 'ስጦታ', labelEn: 'Donation' },
  { value: 'inheritance', labelAm: 'ውርስ', labelEn: 'Inheritance' },
  { value: 'transfer', labelAm: 'ዝውውር', labelEn: 'Transfer' },
  { value: 'found', labelAm: 'የተገኘ', labelEn: 'Found' },
]

// Reuse from immovableOptions
export { AGE_METHODS, QUALITY_LEVELS, ACCESSIBILITY_LEVELS, SEX_TYPES } from './immovableOptions'
```

---

### 4. `frontend/src/components/forms/MovableForm.tsx`

Create this component following the exact same pattern as `ImmovableForm.tsx`. It uses the same `Field`, `FormSection`, `TextInput`, etc. primitives.

**Sections:**

**Section 1 — Identification** (`am="መለያ"` `en="Identification"`)
Fields: `name_amharic` (required), `name_local`, `category` (TextInput — single string for movable, not array), `current_use`, `previous_id`

**Section 2 — Location & Storage** (`am="አካባቢ እና ማስቀመጫ"` `en="Location & Storage"`)
Fields: `location_name`, `woreda`, `kebele`, `house_number`, `owner_type` (Select, MOVABLE_OWNER_TYPES), `owner_name`, `storage_location` (Select, STORAGE_LOCATIONS), `storage_location_other`

**Section 3 — Origin & Acquisition** (`am="አመጣጥ እና ግዥ"` `en="Origin & Acquisition"`)
Fields: `made_by`, `period_made`, `age_method` (Select, MOVABLE_AGE_METHODS), `acquisition_methods` (CheckboxList, ACQUISITION_METHODS)

**Section 4 — Physical Description** (`am="አካላዊ ገጽታ"` `en="Physical Description"`)
Fields: `height_cm`, `width_cm`, `length_cm`, `diameter_cm`, `thickness_cm`, `weight_kg`, `num_pages`, `num_chapters`, `num_illustrations`, `color_type`, `has_decoration` (Switch), `materials` (CheckboxList — use values: `['wood','metal','textile','leather','stone','paper','ceramic','other']` with bilingual labels), `material_other`, `description` (TextArea span=3)

**Section 5 — Significance** (`am="ጠቀሜታ"` `en="Significance"`)
Fields: `notable_because` (CheckboxList — values: `['age','rarity','craftsmanship','historical','religious','cultural']` with bilingual labels), `notable_other`, `significance` (TextArea span=3)

**Section 6 — Condition & Conservation** (`am="ሁኔታ እና ጥበቃ"` `en="Condition & Conservation"`)
Fields: `condition` (Select, MOVABLE_CONDITIONS), `has_threat` (Switch), `threat_description`, `maintenance_done` (Switch), `maintenance_by`, `maintenance_date` (date input), `maintenance_count` (NumberInput), `preventive_level` (Select, QUALITY_LEVELS), `accessibility` (Select, ACCESSIBILITY_LEVELS), `notes` (TextArea span=3)

**Section 7 — Informant & Registrar** (`am="መረጃ ሰጪ እና መዝጋቢ"` `en="Informant & Registrar"`)
Fields: `informant_name`, `informant_sex` (Select, SEX_TYPES), `informant_age` (NumberInput), `informant_occupation`, `caretaker_name`, `caretaker_role`, `registrar_date` (date input)

**Form behavior — same as ImmovableForm:**
- Props: `mode: 'create' | 'edit'`, `initialRecord?: MovableRecord`, `photos?: RecordPhoto[]`
- State: all fields in a single `form` object, initialized from `initialRecord` if editing
- `saveMut`: on create calls `createMovable(form)`, navigates to edit route with new id; on edit calls `updateMovable(id, form)`
- `submitMut`: saves then calls `submitMovable(id)`, invalidates queries, navigates to `/records`
- Required field validation: only `name_amharic` is required before save
- Sticky bottom bar with Cancel / Save Draft / Submit buttons — Submit only shown in edit mode when status is not `approved`
- PhotoUploader shown in edit mode — use `uploadMovablePhoto` and `deleteMovablePhoto` from `src/api/movable.ts`

---

### 5. `frontend/src/components/forms/PhotoUploader.tsx` — make generic

Update `PhotoUploader` to accept upload/delete functions as props instead of calling the immovable-specific ones directly. This makes it reusable for both record types.

```typescript
interface PhotoUploaderProps {
  recordId: string
  photos: RecordPhoto[]
  disabled?: boolean
  onUpload: (file: File) => Promise<RecordPhoto>   // ← new
  onDelete: (photoId: string) => Promise<void>      // ← new
}
```

Update `ImmovableForm` to pass:
```typescript
onUpload={(file) => uploadImmovablePhoto(initialRecord.id, file)}
onDelete={(photoId) => deleteImmovablePhoto(initialRecord.id, photoId)}
```

---

### 6. `frontend/src/components/common/StatusTimeline.tsx`

Create this component:

```typescript
// Props:
interface StatusTimelineProps {
  history: StatusHistoryEntry[]
}
```

Renders a vertical timeline, most recent entry at top. Each entry:
- Colored status badge (use `StatusBadge` component)
- `changed_by_name` — the person who made the change
- Formatted `created_at` date (use `toLocaleDateString`)
- `note` — shown below if present, in smaller muted text
- Arrow from `from_status` → `to_status` (show `from_status` as plain muted text if present)

Visual design: vertical line connecting entries, with a colored dot per entry. Use Tailwind — no external libs.

If `history` is empty: show `EmptyState` with a clock icon and the i18n key `history.empty`.

---

### 7. `frontend/src/components/common/CommentThread.tsx`

Create this component:

```typescript
interface CommentThreadProps {
  comments: RecordComment[]
  recordType: RecordType
  recordId: string
  canComment: boolean   // true for supervisor and manager roles only
}
```

- Renders `RecordComment[]` as a list, newest at bottom
- Each comment: formatted date, `comment_text`, and `author_id` (display as "User" for now — we don't have a name lookup here)
- If `canComment` is true: show a textarea + "Add Comment" button at the bottom
  - On submit: calls `addComment(recordType, recordId, text)` from `src/api/workflow.ts`
  - Body field MUST be `{ comment_text: text }` — this is handled by the `addComment` function
  - Uses `useMutation` from TanStack Query; invalidates `['record-detail', recordType, recordId]` on success
  - Clears textarea on success
  - Shows error if mutation fails
- If `comments` is empty and `canComment` is false: show `EmptyState` with message key `comments.empty`
- If `comments` is empty and `canComment` is true: show a prompt to add the first comment

---

### 8. `frontend/src/components/common/ReturnModal.tsx`

Create this component (used in Prompt 5 by supervisor, built now so it's ready):

```typescript
interface ReturnModalProps {
  open: boolean
  onClose: () => void
  onConfirm: (comment: string) => Promise<void>
  title: string      // e.g. "Return to Registrar"
  titleAm: string    // e.g. "ወደ መዝጋቢ መለስ"
}
```

- Uses shadcn/ui `Dialog` from `src/components/ui/dialog.tsx`
- Contains a `<textarea>` for the required return comment (min 1 character — show inline error if empty on submit)
- "Cancel" button closes the modal
- "Confirm Return" button: disabled while loading; calls `onConfirm(comment)` and closes on success
- Shows spinner on the confirm button while pending
- The comment field label: `am="የመልስ ምክንያት"` `en="Reason for return (required)"`

---

### 9. Update `/src/routes/_authenticated/records.new.$type.tsx`

Replace the "coming soon" movable placeholder with the real form:

```typescript
import { MovableForm } from '@/components/forms/MovableForm'

// In the render:
{type === 'movable' ? <MovableForm mode="create" /> : <ImmovableForm mode="create" />}
```

---

### 10. Update `/src/routes/_authenticated/records.$id.edit.tsx`

This route currently only handles immovable. Update it to:

1. First try to load from immovable: `GET /records/immovable/:id`
2. If that returns 404, try movable: `GET /records/movable/:id`
3. Render the correct form + detail based on which succeeded

Better approach: the `RecordSummary` from the list contains `record_type`. When the user clicks a record card, pass the type. But since we can't change the URL structure right now, do this:

**Use parallel queries — try immovable first, fall back to movable:**

```typescript
// Try immovable
const immovableQ = useQuery({
  queryKey: ['immovable', id],
  queryFn: () => getImmovable(id),
  retry: false,
})

// Try movable only if immovable returned 404 or error
const movableQ = useQuery({
  queryKey: ['movable', id],
  queryFn: () => getMovable(id),
  enabled: immovableQ.isError,
  retry: false,
})
```

Then render based on which query has data.

**The detail page renders:**

```
┌─────────────────────────────────────────────┐
│  [← My Records]  Record ID badge  Status badge │
├─────────────────────────────────────────────┤
│  [Form — left column, 2/3 width]            │
│                                             │
├─────────────────────────────────────────────┤
│  TABS: [History] [Comments] [Photos]        │
│  ← StatusTimeline / CommentThread / grid    │
└─────────────────────────────────────────────┘
```

- Form section: renders `<ImmovableForm mode="edit">` or `<MovableForm mode="edit">` with `initialRecord` and `photos`
- Below the form, render a tabbed panel using shadcn/ui `Tabs` from `src/components/ui/tabs.tsx`:
  - **History tab**: `<StatusTimeline history={detail.history} />`
  - **Comments tab**: `<CommentThread comments={detail.comments} recordType={...} recordId={id} canComment={role === 'supervisor' || role === 'manager'} />`
  - **Photos tab**: shows the photo grid (read-only — editing is handled inside the form already)
- `canComment` comes from `useAuthStore((s) => s.user?.role)`

---

### 11. Update `src/routes/_authenticated/records.index.tsx`

The `RecordCard` currently links to `/records/$id/edit`. This is correct for registrar. Keep it.

---

### 12. i18n additions

Add these keys to **both** `src/i18n/en.json` and `src/i18n/am.json`:

**en.json additions:**
```json
{
  "movable": {
    "title": "Movable Heritage Asset",
    "subtitle": "Fill in all required sections. You can save a draft and continue later."
  },
  "history": {
    "title": "Status History",
    "empty": "No status changes recorded yet."
  },
  "comments": {
    "title": "Comments",
    "empty": "No comments yet.",
    "add": "Add comment",
    "placeholder": "Write your comment...",
    "submit": "Add Comment",
    "error": "Failed to add comment."
  },
  "detail": {
    "tabs": {
      "history": "History",
      "comments": "Comments",
      "photos": "Photos"
    },
    "backToRecords": "My Records"
  },
  "modal": {
    "returnTitle": "Return Record",
    "returnTitleAm": "መዝገቡን መለስ",
    "reasonLabel": "Reason for return (required)",
    "reasonLabelAm": "የመልስ ምክንያት",
    "cancel": "Cancel",
    "confirm": "Confirm Return"
  }
}
```

**am.json additions:**
```json
{
  "movable": {
    "title": "ተንቀሳቃሽ ቅርስ",
    "subtitle": "ሁሉንም አስፈላጊ ክፍሎች ይሙሉ። ረቂቅ አስቀምጠው በኋላ መቀጠል ይችላሉ።"
  },
  "history": {
    "title": "የሁኔታ ታሪክ",
    "empty": "እስካሁን ምንም ለውጥ አልተመዘገበም።"
  },
  "comments": {
    "title": "አስተያየቶች",
    "empty": "እስካሁን አስተያየት የለም።",
    "add": "አስተያየት አክል",
    "placeholder": "አስተያየትዎን ይጻፉ...",
    "submit": "አስተያየት አክል",
    "error": "አስተያየት ማከል አልተቻለም።"
  },
  "detail": {
    "tabs": {
      "history": "ታሪክ",
      "comments": "አስተያየቶች",
      "photos": "ፎቶዎች"
    },
    "backToRecords": "የእኔ መዛግብት"
  },
  "modal": {
    "returnTitle": "Return Record",
    "returnTitleAm": "መዝገቡን መለስ",
    "reasonLabel": "Reason for return (required)",
    "reasonLabelAm": "የመልስ ምክንያት (አስፈላጊ)",
    "cancel": "ሰርዝ",
    "confirm": "አረጋግጥ"
  }
}
```

---

## IMPORTANT RULES — DO NOT VIOLATE

1. **All API calls go through `src/api/client.ts`** — never use `fetch()` or raw `axios` directly in components.

2. **URL correctness** — The existing `src/api/immovable.ts` has a bug: it calls `/immovable` and `/immovable/:id` instead of `/records/immovable` and `/records/immovable/:id`. Fix these while you're here:
   - `listImmovable`: change `/immovable` → `/records/immovable`
   - `getImmovable`: change `/immovable/${id}` → `/records/immovable/${id}`
   - `updateImmovable`: change `/immovable/${id}` → `/records/immovable/${id}`
   - `submitImmovable`: change `/immovable/${id}/submit` → `/records/immovable/${id}/submit` (and change from POST to PUT)
   - `uploadImmovablePhoto`: change `/immovable/${id}/photos` → `/records/immovable/${id}/photos`
   - `deleteImmovablePhoto`: change `/immovable/${recordId}/photos/${photoId}` → `/records/immovable/${recordId}/photos/${photoId}`
   - Also fix `src/api/records.ts`: `listMyRecords` calls `/records/mine` which doesn't exist — change to `/records` (the server scopes it automatically by JWT role)

3. **No hardcoded strings** — every visible string uses `t('key')`.

4. **Every data-fetching component** must handle loading, error, and empty states.

5. **Comment body** — the `addComment` function sends `{ comment_text: "..." }`. Do not change this field name.

6. **Photo upload** — multipart form with field name `"photo"`. Already correct in existing code.

7. **Do not modify** `src/types/index.ts`, `src/stores/authStore.ts`, `src/stores/languageStore.ts`, `src/routes/__root.tsx`, or `src/routes/_authenticated.tsx`.

8. **TanStack Router** — if you add any new route files, they must follow the file-naming convention already in `src/routes/`. Run `bun run build` or let the router plugin auto-generate `routeTree.gen.ts` — do not manually edit `routeTree.gen.ts`.
