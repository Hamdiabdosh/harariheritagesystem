import { useTranslation } from "react-i18next";
import { useNavigate } from "@tanstack/react-router";
import { LogOut, Languages } from "lucide-react";
import type { UserPublic, Language } from "@/types";
import { useAuthStore } from "@/stores/authStore";
import { useLanguageStore } from "@/stores/languageStore";
import { logout as apiLogout } from "@/api/auth";
import { updateMyLanguage } from "@/api/users";
import { AppLogo } from "@/components/common/AppLogo";
import i18n from "@/i18n";

export function Topbar({ user }: { user: UserPublic }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const refreshToken = useAuthStore((s) => s.refreshToken);
  const clearSession = useAuthStore((s) => s.logout);
  const language = useLanguageStore((s) => s.language);
  const setLanguage = useLanguageStore((s) => s.setLanguage);

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

  return (
    <header className="flex h-16 items-center justify-between border-b border-border bg-card px-4 md:px-6">
      <div className="flex items-center gap-2 md:hidden">
        <AppLogo size="sm" />
        <div className="font-amharic text-sm font-bold">{t("app.name")}</div>
      </div>

      <div className="hidden md:block" />

      <div className="flex items-center gap-3">
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
