import { createFileRoute, Navigate } from "@tanstack/react-router";
import { useAuthStore } from "@/stores/authStore";
import { useTranslation } from "react-i18next";
import { RecordsList } from "@/components/records/RecordsList";
import { listRecords } from "@/api/records";

export const Route = createFileRoute("/_authenticated/supervisor/records/")({
  component: SupervisorAllRecordsPage,
});

function SupervisorAllRecordsPage() {
  const { t } = useTranslation();
  const role = useAuthStore((s) => s.user?.role);
  if (role && role !== "supervisor") return <Navigate to="/unauthorized" replace />;

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t("nav.allRecords")}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t("supervisor.allSubtitle")}
        </p>
      </div>
      <RecordsList
        queryKey={["records", "all"]}
        fetcher={listRecords}
        detailPath="/supervisor/records"
      />
    </div>
  );
}
