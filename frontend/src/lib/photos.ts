/** Build a browser URL for a stored photo file_path. */
export function getPhotoUrl(filePath: string): string {
  if (filePath.startsWith("http")) return filePath;
  const apiUrl = import.meta.env.VITE_API_URL ?? "http://localhost:8080/api/v1";
  const origin = apiUrl.replace(/\/api\/v1\/?$/, "");
  return `${origin}/media/${filePath.replace(/^\//, "")}`;
}
