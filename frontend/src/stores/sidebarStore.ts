import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

interface SidebarStoreState {
  open: boolean;
  toggle: () => void;
  setOpen: (open: boolean) => void;
}

const safeStorage = createJSONStorage(() => {
  if (typeof window === "undefined") {
    return { getItem: () => null, setItem: () => {}, removeItem: () => {} };
  }
  return window.localStorage;
});

export const useSidebarStore = create<SidebarStoreState>()(
  persist(
    (set) => ({
      open: true,
      toggle: () => set((state) => ({ open: !state.open })),
      setOpen: (open) => set({ open }),
    }),
    { name: "qirs-sidebar", storage: safeStorage },
  ),
);
