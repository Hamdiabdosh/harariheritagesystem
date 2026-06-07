import { apiClient } from "./client";
import { sanitizeRecordPayload } from "@/lib/recordPayload";
import type {
  ApiResponse,
  MovableCreateResult,
  MovableRecordDetail,
  MovableRecordInput,
  PaginatedMovable,
  RecordListParams,
  RecordPhoto,
  SubmitResult,
} from "@/types";

export async function listMovable(
  params: RecordListParams = {},
): Promise<PaginatedMovable> {
  const { data } = await apiClient.get<ApiResponse<PaginatedMovable>>(
    "/records/movable",
    { params },
  );
  return data.data;
}

export async function getMovable(id: string): Promise<MovableRecordDetail> {
  const { data } = await apiClient.get<ApiResponse<MovableRecordDetail>>(
    `/records/movable/${id}`,
  );
  return data.data;
}

export async function createMovable(
  body: MovableRecordInput,
): Promise<MovableCreateResult> {
  const { data } = await apiClient.post<ApiResponse<MovableCreateResult>>(
    "/records/movable",
    sanitizeRecordPayload(body as Record<string, unknown>),
  );
  return data.data;
}

export async function updateMovable(
  id: string,
  body: MovableRecordInput,
): Promise<void> {
  await apiClient.put(
    `/records/movable/${id}`,
    sanitizeRecordPayload(body as Record<string, unknown>),
  );
}

export async function submitMovable(id: string): Promise<SubmitResult> {
  const { data } = await apiClient.put<ApiResponse<SubmitResult>>(
    `/records/movable/${id}/submit`,
  );
  return data.data;
}

export async function uploadMovablePhoto(
  id: string,
  file: File,
): Promise<RecordPhoto> {
  const fd = new FormData();
  fd.append("photo", file);
  const { data } = await apiClient.post<ApiResponse<RecordPhoto>>(
    `/records/movable/${id}/photos`,
    fd,
    { headers: { "Content-Type": "multipart/form-data" } },
  );
  return data.data;
}

export async function deleteMovablePhoto(
  recordId: string,
  photoId: string,
): Promise<void> {
  await apiClient.delete(`/records/movable/${recordId}/photos/${photoId}`);
}
