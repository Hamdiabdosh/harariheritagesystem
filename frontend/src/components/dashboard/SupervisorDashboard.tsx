import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { Link } from "@tanstack/react-router";
import {
  ListChecks,
  Eye,
  CheckCircle2,
  Building2,
  Package,
  ArrowRight,
} from "lucide-react";
import { getDashboardStats } from "@/api/dashboard";
import { StatCard } from "./StatCard";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { useAuthStore } from "@/stores/authStore";

export function SupervisorDashboard() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["dashboard", "stats"],
    queryFn: getDashboardStats,
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t("nav.dashboard")}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {user?.full_name} · {t("roles.supervisor")}
        </p>
      </div>

      {isLoading && (
        <div className="flex justify-center rounded-xl border border-border bg-card p-12">
          <LoadingSpinner />
        </div>
      )}

      {isError && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-900">
          <div className="font-medium">{t("dashboard.loadError")}</div>
          <div className="mt-1 text-xs opacity-80">{(error as Error)?.message}</div>
          <button
            onClick={() => refetch()}
            className="mt-3 inline-flex items-center rounded-md border border-rose-300 bg-white px-3 py-1.5 text-xs font-medium text-rose-900 hover:bg-rose-100"
          >
            {t("common.retry")}
          </button>
        </div>
      )}

      {data && (
        <>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              label={t("status.pending_review")}
              value={data.by_status.pending_review}
              icon={<ListChecks className="h-5 w-5" />}
              tone="warning"
            />
            <StatCard
              label={t("status.under_review")}
              value={data.by_status.under_review}
              icon={<Eye className="h-5 w-5" />}
              tone="info"
            />
            <StatCard
              label={t("status.approved")}
              value={data.by_status.approved}
              icon={<CheckCircle2 className="h-5 w-5" />}
              tone="success"
            />
          </div>

          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
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
          </div>

          <Link
            to="/pending"
            className="font-amharic flex w-full items-center justify-center gap-2 rounded-xl bg-primary px-4 py-3 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary/90"
          >
            {t("supervisor.pendingCta")}
            <ArrowRight className="h-4 w-4" />
          </Link>
        </>
      )}
    </div>
  );
}
