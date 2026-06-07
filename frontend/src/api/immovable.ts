import { apiClient } from "./client";
import { sanitizeRecordPayload } from "@/lib/recordPayload";
import type {
  ApiResponse,
  ImmovableCreateResult,
  ImmovableRecordDetail,
  ImmovableRecordInput,
  PaginatedImmovable,
  RecordListParams,
  RecordPhoto,
  SubmitResult,
} from "@/types";

export async function listImmovable(
  params: RecordListParams = {},
): Promise<PaginatedImmovable> {
  const { data } = await apiClient.get<ApiResponse<PaginatedImmovable>>(
    "/records/immovable",
    { params },
  );
  return data.data;
}

export async function getImmovable(id: string): Promise<ImmovableRecordDetail> {
  const { data } = await apiClient.get<ApiResponse<ImmovableRecordDetail>>(
    `/records/immovable/${id}`,
  );
  return data.data;
}

export async function createImmovable(
  body: ImmovableRecordInput,
): Promise<ImmovableCreateResult> {
  const { data } = await apiClient.post<ApiResponse<ImmovableCreateResult>>(
    "/records/immovable",
    sanitizeRecordPayload(body as Record<string, unknown>),
  );
  return data.data;
}

export async function updateImmovable(
  id: string,
  body: ImmovableRecordInput,
): Promise<void> {
  await apiClient.put(
    `/records/immovable/${id}`,
    sanitizeRecordPayload(body as Record<string, unknown>),
  );
}

export async function submitImmovable(id: string): Promise<SubmitResult> {
  const { data } = await apiClient.put<ApiResponse<SubmitResult>>(
    `/records/immovable/${id}/submit`,
  );
  return data.data;
}

export async function uploadImmovablePhoto(
  id: string,
  file: File,
): Promise<RecordPhoto> {
  const fd = new FormData();
  fd.append("photo", file);
  const { data } = await apiClient.post<ApiResponse<RecordPhoto>>(
    `/records/immovable/${id}/photos`,
    fd,
    { headers: { "Content-Type": "multipart/form-data" } },
  );
  return data.data;
}

export async function deleteImmovablePhoto(
  recordId: string,
  photoId: string,
): Promise<void> {
  await apiClient.delete(`/records/immovable/${recordId}/photos/${photoId}`);
}
