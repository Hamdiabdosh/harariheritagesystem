import { createFileRoute, useNavigate, Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { ArrowLeft } from "lucide-react";
import { ImmovableForm } from "@/components/forms/ImmovableForm";
import { MovableForm } from "@/components/forms/MovableForm";

export const Route = createFileRoute("/_authenticated/records/new/$type")({
  component: NewRecordPage,
});

function NewRecordPage() {
  const { type } = Route.useParams();
  const { t } = useTranslation();
  const navigate = useNavigate();

  if (type !== "immovable" && type !== "movable") {
    navigate({ to: "/records" });
    return null;
  }

  return (
    <div className="space-y-4">
      <Link
        to="/records"
        className="font-amharic inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="h-4 w-4" />
        {t("nav.myRecords")}
      </Link>

      {type === "immovable" ? (
        <ImmovableForm mode="create" />
      ) : (
        <MovableForm mode="create" />
      )}
    </div>
  );
}
