import { apiClient } from "./client";
import type { UnifiedListParams } from "@/types";
import { getApiErrorMessage } from "@/lib/apiError";

async function readBlobError(blob: Blob): Promise<string> {
  try {
    const text = await blob.text();
    const data = JSON.parse(text) as { error?: string };
    return data.error ?? "Export failed";
  } catch {
    return "Export failed";
  }
}

async function fetchBlob(path: string, params?: UnifiedListParams): Promise<Blob> {
  try {
    const response = await apiClient.get(path, {
      params,
      responseType: "blob",
    });
    return response.data as Blob;
  } catch (err: unknown) {
    if (
      typeof err === "object" &&
      err !== null &&
      "response" in err &&
      typeof (err as { response?: { data?: Blob } }).response?.data === "object"
    ) {
      const blob = (err as { response: { data: Blob } }).response.data;
      throw new Error(await readBlobError(blob));
    }
    throw new Error(getApiErrorMessage(err, "Export failed"));
  }
}

export async function exportCSV(params: UnifiedListParams = {}): Promise<Blob> {
  return fetchBlob("/export/records/csv", params);
}

export async function exportPDF(
  recordType: "immovable" | "movable",
  recordId: string,
): Promise<Blob> {
  return fetchBlob(`/records/${recordType}/${recordId}/pdf`);
}

export function downloadBlob(blob: Blob, filename: string): void {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}
