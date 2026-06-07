import { Search, X } from "lucide-react";
import { useTranslation } from "react-i18next";
import type { RecordStatus, RecordType, UnifiedListParams } from "@/types";

interface RecordFiltersProps {
  value: UnifiedListParams;
  onChange: (next: UnifiedListParams) => void;
}

const STATUSES: RecordStatus[] = [
  "draft",
  "pending_review",
  "under_review",
  "returned",
  "approved",
];
const TYPES: RecordType[] = ["immovable", "movable"];

export function RecordFilters({ value, onChange }: RecordFiltersProps) {
  const { t } = useTranslation();
  const hasActive = Boolean(
    value.search || value.status || value.type || value.woreda,
  );

  return (
    <div className="space-y-3 rounded-xl border border-border bg-card p-3">
      <div className="flex flex-wrap items-center gap-2">
        <div className="relative min-w-[200px] flex-1">
          <Search className="absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <input
            type="text"
            placeholder={t("filters.searchPlaceholder")}
            value={value.search ?? ""}
            onChange={(e) =>
              onChange({ ...value, search: e.target.value || undefined, page: 1 })
            }
            className="font-amharic h-9 w-full rounded-md border border-input bg-background pl-8 pr-3 text-sm placeholder:text-muted-foreground focus:border-primary focus:outline-none"
          />
        </div>

        <select
          value={value.type ?? ""}
          onChange={(e) =>
            onChange({
              ...value,
              type: (e.target.value || undefined) as RecordType | undefined,
              page: 1,
            })
          }
          className="font-amharic h-9 rounded-md border border-input bg-background px-2 text-sm"
        >
          <option value="">{t("filters.allTypes")}</option>
          {TYPES.map((tp) => (
            <option key={tp} value={tp}>
              {t(`recordType.${tp}`)}
            </option>
          ))}
        </select>

        <select
          value={value.status ?? ""}
          onChange={(e) =>
            onChange({
              ...value,
              status: (e.target.value || undefined) as RecordStatus | undefined,
              page: 1,
            })
          }
          className="font-amharic h-9 rounded-md border border-input bg-background px-2 text-sm"
        >
          <option value="">{t("filters.allStatuses")}</option>
          {STATUSES.map((s) => (
            <option key={s} value={s}>
              {t(`status.${s}`)}
            </option>
          ))}
        </select>

        <input
          type="text"
          placeholder={t("filters.woreda")}
          value={value.woreda ?? ""}
          onChange={(e) =>
            onChange({ ...value, woreda: e.target.value || undefined, page: 1 })
          }
          className="font-amharic h-9 w-32 rounded-md border border-input bg-background px-2 text-sm placeholder:text-muted-foreground"
        />

        {hasActive && (
          <button
            type="button"
            onClick={() =>
              onChange({ page: 1, limit: value.limit ?? 20 })
            }
            className="inline-flex h-9 items-center gap-1 rounded-md border border-input bg-background px-2 text-xs text-muted-foreground hover:bg-accent"
          >
            <X className="h-3.5 w-3.5" />
            {t("filters.clear")}
          </button>
        )}
      </div>
    </div>
  );
}
