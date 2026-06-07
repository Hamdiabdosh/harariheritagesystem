import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { Language } from "@/types";

interface LanguageStoreState {
  language: Language;
  setLanguage: (lang: Language) => void;
}

const safeStorage = createJSONStorage(() => {
  if (typeof window === "undefined") {
    return { getItem: () => null, setItem: () => {}, removeItem: () => {} };
  }
  return window.localStorage;
});

export const useLanguageStore = create<LanguageStoreState>()(
  persist(
    (set) => ({
      language: "am",
      setLanguage: (language) => set({ language }),
    }),
    { name: "qirs-language", storage: safeStorage },
  ),
);
