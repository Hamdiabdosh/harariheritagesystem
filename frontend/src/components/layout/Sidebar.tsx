import { Link, useRouterState } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import {
  LayoutDashboard,
  FileText,
  ListChecks,
  CheckCircle2,
  Users,
  BarChart3,
} from "lucide-react";
import type { Role } from "@/types";
import { cn } from "@/lib/utils";
import { AppLogo } from "@/components/common/AppLogo";
import { getDashboardStats } from "@/api/dashboard";
import { useSidebarStore } from "@/stores/sidebarStore";

type StatusCountKey =
  | "pending_review"
  | "under_review"
  | "draft"
  | "returned"
  | "approved";

interface NavItem {
  to: string;
  labelKey: string;
  icon: typeof LayoutDashboard;
  roles: Role[];
  showCount?: StatusCountKey;
}

const NAV: NavItem[] = [
  {
    to: "/dashboard",
    labelKey: "nav.dashboard",
    icon: LayoutDashboard,
    roles: ["registrar", "supervisor", "manager"],
  },
  { to: "/records", labelKey: "nav.myRecords", icon: FileText, roles: ["registrar"] },
  {
    to: "/supervisor/records",
    labelKey: "nav.allRecords",
    icon: FileText,
    roles: ["supervisor"],
  },
  {
    to: "/pending",
    labelKey: "nav.pendingReview",
    icon: ListChecks,
    roles: ["supervisor"],
    showCount: "pending_review",
  },
  {
    to: "/manager/records",
    labelKey: "nav.allRecords",
    icon: FileText,
    roles: ["manager"],
  },
  {
    to: "/reviewed",
    labelKey: "nav.reviewed",
    icon: CheckCircle2,
    roles: ["manager"],
    showCount: "under_review",
  },
  { to: "/users", labelKey: "nav.users", icon: Users, roles: ["manager"] },
  {
    to: "/reports",
    labelKey: "nav.reports",
    icon: BarChart3,
    roles: ["supervisor", "manager"],
  },
];

export function Sidebar({ role }: { role: Role }) {
  const { t } = useTranslation();
  const open = useSidebarStore((s) => s.open);
  const pathname = useRouterState({ select: (s) => s.location.pathname });

  const { data: stats } = useQuery({
    queryKey: ["dashboard", "stats"],
    queryFn: getDashboardStats,
    enabled: role === "supervisor" || role === "manager",
    staleTime: 60_000,
  });

  const items = NAV.filter((i) => i.roles.includes(role));

  return (
    <aside
      className={cn(
        "fixed inset-y-0 left-0 z-40 hidden h-svh w-64 flex-col border-r border-border bg-sidebar transition-transform duration-200 ease-in-out md:flex",
        open ? "translate-x-0" : "-translate-x-full",
      )}
    >
      <div className="flex h-16 items-center gap-2.5 border-b border-border px-4">
        <AppLogo size="md" />
        <div className="leading-tight">
          <div className="font-amharic text-sm font-bold text-sidebar-foreground">
            {t("app.name")}
          </div>
          <div className="text-[10px] uppercase tracking-wide text-muted-foreground">
            Qirs Mezgeb
          </div>
        </div>
      </div>

      <nav className="flex-1 space-y-1 overflow-y-auto p-3">
        {items.map((item) => {
          const active = pathname === item.to || pathname.startsWith(item.to + "/");
          const Icon = item.icon;
          const count =
            item.showCount && stats ? stats.by_status[item.showCount] : 0;
          const showBadge = Boolean(item.showCount && count > 0);
          return (
            <Link
              key={item.to + item.labelKey}
              to={item.to}
              className={cn(
                "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                active
                  ? "bg-sidebar-accent text-sidebar-accent-foreground"
                  : "text-sidebar-foreground/80 hover:bg-sidebar-accent/60 hover:text-sidebar-accent-foreground",
              )}
            >
              <Icon className="h-4 w-4 shrink-0" />
              <span className="font-amharic flex-1">{t(item.labelKey)}</span>
              {showBadge && (
                <span className="ml-auto rounded-full bg-amber-500 px-1.5 py-0.5 text-[10px] font-bold text-white tabular-nums">
                  {count}
                </span>
              )}
            </Link>
          );
        })}
      </nav>

      <div className="border-t border-border p-3 text-[11px] text-muted-foreground">
        <div className="font-amharic">{t("app.subtitle")}</div>
      </div>
    </aside>
  );
}
