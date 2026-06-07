import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Upload, X, ImageIcon, Loader2 } from "lucide-react";
import { useMutation } from "@tanstack/react-query";
import type { RecordPhoto } from "@/types";
import { getPhotoUrl } from "@/lib/photos";

const MAX_SIDE = 1200;
const JPEG_QUALITY = 0.82;

async function compressImage(file: File): Promise<File> {
  if (file.type !== "image/jpeg" && file.type !== "image/png") {
    return file;
  }

  try {
    const bitmap = await createImageBitmap(file);
    const scale = Math.min(1, MAX_SIDE / Math.max(bitmap.width, bitmap.height));
    const width = Math.round(bitmap.width * scale);
    const height = Math.round(bitmap.height * scale);

    const canvas = document.createElement("canvas");
    canvas.width = width;
    canvas.height = height;

    const ctx = canvas.getContext("2d");
    if (!ctx) {
      bitmap.close();
      return file;
    }

    ctx.drawImage(bitmap, 0, 0, width, height);
    bitmap.close();

    const blob = await new Promise<Blob | null>((resolve) => {
      canvas.toBlob(resolve, "image/jpeg", JPEG_QUALITY);
    });
    if (!blob) return file;

    return new File([blob], file.name, {
      type: "image/jpeg",
      lastModified: Date.now(),
    });
  } catch {
    return file;
  }
}

interface PhotoUploaderProps {
  recordId: string;
  photos: RecordPhoto[];
  disabled?: boolean;
  onUpload: (file: File) => Promise<RecordPhoto>;
  onDelete: (photoId: string) => Promise<void>;
  onPhotosChange?: () => void;
}

export function PhotoUploader({
  recordId: _recordId,
  photos,
  disabled,
  onUpload,
  onDelete,
  onPhotosChange,
}: PhotoUploaderProps) {
  const { t } = useTranslation();
  const inputRef = useRef<HTMLInputElement>(null);
  const [error, setError] = useState<string | null>(null);
  const [compressing, setCompressing] = useState(false);

  const uploadMut = useMutation({
    mutationFn: onUpload,
    onSuccess: () => onPhotosChange?.(),
    onError: (e: Error) => setError(e.message),
  });

  const deleteMut = useMutation({
    mutationFn: onDelete,
    onSuccess: () => onPhotosChange?.(),
  });

  const handleFiles = async (files: FileList | null) => {
    if (!files) return;
    setError(null);
    for (const file of Array.from(files)) {
      if (file.size > 10 * 1024 * 1024) {
        setError(t("photos.tooLarge", { name: file.name }));
        continue;
      }
      setCompressing(true);
      let uploadFile = file;
      try {
        uploadFile = await compressImage(file);
      } finally {
        setCompressing(false);
      }
      await uploadMut.mutateAsync(uploadFile).catch(() => {});
    }
    if (inputRef.current) inputRef.current.value = "";
  };

  return (
    <div className="space-y-3 rounded-xl border border-border bg-card p-4">
      <div className="flex items-center justify-between gap-2">
        <div>
          <h3 className="font-amharic text-sm font-semibold text-foreground">
            {t("photos.title")}
          </h3>
          <p className="text-xs text-muted-foreground">
            {t("photos.count", { n: photos.length })}
          </p>
        </div>
        <button
          type="button"
          disabled={disabled || compressing || uploadMut.isPending}
          onClick={() => inputRef.current?.click()}
          className="inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
        >
          {compressing || uploadMut.isPending ? (
            <Loader2 className="h-3.5 w-3.5 animate-spin" />
          ) : (
            <Upload className="h-3.5 w-3.5" />
          )}
          <span className="font-amharic">
            {compressing
              ? t("photos.compressing", { defaultValue: "Compressing..." })
              : t("photos.upload")}
          </span>
        </button>
        <input
          ref={inputRef}
          type="file"
          accept="image/*"
          multiple
          className="hidden"
          onChange={(e) => handleFiles(e.target.files)}
        />
      </div>

      {error && (
        <div className="rounded-md bg-rose-50 px-3 py-2 text-xs text-rose-900">
          {error}
        </div>
      )}

      {photos.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-md border border-dashed border-border bg-muted/30 py-8 text-center">
          <ImageIcon className="mb-2 h-6 w-6 text-muted-foreground" />
          <p className="font-amharic text-xs text-muted-foreground">{t("photos.empty")}</p>
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-3 md:grid-cols-4">
          {photos.map((p) => {
            const src = getPhotoUrl(p.file_path);
            return (
              <div
                key={p.id}
                className="group relative aspect-square overflow-hidden rounded-md border border-border bg-muted"
              >
                <img
                  src={src}
                  alt={p.file_name ?? "photo"}
                  className="h-full w-full object-cover"
                  loading="lazy"
                />
                {!disabled && (
                  <button
                    type="button"
                    onClick={() => deleteMut.mutate(p.id)}
                    className="absolute right-1 top-1 flex h-6 w-6 items-center justify-center rounded-full bg-rose-600 text-white opacity-0 transition group-hover:opacity-100"
                    aria-label="Delete photo"
                  >
                    <X className="h-3.5 w-3.5" />
                  </button>
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
