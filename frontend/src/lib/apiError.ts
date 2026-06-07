import axios from "axios";

export function getApiErrorMessage(err: unknown, fallback = "Request failed"): string {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as
      | { error?: string; fields?: Record<string, string> }
      | undefined;

    if (data?.fields && Object.keys(data.fields).length > 0) {
      return `${data.error ?? "Validation failed"} (${Object.keys(data.fields).join(", ")})`;
    }

    return data?.error ?? err.message;
  }

  if (err instanceof Error) return err.message;
  return fallback;
}
