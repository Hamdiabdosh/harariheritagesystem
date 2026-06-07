import { createFileRoute, Navigate } from "@tanstack/react-router";
import { useAuthStore } from "@/stores/authStore";
import { ReportsPage } from "@/components/reports/ReportsPage";

export const Route = createFileRoute("/_authenticated/reports")({
  component: Reports,
});

function Reports() {
  const role = useAuthStore((s) => s.user?.role);
  if (role && role !== "manager" && role !== "supervisor") {
    return <Navigate to="/unauthorized" replace />;
  }
  return <ReportsPage />;
}
