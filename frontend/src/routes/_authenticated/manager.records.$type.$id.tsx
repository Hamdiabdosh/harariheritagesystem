import { createFileRoute, Navigate, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { ArrowLeft } from "lucide-react";
import { useAuthStore } from "@/stores/authStore";
import { ManagerRecordDetail } from "@/components/records/ManagerRecordDetail";
import type { RecordType } from "@/types";

export const Route = createFileRoute("/_authenticated/manager/records/$type/$id")({
  component: ManagerRecordDetailPage,
});

function ManagerRecordDetailPage() {
  const { type, id } = Route.useParams();
  const { t } = useTranslation();
  const role = useAuthStore((s) => s.user?.role);

  if (role && role !== "manager") return <Navigate to="/unauthorized" replace />;
  if (type !== "immovable" && type !== "movable") {
    return <Navigate to="/reviewed" replace />;
  }

  return (
    <div className="space-y-4">
      <Link
        to="/reviewed"
        className="font-amharic inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        {t("nav.reviewed")}
      </Link>
      <ManagerRecordDetail recordId={id} recordType={type as RecordType} />
    </div>
  );
}
