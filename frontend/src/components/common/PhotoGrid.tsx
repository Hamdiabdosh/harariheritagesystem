import { ImageIcon } from "lucide-react";
import { useTranslation } from "react-i18next";
import type { RecordPhoto } from "@/types";
import { getPhotoUrl } from "@/lib/photos";
import { EmptyState } from "@/components/common/EmptyState";

export function PhotoGrid({ photos }: { photos: RecordPhoto[] }) {
  const { t } = useTranslation();

  if (photos.length === 0) {
    return (
      <EmptyState
        icon={<ImageIcon className="h-6 w-6" />}
        title={t("photos.empty")}
      />
    );
  }

  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4">
      {photos.map((p) => {
        const src = getPhotoUrl(p.file_path);
        return (
          <div
            key={p.id}
            className="overflow-hidden rounded-lg border border-border bg-muted aspect-square"
          >
            <img
              src={src}
              alt={p.file_name ?? "photo"}
              className="h-full w-full object-cover"
              loading="lazy"
            />
          </div>
        );
      })}
    </div>
  );
}
