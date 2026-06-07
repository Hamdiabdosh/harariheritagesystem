import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { UserPublic } from "@/types";

interface AuthStoreState {
  user: UserPublic | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  hydrated: boolean;
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
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      hydrated: false,
      setSession: ({ user, access_token, refresh_token }) =>
        set({
          user,
          accessToken: access_token,
          refreshToken: refresh_token,
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
      partialize: (s) => ({
        user: s.user,
        accessToken: s.accessToken,
        refreshToken: s.refreshToken,
        isAuthenticated: s.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        state?._setHydrated(true);
      },
    },
  ),
);
