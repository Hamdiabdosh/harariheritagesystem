import { apiClient } from "./client";
import type { ApiResponse, TokenPair } from "@/types";

export async function login(email: string, password: string): Promise<TokenPair> {
  const { data } = await apiClient.post<ApiResponse<TokenPair>>("/auth/login", {
    email,
    password,
  });
  return data.data;
}

export async function refreshToken(refresh_token: string): Promise<string> {
  const { data } = await apiClient.post<ApiResponse<{ access_token: string }>>(
    "/auth/refresh",
    { refresh_token },
  );
  return data.data.access_token;
}

export async function logout(refresh_token: string): Promise<void> {
  await apiClient.post("/auth/logout", { refresh_token });
}
