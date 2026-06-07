import { createFileRoute } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { RecordsList } from "@/components/records/RecordsList";
import { listMyRecords } from "@/api/records";

export const Route = createFileRoute("/_authenticated/records/")({
  component: RecordsPage,
});

function RecordsPage() {
  const { t } = useTranslation();
  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t("nav.myRecords")}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t("records.mySubtitle")}
        </p>
      </div>
      <RecordsList queryKey={["records", "mine"]} fetcher={listMyRecords} />
    </div>
  );
}
