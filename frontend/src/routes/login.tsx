import { createFileRoute, Navigate, useNavigate } from "@tanstack/react-router";
import { useState, type FormEvent } from "react";
import { useTranslation } from "react-i18next";
import { Languages } from "lucide-react";
import { useAuthStore } from "@/stores/authStore";
import { useLanguageStore } from "@/stores/languageStore";
import { login as apiLogin } from "@/api/auth";
import { AppLogo } from "@/components/common/AppLogo";
import i18n from "@/i18n";
import type { Language } from "@/types";

export const Route = createFileRoute("/login")({
  component: LoginPage,
});

function LoginPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const hydrated = useAuthStore((s) => s.hydrated);
  const setSession = useAuthStore((s) => s.setSession);
  const language = useLanguageStore((s) => s.language);
  const setLanguage = useLanguageStore((s) => s.setLanguage);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  if (hydrated && isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const result = await apiLogin(email, password);
      setSession(result);
      if (result.user.language) {
        setLanguage(result.user.language);
        void i18n.changeLanguage(result.user.language);
      }
      navigate({ to: "/dashboard", replace: true });
    } catch {
      setError(t("auth.invalidCredentials"));
    } finally {
      setLoading(false);
    }
  };

  const toggleLang = () => {
    const next: Language = language === "am" ? "en" : "am";
    setLanguage(next);
    void i18n.changeLanguage(next);
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-secondary to-background px-4">
      <div className="absolute right-4 top-4">
        <button
          onClick={toggleLang}
          className="inline-flex items-center gap-1.5 rounded-md border border-border bg-card px-3 py-1.5 text-xs font-medium hover:bg-accent"
        >
          <Languages className="h-3.5 w-3.5" />
          {language === "am" ? "English" : "አማርኛ"}
        </button>
      </div>

      <div className="w-full max-w-md">
        <div className="mb-6 flex flex-col items-center text-center">
          <AppLogo size="lg" className="mb-3 drop-shadow-sm" />
          <h1 className="font-amharic text-2xl font-bold text-foreground">{t("app.name")}</h1>
          <p className="mt-1 font-amharic text-sm text-muted-foreground">{t("app.subtitle")}</p>
        </div>

        <div className="rounded-xl border border-border bg-card p-6 shadow-sm">
          <h2 className="font-amharic text-lg font-semibold text-foreground">
            {t("auth.loginTitle")}
          </h2>
          <p className="mt-1 font-amharic text-sm text-muted-foreground">
            {t("auth.loginSubtitle")}
          </p>

          <form onSubmit={handleSubmit} className="mt-5 space-y-4">
            <div>
              <label className="font-amharic mb-1.5 block text-sm font-medium text-foreground">
                {t("auth.email")}
                <span className="ml-1 text-xs text-muted-foreground">/ Email</span>
              </label>
              <input
                type="email"
                required
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm outline-none transition focus:border-ring focus:ring-2 focus:ring-ring/30"
                placeholder="name@harari.gov.et"
              />
            </div>

            <div>
              <label className="font-amharic mb-1.5 block text-sm font-medium text-foreground">
                {t("auth.password")}
                <span className="ml-1 text-xs text-muted-foreground">/ Password</span>
              </label>
              <input
                type="password"
                required
                autoComplete="current-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm outline-none transition focus:border-ring focus:ring-2 focus:ring-ring/30"
              />
            </div>

            {error && (
              <div className="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="font-amharic w-full rounded-md bg-primary px-4 py-2.5 text-sm font-semibold text-primary-foreground shadow transition hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {loading ? t("auth.signingIn") : t("auth.signIn")}
            </button>
          </form>
        </div>

        <p className="mt-6 text-center text-xs text-muted-foreground">
          Harari Regional State Culture, Heritage and Tourism Bureau
        </p>
      </div>
    </div>
  );
}
