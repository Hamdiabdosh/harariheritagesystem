import { apiClient } from "./client";
import type {
  ApiResponse,
  UserItem,
  PaginatedUsers,
  UserListParams,
  Language,
  CreateUserBody,
  UpdateUserBody,
} from "@/types";

export async function getMe(): Promise<UserItem> {
  const { data } = await apiClient.get<ApiResponse<{ user: UserItem }>>("/users/me");
  return data.data.user;
}

export async function updateMyLanguage(language: Language): Promise<void> {
  await apiClient.put("/users/me/language", { language });
}

export async function listUsers(params: UserListParams = {}): Promise<PaginatedUsers> {
  const { data } = await apiClient.get<ApiResponse<PaginatedUsers>>("/users", { params });
  return data.data;
}

export async function createUser(body: CreateUserBody): Promise<UserItem> {
  const { data } = await apiClient.post<ApiResponse<{ user: UserItem }>>("/users", body);
  return data.data.user;
}

export async function updateUser(id: string, body: UpdateUserBody): Promise<UserItem> {
  const { data } = await apiClient.put<ApiResponse<{ user: UserItem }>>(`/users/${id}`, body);
  return data.data.user;
}

export async function deactivateUser(id: string): Promise<void> {
  await apiClient.delete(`/users/${id}`);
}
