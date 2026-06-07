import { createFileRoute, Navigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useAuthStore } from "@/stores/authStore";
import { RecordsList } from "@/components/records/RecordsList";
import { listRecords } from "@/api/records";

export const Route = createFileRoute("/_authenticated/manager/records/")({
  component: ManagerAllRecordsPage,
});

function ManagerAllRecordsPage() {
  const { t } = useTranslation();
  const role = useAuthStore((s) => s.user?.role);
  if (role && role !== "manager") return <Navigate to="/unauthorized" replace />;

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t("nav.allRecords")}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t("manager.allSubtitle")}
        </p>
      </div>
      <RecordsList
        queryKey={["records", "manager-all"]}
        fetcher={listRecords}
        detailPath="/manager/records"
      />
    </div>
  );
}
