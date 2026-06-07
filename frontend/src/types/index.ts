// ============================================================
// Qirs Mezgeb types — generated from Go backend structs
// ============================================================

export type Role = "registrar" | "supervisor" | "manager";
export type Language = "am" | "en";
export type RecordStatus =
  | "draft"
  | "pending_review"
  | "under_review"
  | "returned"
  | "approved";
export type RecordType = "immovable" | "movable";

export type ImmovableOwnerType =
  | "public"
  | "government"
  | "religion"
  | "private"
  | "waqf"
  | string;
export type AgeMethod = "estimated" | "exact" | "relative" | string;
export type OverallCondition =
  | "very_good"
  | "good"
  | "damaged"
  | "severely_damaged"
  | string;
export type DamageLevel = "minor" | "moderate" | "medium" | "severe" | string;
export type QualityLevel =
  | "very_good"
  | "good"
  | "medium"
  | "low"
  | "very_low"
  | string;
export type AccessibilityLevel =
  | "very_good"
  | "good"
  | "medium"
  | "low"
  | "very_low"
  | "none"
  | string;
export type SexType = "male" | "female" | string;
export type MovableOwnerType = "public" | "government" | "religion" | "private" | string;
export type StorageLocation =
  | "museum"
  | "store"
  | "church"
  | "private_home"
  | "other"
  | string;
export type MovableCondition = "good" | "fair" | "damaged" | "incomplete" | string;

// ── API envelope ─────────────────────────────────────────────
export interface ApiResponse<T> {
  success: true;
  data: T;
  message: string;
}
export interface ApiError {
  success: false;
  error: string;
  code: number;
}
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

// ── Auth / Users ─────────────────────────────────────────────
export interface UserPublic {
  id: string;
  full_name: string;
  email: string;
  role: Role;
  language: Language;
}

export interface UserItem {
  id: string;
  full_name: string; // /users list uses "full_name"
  email: string;
  role: Role;
  language: Language;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export type PaginatedUsers = PaginatedResponse<UserItem>;

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  user: UserPublic;
}

export interface CreateUserBody {
  full_name: string;
  email: string;
  password: string;
  role: Role;
  language?: Language;
}
export interface UpdateUserBody {
  full_name?: string;
  role?: Role;
  language?: Language;
  is_active?: boolean;
}
export interface UpdateLanguageBody {
  language: Language;
}

// ── Records (shared sub-resources) ───────────────────────────
export interface RecordPhoto {
  id: string;
  record_type: RecordType;
  record_id: string;
  file_path: string;
  file_name?: string;
  file_size_bytes?: number;
  uploaded_by: string;
  created_at: string;
}

export interface RecordComment {
  id: string;
  record_type: RecordType;
  record_id: string;
  author_id: string;
  comment_text: string;
  created_at: string;
}

export interface StatusHistoryEntry {
  id: string;
  record_type: RecordType;
  record_id: string;
  changed_by: string;
  changed_by_name: string;
  from_status?: string;
  to_status: string;
  note?: string;
  created_at: string;
}

// ── Immovable ────────────────────────────────────────────────
export interface ImmovableRecord {
  id: string;
  record_id: string;
  registrar_id: string;
  status: RecordStatus;
  name_amharic: string;
  name_local?: string;
  category?: string[];
  current_use?: string[];
  current_use_other?: string;
  previous_id?: string;
  woreda: string;
  kebele: string;
  house_number?: string;
  street_number?: string;
  gate?: string;
  owner_type?: ImmovableOwnerType;
  owner_name?: string;
  map_reference?: string;
  gps_east?: number;
  gps_north?: number;
  elevation_m?: number;
  built_by?: string;
  construction_period?: string;
  age_method?: AgeMethod;
  height_m?: number;
  length_m?: number;
  width_m?: number;
  num_doors?: number;
  num_windows?: number;
  num_rooms?: number;
  material?: string;
  description?: string;
  harari_house_grades?: string[];
  neighborhood_type?: string;
  overall_condition?: OverallCondition;
  damage_roof?: DamageLevel;
  damage_cornice?: DamageLevel;
  damage_wall?: DamageLevel;
  damage_floor?: DamageLevel;
  damage_door?: DamageLevel;
  damage_cupboard?: DamageLevel;
  damage_upper_floor?: DamageLevel;
  damage_dera?: DamageLevel;
  damage_pillar?: DamageLevel;
  value_historical?: string;
  value_craftsmanship?: string;
  value_artistic?: string;
  value_scientific?: string;
  value_cultural?: string;
  has_threat?: boolean;
  maintenance_done?: boolean;
  maintenance_reason?: string;
  maintenance_by?: string;
  maintenance_date?: string;
  maintenance_count?: number;
  preventive_level?: QualityLevel;
  accessibility?: AccessibilityLevel;
  notes?: string;
  related_docs?: string[];
  has_oral_history?: boolean;
  caretaker_name?: string;
  caretaker_role?: string;
  informant_name?: string;
  informant_sex?: SexType;
  informant_age?: number;
  registrar_date?: string;
  approved_at?: string;
  approved_by?: string;
  created_at: string;
  updated_at: string;
}

export type ImmovableRecordInput = Partial<
  Omit<
    ImmovableRecord,
    | "id"
    | "record_id"
    | "registrar_id"
    | "status"
    | "approved_at"
    | "approved_by"
    | "created_at"
    | "updated_at"
  >
>;

export interface ImmovableRecordDetail {
  record: ImmovableRecord;
  photos: RecordPhoto[];
  comments: RecordComment[];
  history: StatusHistoryEntry[];
}

export interface ImmovableCreateResult {
  id: string;
  record_id: string;
  status: RecordStatus;
}

export interface SubmitResult {
  status: RecordStatus;
}

// ── Movable ──────────────────────────────────────────────────
export interface MovableRecord {
  id: string;
  record_id: string;
  registrar_id: string;
  status: RecordStatus;
  name_amharic: string;
  name_local?: string;
  category?: string;
  location_name?: string;
  woreda?: string;
  kebele?: string;
  house_number?: string;
  current_use?: string;
  previous_id?: string;
  owner_type?: MovableOwnerType;
  owner_name?: string;
  storage_location?: StorageLocation;
  storage_location_other?: string;
  made_by?: string;
  period_made?: string;
  age_method?: AgeMethod;
  acquisition_methods?: string[];
  height_cm?: number;
  width_cm?: number;
  length_cm?: number;
  diameter_cm?: number;
  thickness_cm?: number;
  weight_kg?: number;
  num_pages?: number;
  num_chapters?: number;
  num_illustrations?: number;
  color_type?: string;
  has_decoration?: boolean;
  materials?: string[];
  material_other?: string;
  description?: string;
  notable_because?: string[];
  notable_other?: string;
  significance?: string;
  condition?: MovableCondition;
  has_threat?: boolean;
  threat_description?: string;
  maintenance_done?: boolean;
  maintenance_by?: string;
  maintenance_date?: string;
  maintenance_count?: number;
  preventive_level?: QualityLevel;
  accessibility?: AccessibilityLevel;
  notes?: string;
  related_docs?: string[];
  informant_name?: string;
  informant_sex?: SexType;
  informant_age?: number;
  informant_occupation?: string;
  caretaker_name?: string;
  caretaker_role?: string;
  registrar_date?: string;
  approved_at?: string;
  approved_by?: string;
  created_at: string;
  updated_at: string;
}

export type MovableRecordInput = Partial<
  Omit<
    MovableRecord,
    | "id"
    | "record_id"
    | "registrar_id"
    | "status"
    | "approved_at"
    | "approved_by"
    | "created_at"
    | "updated_at"
  >
>;

export interface MovableRecordDetail {
  record: MovableRecord;
  photos: RecordPhoto[];
  comments: RecordComment[];
  history: StatusHistoryEntry[];
}

export interface MovableCreateResult {
  id: string;
  record_id: string;
  status: RecordStatus;
}

// ── Workflow bodies ──────────────────────────────────────────
export interface ApproveBody {
  comment?: string;
}
export interface ReturnBody {
  comment: string;
}
export interface AddCommentBody {
  comment_text: string;
}

// ── Dashboard ────────────────────────────────────────────────
export interface StatusCounts {
  draft: number;
  pending_review: number;
  under_review: number;
  returned: number;
  approved: number;
}
export interface DashboardStats {
  total_immovable: number;
  total_movable: number;
  by_status: StatusCounts;
  pending_my_action?: number;
}
export interface RecordSummary {
  id: string;
  record_type: RecordType;
  record_id: string;
  name_amharic: string;
  status: RecordStatus;
  woreda?: string;
  kebele?: string;
  registrar_id: string;
  created_at: string;
  updated_at: string;
}

export type PaginatedRecordSummaries = PaginatedResponse<RecordSummary>;
export type PaginatedImmovable = PaginatedResponse<ImmovableRecord>;
export type PaginatedMovable = PaginatedResponse<MovableRecord>;

// ── Query params ─────────────────────────────────────────────
export interface RecordListParams {
  page?: number;
  limit?: number;
  status?: RecordStatus;
  woreda?: string;
  search?: string;
  date_from?: string;
  date_to?: string;
}
export interface UnifiedListParams extends RecordListParams {
  type?: RecordType;
  kebele?: string;
}
export interface UserListParams {
  page?: number;
  limit?: number;
  role?: Role;
  is_active?: boolean;
}

// ── Auth store ───────────────────────────────────────────────
export interface AuthState {
  user: UserPublic | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
}
