import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

interface StatCardProps {
  label: string;
  value: number | string;
  icon?: ReactNode;
  tone?: "default" | "primary" | "warning" | "info" | "success" | "danger";
  hint?: string;
}

const toneStyles: Record<NonNullable<StatCardProps["tone"]>, string> = {
  default: "bg-card border-border",
  primary: "bg-primary/5 border-primary/20",
  warning: "bg-amber-50 border-amber-200",
  info: "bg-blue-50 border-blue-200",
  success: "bg-emerald-50 border-emerald-200",
  danger: "bg-rose-50 border-rose-200",
};

const iconToneStyles: Record<NonNullable<StatCardProps["tone"]>, string> = {
  default: "bg-muted text-muted-foreground",
  primary: "bg-primary/10 text-primary",
  warning: "bg-amber-100 text-amber-900",
  info: "bg-blue-100 text-blue-900",
  success: "bg-emerald-100 text-emerald-900",
  danger: "bg-rose-100 text-rose-900",
};

export function StatCard({
  label,
  value,
  icon,
  tone = "default",
  hint,
}: StatCardProps) {
  return (
    <div
      className={cn(
        "flex items-start gap-3 rounded-xl border p-4 transition-shadow hover:shadow-sm",
        toneStyles[tone],
      )}
    >
      {icon && (
        <div
          className={cn(
            "flex h-10 w-10 shrink-0 items-center justify-center rounded-lg",
            iconToneStyles[tone],
          )}
        >
          {icon}
        </div>
      )}
      <div className="min-w-0 flex-1">
        <div className="font-amharic text-xs text-muted-foreground">{label}</div>
        <div className="mt-1 text-2xl font-bold text-foreground tabular-nums">
          {value}
        </div>
        {hint && (
          <div className="font-amharic mt-0.5 text-xs text-muted-foreground">
            {hint}
          </div>
        )}
      </div>
    </div>
  );
}
