import { apiClient } from "./client";
import type {
  ApiResponse,
  PaginatedRecordSummaries,
  UnifiedListParams,
} from "@/types";

/**
 * Unified records list — GET /records
 * Returns both immovable & movable in a single paginated envelope.
 */
export async function listRecords(
  params: UnifiedListParams = {},
): Promise<PaginatedRecordSummaries> {
  const { data } = await apiClient.get<ApiResponse<PaginatedRecordSummaries>>(
    "/records",
    { params },
  );
  return data.data;
}

/**
 * "My records" — same endpoint, scoped by the registrar's own ID server-side
 * when the JWT role is `registrar`. The API does the scoping; the client just
 * forwards filters.
 */
export async function listMyRecords(
  params: UnifiedListParams = {},
): Promise<PaginatedRecordSummaries> {
  return listRecords(params);
}
