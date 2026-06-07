import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { LogOut, Languages, Bell, RotateCcw, ListChecks, Eye } from "lucide-react";
import type { RecordStatus, UserPublic, Language } from "@/types";
import { useAuthStore } from "@/stores/authStore";
import { useLanguageStore } from "@/stores/languageStore";
import { logout as apiLogout } from "@/api/auth";
import { updateMyLanguage } from "@/api/users";
import { getDashboardStats } from "@/api/dashboard";
import { listRecords } from "@/api/records";
import { AppLogo } from "@/components/common/AppLogo";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import i18n from "@/i18n";

function truncate(text: string, max: number): string {
  return text.length > max ? `${text.slice(0, max)}…` : text;
}

function daysAgo(iso: string): number {
  return Math.floor(
    (Date.now() - new Date(iso).getTime()) / (1000 * 60 * 60 * 24),
  );
}

function notificationStatusForRole(role: UserPublic["role"]): RecordStatus | null {
  switch (role) {
    case "registrar":
      return "returned";
    case "supervisor":
      return "pending_review";
    case "manager":
      return "under_review";
    default:
      return null;
  }
}

function NotificationIcon({ status }: { status: RecordStatus }) {
  switch (status) {
    case "returned":
      return <RotateCcw className="h-4 w-4 shrink-0 text-amber-600" />;
    case "pending_review":
      return <ListChecks className="h-4 w-4 shrink-0 text-amber-600" />;
    case "under_review":
      return <Eye className="h-4 w-4 shrink-0 text-blue-600" />;
    default:
      return <Bell className="h-4 w-4 shrink-0 text-muted-foreground" />;
  }
}

export function Topbar({ user }: { user: UserPublic }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const refreshToken = useAuthStore((s) => s.refreshToken);
  const clearSession = useAuthStore((s) => s.logout);
  const language = useLanguageStore((s) => s.language);
  const setLanguage = useLanguageStore((s) => s.setLanguage);

  const [open, setOpen] = useState(false);
  const panelRef = useRef<HTMLDivElement>(null);

  const statusFilter = notificationStatusForRole(user.role);

  const { data: stats } = useQuery({
    queryKey: ["dashboard", "stats"],
    queryFn: getDashboardStats,
    staleTime: 60_000,
  });

  const { data: notifications, isLoading: notificationsLoading } = useQuery({
    queryKey: ["records", "notifications", user.role, statusFilter],
    queryFn: () => listRecords({ status: statusFilter!, limit: 5, page: 1 }),
    enabled: open && statusFilter !== null,
  });

  useEffect(() => {
    if (!open) return;
    const handler = (e: MouseEvent) => {
      if (panelRef.current && !panelRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, [open]);

  const handleLogout = async () => {
    try {
      if (refreshToken) await apiLogout(refreshToken);
    } catch {
      // ignore
    }
    clearSession();
    navigate({ to: "/login", replace: true });
  };

  const toggleLang = () => {
    const next: Language = language === "am" ? "en" : "am";
    setLanguage(next);
    void i18n.changeLanguage(next);
    void updateMyLanguage(next).catch(() => undefined);
  };

  const showDot = (stats?.pending_my_action ?? 0) > 0;
  const items = notifications?.items ?? [];

  return (
    <header className="flex h-16 items-center justify-between border-b border-border bg-card px-4 md:px-6">
      <div className="flex items-center gap-2 md:hidden">
        <AppLogo size="sm" />
        <div className="font-amharic text-sm font-bold">{t("app.name")}</div>
      </div>

      <div className="hidden md:block" />

      <div className="flex items-center gap-3">
        {statusFilter && (
          <div className="relative" ref={panelRef}>
            <button
              type="button"
              onClick={() => setOpen((v) => !v)}
              className="relative inline-flex h-9 w-9 items-center justify-center rounded-md border border-border bg-background text-foreground transition-colors hover:bg-accent"
              aria-label={t("notifications.title")}
              aria-expanded={open}
            >
              <Bell className="h-4 w-4" />
              {showDot && (
                <span className="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-red-500" />
              )}
            </button>

            {open && (
              <div className="absolute right-0 z-50 mt-2 w-80 max-w-sm rounded-xl border border-border bg-card shadow-lg">
                <div className="border-b border-border px-4 py-3">
                  <h3 className="font-amharic text-sm font-semibold text-foreground">
                    {t("notifications.title")}
                  </h3>
                </div>

                <div className="max-h-80 overflow-y-auto">
                  {notificationsLoading && (
                    <div className="flex justify-center py-6">
                      <LoadingSpinner />
                    </div>
                  )}

                  {!notificationsLoading && items.length === 0 && (
                    <p className="font-amharic px-4 py-6 text-center text-sm text-muted-foreground">
                      {t("notifications.empty")}
                    </p>
                  )}

                  {!notificationsLoading &&
                    items.map((record) => (
                      <Link
                        key={record.id}
                        to="/records/$id/edit"
                        params={{ id: record.id }}
                        onClick={() => setOpen(false)}
                        className="flex items-start gap-3 border-b border-border px-4 py-3 transition-colors last:border-b-0 hover:bg-accent/50"
                      >
                        <NotificationIcon status={record.status} />
                        <div className="min-w-0 flex-1">
                          <p className="font-amharic truncate text-sm font-medium text-foreground">
                            {truncate(record.name_amharic, 30)}
                          </p>
                          <p className="text-xs text-muted-foreground tabular-nums">
                            {record.record_id}
                          </p>
                          <p className="mt-0.5 text-xs text-muted-foreground">
                            {t("notifications.timeAgo", {
                              n: daysAgo(record.updated_at),
                            })}
                          </p>
                        </div>
                      </Link>
                    ))}
                </div>

                <div className="border-t border-border px-4 py-2.5">
                  <Link
                    to="/records"
                    onClick={() => setOpen(false)}
                    className="font-amharic text-xs font-medium text-primary hover:text-primary/80"
                  >
                    {t("notifications.viewAll")} →
                  </Link>
                </div>
              </div>
            )}
          </div>
        )}

        <button
          onClick={toggleLang}
          className="inline-flex items-center gap-1.5 rounded-md border border-border bg-background px-2.5 py-1.5 text-xs font-medium text-foreground transition-colors hover:bg-accent"
          aria-label={t("common.language")}
        >
          <Languages className="h-3.5 w-3.5" />
          {language === "am" ? "EN" : "አማ"}
        </button>

        <div className="hidden text-right sm:block">
          <div className="text-sm font-medium text-foreground">{user.full_name}</div>
          <div className="text-[11px] text-muted-foreground">{t(`roles.${user.role}`)}</div>
        </div>

        <div className="flex h-9 w-9 items-center justify-center rounded-full bg-primary text-sm font-semibold text-primary-foreground">
          {user.full_name?.charAt(0)?.toUpperCase() ?? "?"}
        </div>

        <button
          onClick={handleLogout}
          className="inline-flex items-center gap-1.5 rounded-md px-2 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
          aria-label={t("auth.logout")}
        >
          <LogOut className="h-4 w-4" />
        </button>
      </div>
    </header>
  );
}
