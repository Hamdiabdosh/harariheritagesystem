import { apiClient } from "./client";
import type { ApiResponse, DashboardStats } from "@/types";

export async function getDashboardStats(): Promise<DashboardStats> {
  const { data } = await apiClient.get<ApiResponse<DashboardStats>>("/dashboard/stats");
  return data.data;
}
