import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import {
  Building2,
  Package,
  FileEdit,
  Clock,
  RotateCcw,
  CheckCircle2,
  Plus,
} from "lucide-react";
import { getDashboardStats } from "@/api/dashboard";
import { StatCard } from "./StatCard";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { useAuthStore } from "@/stores/authStore";

export function RegistrarDashboard() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["dashboard", "stats", "registrar"],
    queryFn: getDashboardStats,
  });

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="font-amharic text-2xl font-bold text-foreground">
            {t("nav.dashboard")}
          </h1>
          <p className="font-amharic mt-1 text-sm text-muted-foreground">
            {user?.full_name} · {user && t(`roles.${user.role}`)}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Link
            to="/records/new/$type"
            params={{ type: "immovable" }}
            className="inline-flex items-center gap-2 rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
          >
            <Plus className="h-4 w-4" />
            <span className="font-amharic">{t("nav.newImmovable")}</span>
          </Link>
          <Link
            to="/records/new/$type"
            params={{ type: "movable" }}
            className="inline-flex items-center gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm font-medium text-foreground transition-colors hover:bg-accent"
          >
            <Plus className="h-4 w-4" />
            <span className="font-amharic">{t("nav.newMovable")}</span>
          </Link>
        </div>
      </div>

      {/* Loading / error */}
      {isLoading && (
        <div className="flex justify-center rounded-xl border border-border bg-card p-12">
          <LoadingSpinner />
        </div>
      )}

      {isError && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-900">
          <div className="font-medium">{t("dashboard.loadError")}</div>
          <div className="mt-1 text-xs opacity-80">
            {(error as Error)?.message}
          </div>
          <button
            onClick={() => refetch()}
            className="mt-3 inline-flex items-center rounded-md border border-rose-300 bg-white px-3 py-1.5 text-xs font-medium text-rose-900 hover:bg-rose-100"
          >
            {t("common.retry")}
          </button>
        </div>
      )}

      {/* Stats */}
      {data && (
        <>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              label={t("dashboard.totalImmovable")}
              value={data.total_immovable}
              icon={<Building2 className="h-5 w-5" />}
              tone="primary"
            />
            <StatCard
              label={t("dashboard.totalMovable")}
              value={data.total_movable}
              icon={<Package className="h-5 w-5" />}
              tone="primary"
            />
            <StatCard
              label={t("dashboard.totalRecords")}
              value={data.total_immovable + data.total_movable}
              icon={<FileEdit className="h-5 w-5" />}
            />
          </div>

          <div>
            <h2 className="font-amharic mb-3 text-sm font-semibold text-foreground">
              {t("dashboard.byStatus")}
            </h2>
            <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
              <StatCard
                label={t("status.draft")}
                value={data.by_status.draft}
                icon={<FileEdit className="h-5 w-5" />}
              />
              <StatCard
                label={t("status.pending_review")}
                value={data.by_status.pending_review}
                icon={<Clock className="h-5 w-5" />}
                tone="warning"
              />
              <StatCard
                label={t("status.under_review")}
                value={data.by_status.under_review}
                icon={<Clock className="h-5 w-5" />}
                tone="info"
              />
              <StatCard
                label={t("status.returned")}
                value={data.by_status.returned}
                icon={<RotateCcw className="h-5 w-5" />}
                tone="danger"
              />
              <StatCard
                label={t("status.approved")}
                value={data.by_status.approved}
                icon={<CheckCircle2 className="h-5 w-5" />}
                tone="success"
              />
            </div>
          </div>
        </>
      )}
    </div>
  );
}
