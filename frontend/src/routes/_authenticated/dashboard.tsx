import { createFileRoute } from "@tanstack/react-router";
import { useAuthStore } from "@/stores/authStore";
import { RegistrarDashboard } from "@/components/dashboard/RegistrarDashboard";
import { SupervisorDashboard } from "@/components/dashboard/SupervisorDashboard";
import { ManagerDashboard } from "@/components/dashboard/ManagerDashboard";

export const Route = createFileRoute("/_authenticated/dashboard")({
  component: DashboardPage,
});

function DashboardPage() {
  const user = useAuthStore((s) => s.user);
  if (!user) return null;
  if (user.role === "supervisor") return <SupervisorDashboard />;
  if (user.role === "manager") return <ManagerDashboard />;
  return <RegistrarDashboard />;
}
