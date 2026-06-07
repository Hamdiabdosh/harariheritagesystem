import { createFileRoute, Navigate, Link } from "@tanstack/react-router";
import { useAuthStore } from "@/stores/authStore";
import { useTranslation } from "react-i18next";
import { ArrowLeft } from "lucide-react";
import { SupervisorRecordDetail } from "@/components/records/SupervisorRecordDetail";
import type { RecordType } from "@/types";

export const Route = createFileRoute("/_authenticated/supervisor/records/$type/$id")({
  component: SupervisorRecordDetailPage,
});

function SupervisorRecordDetailPage() {
  const { type, id } = Route.useParams();
  const { t } = useTranslation();
  const role = useAuthStore((s) => s.user?.role);

  if (role && role !== "supervisor") return <Navigate to="/unauthorized" replace />;
  if (type !== "immovable" && type !== "movable") {
    return <Navigate to="/pending" replace />;
  }

  return (
    <div className="space-y-4">
      <Link
        to="/pending"
        className="font-amharic inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        {t("nav.pendingReview")}
      </Link>
      <SupervisorRecordDetail recordId={id} recordType={type as RecordType} />
    </div>
  );
}
