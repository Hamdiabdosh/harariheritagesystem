import { Link } from "@tanstack/react-router";
import { Building2, Package, MapPin, Calendar } from "lucide-react";
import { useTranslation } from "react-i18next";
import { StatusBadge } from "@/components/common/StatusBadge";
import type { RecordSummary } from "@/types";

function formatDate(iso: string, locale: string) {
  try {
    return new Date(iso).toLocaleDateString(locale === "am" ? "am-ET" : "en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  } catch {
    return iso;
  }
}

export function RecordCard({
  record,
  detailPath,
}: {
  record: RecordSummary;
  detailPath?: string;
}) {
  const { t, i18n } = useTranslation();
  const Icon = record.record_type === "immovable" ? Building2 : Package;

  const linkProps = detailPath
    ? detailPath === "/manager/records"
      ? {
          to: "/manager/records/$type/$id" as const,
          params: { type: record.record_type, id: record.id },
        }
      : {
          to: "/supervisor/records/$type/$id" as const,
          params: { type: record.record_type, id: record.id },
        }
    : {
        to: "/records/$id/edit" as const,
        params: { id: record.id },
      };

  return (
    <Link
      {...linkProps}
      className="block rounded-xl border border-border bg-card p-4 transition-all hover:border-primary/40 hover:shadow-sm"
    >
      <div className="flex items-start gap-3">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <Icon className="h-5 w-5" />
        </div>
        <div className="min-w-0 flex-1">
          <div className="flex items-start justify-between gap-2">
            <div className="min-w-0">
              <div className="font-amharic truncate text-base font-semibold text-foreground">
                {record.name_amharic}
              </div>
              <div className="mt-0.5 text-xs text-muted-foreground tabular-nums">
                {record.record_id}
              </div>
            </div>
            <StatusBadge status={record.status} />
          </div>
          <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
            {(record.woreda || record.kebele) && (
              <span className="font-amharic inline-flex items-center gap-1">
                <MapPin className="h-3 w-3" />
                {[record.woreda, record.kebele].filter(Boolean).join(" / ")}
              </span>
            )}
            <span className="inline-flex items-center gap-1">
              <Calendar className="h-3 w-3" />
              {formatDate(record.updated_at, i18n.language)}
            </span>
            <span className="font-amharic capitalize">
              {t(`recordType.${record.record_type}`)}
            </span>
          </div>
        </div>
      </div>
    </Link>
  );
}
