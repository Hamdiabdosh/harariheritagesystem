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
import { listRecords } from "@/api/records";
import { StatCard } from "./StatCard";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { EmptyState } from "@/components/common/EmptyState";
import { useAuthStore } from "@/stores/authStore";
import type { RecordSummary } from "@/types";

function daysWaiting(createdAt: string): number {
  return Math.floor(
    (Date.now() - new Date(createdAt).getTime()) / (1000 * 60 * 60 * 24),
  );
}

function ageBadgeClasses(days: number): string {
  if (days <= 3) return "bg-emerald-100 text-emerald-800";
  if (days <= 7) return "bg-amber-100 text-amber-800";
  return "bg-red-100 text-red-800";
}

function PendingAgingQueue() {
  const { t } = useTranslation();

  const { data, isLoading } = useQuery({
    queryKey: ["records", "pending_review", "aging"],
    queryFn: () => listRecords({ status: "pending_review", limit: 10, page: 1 }),
  });

  const records = data?.items ?? [];

  return (
    <div className="rounded-xl border border-border bg-card p-4">
      <h2 className="font-amharic text-base font-semibold text-foreground">
        {t("supervisor.waitingTime")}
      </h2>

      {isLoading && (
        <div className="flex justify-center py-8">
          <LoadingSpinner />
        </div>
      )}

      {!isLoading && records.length === 0 && (
        <div className="mt-3">
          <EmptyState title={t("supervisor.noPending")} icon={<ListChecks className="h-6 w-6" />} />
        </div>
      )}

      {!isLoading && records.length > 0 && (
        <div className="mt-3 space-y-2">
          <ul className="divide-y divide-border rounded-lg border border-border">
            {records.map((record) => (
              <PendingAgingRow key={record.id} record={record} />
            ))}
          </ul>
          <Link
            to="/pending"
            className="font-amharic inline-flex items-center gap-1 text-sm font-medium text-primary hover:text-primary/80"
          >
            {t("nav.pendingReview")}
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </div>
      )}
    </div>
  );
}

function PendingAgingRow({ record }: { record: RecordSummary }) {
  const { t } = useTranslation();
  const days = daysWaiting(record.created_at);
  const TypeIcon = record.record_type === "immovable" ? Building2 : Package;

  const badgeLabel =
    days >= 8
      ? `${days} ${t("supervisor.days")} — ${t("supervisor.overdue")}`
      : `${days} ${t("supervisor.days")}`;

  return (
    <li className="flex flex-wrap items-center gap-2 px-3 py-2.5 text-sm sm:flex-nowrap">
      <span
        className={`shrink-0 rounded-md px-2 py-0.5 text-xs font-medium tabular-nums ${ageBadgeClasses(days)}`}
      >
        {badgeLabel}
      </span>
      <span className="shrink-0 text-xs text-muted-foreground tabular-nums">
        {record.record_id}
      </span>
      <span className="font-amharic min-w-0 flex-1 truncate text-foreground">
        {record.name_amharic}
      </span>
      <span className="inline-flex shrink-0 items-center gap-1 rounded-md bg-muted px-2 py-0.5 text-xs text-muted-foreground">
        <TypeIcon className="h-3 w-3" />
        {t(`recordType.${record.record_type}`)}
      </span>
      <Link
        to="/supervisor/records/$type/$id"
        params={{ type: record.record_type, id: record.id }}
        className="font-amharic inline-flex shrink-0 items-center gap-1 text-xs font-medium text-primary hover:text-primary/80"
      >
        {t("actions.viewDetail")}
        <ArrowRight className="h-3 w-3" />
      </Link>
    </li>
  );
}

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

          <PendingAgingQueue />

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
