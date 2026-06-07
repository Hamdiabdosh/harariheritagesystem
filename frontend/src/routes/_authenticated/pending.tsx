import { createFileRoute, Navigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useAuthStore } from "@/stores/authStore";
import { RecordsList } from "@/components/records/RecordsList";
import { listRecords } from "@/api/records";
import type { UnifiedListParams } from "@/types";

export const Route = createFileRoute("/_authenticated/pending")({
  component: PendingPage,
});

function PendingPage() {
  const { t } = useTranslation();
  const role = useAuthStore((s) => s.user?.role);
  if (role && role !== "supervisor") return <Navigate to="/unauthorized" replace />;

  const fetcher = (params: UnifiedListParams) =>
    listRecords({ ...params, status: "pending_review" });

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t("nav.pendingReview")}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t("supervisor.pendingSubtitle")}
        </p>
      </div>
      <RecordsList
        queryKey={["records", "pending_review"]}
        fetcher={fetcher}
        detailPath="/supervisor/records"
      />
    </div>
  );
}
