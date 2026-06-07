const DATE_FIELDS = new Set(["maintenance_date", "registrar_date"]);

export function toDateInputValue(iso?: string | null): string {
  if (!iso) return "";
  const date = iso.slice(0, 10);
  return /^\d{4}-\d{2}-\d{2}$/.test(date) ? date : "";
}

/** Strip empty optional values and normalize date fields for Go JSON binding. */
export function sanitizeRecordPayload(
  body: Record<string, unknown>,
): Record<string, unknown> {
  const out: Record<string, unknown> = {};

  for (const [key, value] of Object.entries(body)) {
    if (value === "" || value === null || value === undefined) continue;

    if (DATE_FIELDS.has(key) && typeof value === "string") {
      out[key] = `${value.slice(0, 10)}T00:00:00Z`;
      continue;
    }

    out[key] = value;
  }

  return out;
}
