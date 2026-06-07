import { useTranslation } from "react-i18next";
import type { RecordStatus } from "@/types";
import { cn } from "@/lib/utils";

const styles: Record<RecordStatus, string> = {
  draft: "bg-muted text-muted-foreground border-border",
  pending_review: "bg-amber-100 text-amber-900 border-amber-200",
  under_review: "bg-blue-100 text-blue-900 border-blue-200",
  returned: "bg-rose-100 text-rose-900 border-rose-200",
  approved: "bg-emerald-100 text-emerald-900 border-emerald-200",
};

export function StatusBadge({
  status,
  className,
}: {
  status: RecordStatus;
  className?: string;
}) {
  const { t } = useTranslation();
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium",
        styles[status],
        className,
      )}
    >
      {t(`status.${status}`)}
    </span>
  );
}
