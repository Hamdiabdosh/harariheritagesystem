import { apiClient } from "./client";
import type { ApiResponse, RecordComment, RecordStatus, RecordType } from "@/types";

interface StatusResult {
  status: RecordStatus;
  approved_at?: string;
}

export async function reviewApprove(
  recordType: RecordType,
  recordId: string,
  comment?: string,
): Promise<StatusResult> {
  const { data } = await apiClient.put<ApiResponse<StatusResult>>(
    `/records/${recordType}/${recordId}/review-approve`,
    comment ? { comment_text: comment } : {},
  );
  return data.data;
}

export async function reviewReturn(
  recordType: RecordType,
  recordId: string,
  comment: string,
): Promise<StatusResult> {
  const { data } = await apiClient.put<ApiResponse<StatusResult>>(
    `/records/${recordType}/${recordId}/review-return`,
    { comment_text: comment },
  );
  return data.data;
}

export async function finalApprove(
  recordType: RecordType,
  recordId: string,
  comment?: string,
): Promise<StatusResult> {
  const { data } = await apiClient.put<ApiResponse<StatusResult>>(
    `/records/${recordType}/${recordId}/final-approve`,
    comment ? { comment_text: comment } : {},
  );
  return data.data;
}

export async function finalReturn(
  recordType: RecordType,
  recordId: string,
  comment: string,
): Promise<StatusResult> {
  const { data } = await apiClient.put<ApiResponse<StatusResult>>(
    `/records/${recordType}/${recordId}/final-return`,
    { comment_text: comment },
  );
  return data.data;
}

export async function addComment(
  recordType: RecordType,
  recordId: string,
  commentText: string,
): Promise<RecordComment> {
  const { data } = await apiClient.post<ApiResponse<{ comment: RecordComment }>>(
    `/records/${recordType}/${recordId}/comments`,
    { comment_text: commentText },
  );
  return data.data.comment;
}
