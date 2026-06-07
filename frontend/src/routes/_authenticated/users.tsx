import { createFileRoute, Navigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useAuthStore } from "@/stores/authStore";
import { UserManagement } from "@/components/users/UserManagement";

export const Route = createFileRoute("/_authenticated/users")({
  component: UsersPage,
});

function UsersPage() {
  const { t } = useTranslation();
  const role = useAuthStore((s) => s.user?.role);
  if (role && role !== "manager") return <Navigate to="/unauthorized" replace />;

  return (
    <div className="space-y-4">
      <div>
        <h1 className="font-amharic text-2xl font-bold text-foreground">
          {t("nav.users")}
        </h1>
        <p className="font-amharic mt-1 text-sm text-muted-foreground">
          {t("manager.usersSubtitle")}
        </p>
      </div>
      <UserManagement />
    </div>
  );
}
