import { useEffect, type ReactNode } from "react";
import { Sidebar } from "./Sidebar";
import { Topbar } from "./Topbar";
import type { UserPublic } from "@/types";
import { cn } from "@/lib/utils";
import { useSidebarStore } from "@/stores/sidebarStore";

export function AppLayout({ user, children }: { user: UserPublic; children: ReactNode }) {
  const open = useSidebarStore((s) => s.open);
  const toggle = useSidebarStore((s) => s.toggle);

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "b" && (event.metaKey || event.ctrlKey)) {
        event.preventDefault();
        toggle();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [toggle]);

  return (
    <div className="min-h-screen w-full bg-background">
      <Sidebar role={user.role} />
      <div
        className={cn(
          "flex min-h-screen min-w-0 flex-col transition-[margin] duration-200 ease-in-out",
          open && "md:ml-64",
        )}
      >
        <Topbar user={user} />
        <main className="flex-1 overflow-x-hidden p-4 md:p-6">{children}</main>
      </div>
    </div>
  );
}
