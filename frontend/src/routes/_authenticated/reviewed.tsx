import { createFileRoute, Navigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useAuthStore } from "@/stores/authStore";
import { RecordsList } from "@/components/records/RecordsList";
import { listRecords } from "@/api/records";
import type { UnifiedListParams } from "@/types";

export const Route = createFileRoute("/_authenticated/reviewed")({
  component: ReviewedPage,
});

function ReviewedPage() {
  const { t } = useTranslation();
  const role = useAuthStore((s) => s.user?.role);
  if (role && role !== "manager") return <Navigate to="/unauthorized" replace />;

  const fetcher = (params: UnifiedListParams) =>
    listRecords({ ...params, status: "under_review" });

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t("nav.reviewed")}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t("manager.reviewedSubtitle")}
        </p>
      </div>
      <RecordsList
        queryKey={["records", "under_review"]}
        fetcher={fetcher}
        detailPath="/manager/records"
      />
    </div>
  );
}
