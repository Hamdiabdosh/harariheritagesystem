import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { getDashboardStats } from "@/api/dashboard";
import { exportCSV, downloadBlob } from "@/api/export";
import type { RecordStatus, RecordType, UnifiedListParams } from "@/types";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { StatCard } from "@/components/dashboard/StatCard";
import { Building2, Package } from "lucide-react";

const STATUS_BAR_COLORS: Record<keyof typeof defaultCounts, string> = {
  draft: "bg-muted-foreground/60",
  pending_review: "bg-amber-500",
  under_review: "bg-blue-500",
  returned: "bg-rose-500",
  approved: "bg-emerald-500",
};

const defaultCounts = {
  draft: 0,
  pending_review: 0,
  under_review: 0,
  returned: 0,
  approved: 0,
};

export function ReportsPage() {
  const { t } = useTranslation();
  const [exporting, setExporting] = useState(false);
  const [filters, setFilters] = useState<UnifiedListParams>({});

  const { data: stats, isLoading } = useQuery({
    queryKey: ["dashboard", "stats"],
    queryFn: getDashboardStats,
  });

  const handleExport = async () => {
    setExporting(true);
    try {
      const blob = await exportCSV(filters);
      downloadBlob(blob, "qirs-mezgeb-export.csv");
    } catch {
      toast.error(t("toast.error"));
    } finally {
      setExporting(false);
    }
  };

  const counts = stats?.by_status ?? defaultCounts;
  const total = stats ? stats.total_immovable + stats.total_movable : 0;

  const statusRows: { key: keyof typeof defaultCounts; label: string }[] = [
    { key: "draft", label: t("status.draft") },
    { key: "pending_review", label: t("status.pending_review") },
    { key: "under_review", label: t("status.under_review") },
    { key: "returned", label: t("status.returned") },
    { key: "approved", label: t("status.approved") },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t("reports.title")}
        </h1>
      </div>

      <section className="rounded-xl border border-border bg-card p-4 space-y-4">
        <h2 className="font-amharic text-sm font-semibold text-foreground">
          {t("reports.exportSection")}
        </h2>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <label className="block text-sm">
            <span className="font-amharic text-xs text-muted-foreground">
              {t("reports.filterType")}
            </span>
            <select
              value={filters.type ?? ""}
              onChange={(e) =>
                setFilters((f) => ({
                  ...f,
                  type: (e.target.value || undefined) as RecordType | undefined,
                }))
              }
              className="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="">{t("filters.allTypes")}</option>
              <option value="immovable">{t("recordType.immovable")}</option>
              <option value="movable">{t("recordType.movable")}</option>
            </select>
          </label>

          <label className="block text-sm">
            <span className="font-amharic text-xs text-muted-foreground">
              {t("reports.filterStatus")}
            </span>
            <select
              value={filters.status ?? ""}
              onChange={(e) =>
                setFilters((f) => ({
                  ...f,
                  status: (e.target.value || undefined) as RecordStatus | undefined,
                }))
              }
              className="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            >
              <option value="">{t("filters.allStatuses")}</option>
              <option value="draft">{t("status.draft")}</option>
              <option value="pending_review">{t("status.pending_review")}</option>
              <option value="under_review">{t("status.under_review")}</option>
              <option value="returned">{t("status.returned")}</option>
              <option value="approved">{t("status.approved")}</option>
            </select>
          </label>

          <label className="block text-sm">
            <span className="font-amharic text-xs text-muted-foreground">
              {t("reports.filterWoreda")}
            </span>
            <input
              type="text"
              value={filters.woreda ?? ""}
              onChange={(e) =>
                setFilters((f) => ({ ...f, woreda: e.target.value || undefined }))
              }
              className="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            />
          </label>

          <label className="block text-sm">
            <span className="font-amharic text-xs text-muted-foreground">
              {t("reports.filterDateFrom")}
            </span>
            <input
              type="date"
              value={filters.date_from ?? ""}
              onChange={(e) =>
                setFilters((f) => ({ ...f, date_from: e.target.value || undefined }))
              }
              className="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            />
          </label>

          <label className="block text-sm">
            <span className="font-amharic text-xs text-muted-foreground">
              {t("reports.filterDateTo")}
            </span>
            <input
              type="date"
              value={filters.date_to ?? ""}
              onChange={(e) =>
                setFilters((f) => ({ ...f, date_to: e.target.value || undefined }))
              }
              className="mt-1 block w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
            />
          </label>
        </div>

        <button
          type="button"
          onClick={() => void handleExport()}
          disabled={exporting}
          className="font-amharic inline-flex items-center gap-1.5 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
        >
          {exporting ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" />
              {t("reports.exporting")}
            </>
          ) : (
            t("reports.exportButton")
          )}
        </button>
      </section>

      <section className="rounded-xl border border-border bg-card p-4 space-y-4">
        <h2 className="font-amharic text-sm font-semibold text-foreground">
          {t("reports.statsSection")}
        </h2>

        {isLoading && (
          <div className="flex justify-center py-8">
            <LoadingSpinner />
          </div>
        )}

        {stats && (
          <>
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
              <StatCard
                label={t("reports.totalRecords")}
                value={total}
                tone="primary"
              />
              <StatCard
                label={t("dashboard.totalImmovable")}
                value={stats.total_immovable}
                icon={<Building2 className="h-5 w-5" />}
              />
              <StatCard
                label={t("dashboard.totalMovable")}
                value={stats.total_movable}
                icon={<Package className="h-5 w-5" />}
              />
            </div>

            <div className="space-y-3">
              {statusRows.map(({ key, label }) => {
                const count = counts[key];
                const pct = total > 0 ? ((count / total) * 100).toFixed(1) : "0.0";
                return (
                  <div key={key}>
                    <div className="mb-1 flex items-center justify-between text-sm">
                      <span className="font-amharic text-foreground">{label}</span>
                      <span className="tabular-nums text-muted-foreground">
                        {count} ({pct}%)
                      </span>
                    </div>
                    <div className="h-2 overflow-hidden rounded-full bg-muted">
                      <div
                        className={`h-full rounded-full transition-all ${STATUS_BAR_COLORS[key]}`}
                        style={{ width: `${total > 0 ? (count / total) * 100 : 0}%` }}
                      />
                    </div>
                  </div>
                );
              })}
            </div>
          </>
        )}
      </section>
    </div>
  );
}
