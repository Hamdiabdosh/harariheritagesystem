import { apiClient } from "./client";
import type {
  ApiResponse,
  RecordType,
  RecordComment,
  StatusHistoryEntry,
  AddCommentBody,
  ApproveBody,
  ReturnBody,
} from "@/types";

export async function reviewApprove(
  type: RecordType,
  id: string,
  comment?: string,
): Promise<void> {
  const body: ApproveBody = comment ? { comment } : {};
  await apiClient.put(`/records/${type}/${id}/review-approve`, body);
}

export async function reviewReturn(
  type: RecordType,
  id: string,
  comment: string,
): Promise<void> {
  const body: ReturnBody = { comment };
  await apiClient.put(`/records/${type}/${id}/review-return`, body);
}

export async function finalApprove(
  type: RecordType,
  id: string,
  comment?: string,
): Promise<void> {
  const body: ApproveBody = comment ? { comment } : {};
  await apiClient.put(`/records/${type}/${id}/final-approve`, body);
}

export async function finalReturn(
  type: RecordType,
  id: string,
  comment: string,
): Promise<void> {
  const body: ReturnBody = { comment };
  await apiClient.put(`/records/${type}/${id}/final-return`, body);
}

export async function getComments(
  type: RecordType,
  id: string,
): Promise<RecordComment[]> {
  const { data } = await apiClient.get<ApiResponse<{ comments: RecordComment[] }>>(
    `/records/${type}/${id}/comments`,
  );
  return data.data.comments;
}

export async function addComment(
  type: RecordType,
  id: string,
  commentText: string,
): Promise<RecordComment> {
  const body: AddCommentBody = { comment_text: commentText };
  const { data } = await apiClient.post<ApiResponse<{ comment: RecordComment }>>(
    `/records/${type}/${id}/comments`,
    body,
  );
  return data.data.comment;
}

export async function getHistory(
  type: RecordType,
  id: string,
): Promise<StatusHistoryEntry[]> {
  const { data } = await apiClient.get<ApiResponse<{ history: StatusHistoryEntry[] }>>(
    `/records/${type}/${id}/history`,
  );
  return data.data.history;
}
