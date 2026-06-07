import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { UserPublic } from "@/types";

interface AuthStoreState {
  // ── Persisted to localStorage (safe — not sensitive tokens) ──
  user: UserPublic | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  hydrated: boolean;

  // ── Memory only — NEVER written to localStorage ──
  accessToken: string | null;

  // ── Actions ──
  setSession: (s: { user: UserPublic; access_token: string; refresh_token: string }) => void;
  setAccessToken: (token: string) => void;
  setUser: (user: UserPublic) => void;
  logout: () => void;
  _setHydrated: (v: boolean) => void;
}

// Safe storage — falls back to noop on the server during SSR
const safeStorage = createJSONStorage(() => {
  if (typeof window === "undefined") {
    return {
      getItem: () => null,
      setItem: () => {},
      removeItem: () => {},
    };
  }
  return window.localStorage;
});

export const useAuthStore = create<AuthStoreState>()(
  persist(
    (set) => ({
      // Persisted state
      user: null,
      refreshToken: null,
      isAuthenticated: false,
      hydrated: false,

      // Memory-only — starts null on every page load
      // The silent-refresh interceptor in client.ts will populate this
      // automatically on the first API call after a page reload.
      accessToken: null,

      setSession: ({ user, access_token, refresh_token }) =>
        set({
          user,
          accessToken: access_token,   // memory only — not persisted
          refreshToken: refresh_token, // persisted — used to get new access tokens
          isAuthenticated: true,
        }),

      setAccessToken: (token) => set({ accessToken: token }),

      setUser: (user) => set({ user }),

      logout: () =>
        set({
          user: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
        }),

      _setHydrated: (v) => set({ hydrated: v }),
    }),
    {
      name: "qirs-auth",
      storage: safeStorage,

      // ── CRITICAL: accessToken is deliberately excluded here ──
      // Only user identity, refresh token, and auth flag are persisted.
      // accessToken lives in memory only — gone on page reload.
      // The Axios interceptor calls /auth/refresh on the first 401
      // to get a new access token transparently.
      partialize: (s) => ({
        user: s.user,
        refreshToken: s.refreshToken,
        isAuthenticated: s.isAuthenticated,
        // accessToken intentionally omitted
      }),

      onRehydrateStorage: () => (state) => {
        state?._setHydrated(true);
      },
    },
  ),
);
