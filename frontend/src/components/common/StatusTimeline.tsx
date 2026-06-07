import { useTranslation } from "react-i18next";
import { Clock } from "lucide-react";
import type { StatusHistoryEntry, RecordStatus } from "@/types";
import { StatusBadge } from "@/components/common/StatusBadge";
import { EmptyState } from "@/components/common/EmptyState";

interface StatusTimelineProps {
  history: StatusHistoryEntry[];
}

function formatDate(iso: string, locale: string) {
  try {
    return new Date(iso).toLocaleDateString(locale === "am" ? "am-ET" : "en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch {
    return iso;
  }
}

export function StatusTimeline({ history }: StatusTimelineProps) {
  const { t, i18n } = useTranslation();

  if (history.length === 0) {
    return (
      <EmptyState
        icon={<Clock className="h-6 w-6" />}
        title={t("history.empty")}
      />
    );
  }

  const sorted = [...history].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  );

  return (
    <div className="relative space-y-0">
      <div className="absolute bottom-2 left-[11px] top-2 w-px bg-border" aria-hidden />
      {sorted.map((entry) => (
        <div key={entry.id} className="relative flex gap-4 pb-6 last:pb-0">
          <div className="relative z-10 mt-1.5 h-[22px] w-[22px] shrink-0 rounded-full border-2 border-primary bg-background" />
          <div className="min-w-0 flex-1 space-y-1.5 rounded-lg border border-border bg-card p-3">
            <div className="flex flex-wrap items-center gap-2">
              {entry.from_status && (
                <>
                  <span className="text-xs text-muted-foreground">
                    {t(`status.${entry.from_status as RecordStatus}`)}
                  </span>
                  <span className="text-xs text-muted-foreground">→</span>
                </>
              )}
              <StatusBadge status={entry.to_status as RecordStatus} />
            </div>
            <p className="text-sm font-medium text-foreground">{entry.changed_by_name}</p>
            <p className="text-xs text-muted-foreground">
              {formatDate(entry.created_at, i18n.language)}
            </p>
            {entry.note && (
              <p className="text-sm text-muted-foreground">{entry.note}</p>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
