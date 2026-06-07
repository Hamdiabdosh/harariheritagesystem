import type { ReactNode, InputHTMLAttributes, TextareaHTMLAttributes, SelectHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

/**
 * Bilingual label: Amharic on top, English hint below.
 */
export function FieldLabel({
  am,
  en,
  required,
  className,
}: {
  am: string;
  en?: string;
  required?: boolean;
  className?: string;
}) {
  return (
    <label className={cn("mb-1 block text-sm", className)}>
      <span className="font-amharic font-medium text-foreground">{am}</span>
      {en && <span className="ml-1.5 text-xs text-muted-foreground">/ {en}</span>}
      {required && <span className="ml-1 text-destructive">*</span>}
    </label>
  );
}

const baseInput =
  "block w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm outline-none transition focus:border-ring focus:ring-2 focus:ring-ring/30 disabled:cursor-not-allowed disabled:opacity-60";

export function TextInput(
  props: InputHTMLAttributes<HTMLInputElement> & { error?: string },
) {
  const { error, className, ...rest } = props;
  return (
    <>
      <input
        {...rest}
        className={cn(baseInput, error && "border-destructive", className)}
      />
      {error && <p className="mt-1 text-xs text-destructive">{error}</p>}
    </>
  );
}

export function NumberInput(
  props: InputHTMLAttributes<HTMLInputElement> & { error?: string },
) {
  return <TextInput type="number" step="any" {...props} />;
}

export function TextArea(
  props: TextareaHTMLAttributes<HTMLTextAreaElement> & { error?: string },
) {
  const { error, className, ...rest } = props;
  return (
    <>
      <textarea
        rows={3}
        {...rest}
        className={cn(baseInput, "font-amharic resize-y", error && "border-destructive", className)}
      />
      {error && <p className="mt-1 text-xs text-destructive">{error}</p>}
    </>
  );
}

export function Select(
  props: SelectHTMLAttributes<HTMLSelectElement> & {
    options: Array<{ value: string; labelAm: string; labelEn?: string }>;
    placeholderAm?: string;
    error?: string;
  },
) {
  const { options, placeholderAm, error, className, ...rest } = props;
  return (
    <>
      <select
        {...rest}
        className={cn(baseInput, "font-amharic", error && "border-destructive", className)}
      >
        {placeholderAm !== undefined && <option value="">{placeholderAm}</option>}
        {options.map((o) => (
          <option key={o.value} value={o.value}>
            {o.labelAm}
            {o.labelEn ? ` / ${o.labelEn}` : ""}
          </option>
        ))}
      </select>
      {error && <p className="mt-1 text-xs text-destructive">{error}</p>}
    </>
  );
}

/** Wraps a labeled control in a grid cell. */
export function Field({
  am,
  en,
  required,
  children,
  span = 1,
}: {
  am: string;
  en?: string;
  required?: boolean;
  children: ReactNode;
  span?: 1 | 2 | 3;
}) {
  const spanClass =
    span === 3 ? "md:col-span-3" : span === 2 ? "md:col-span-2" : "";
  return (
    <div className={spanClass}>
      <FieldLabel am={am} en={en} required={required} />
      {children}
    </div>
  );
}

/** Section card with collapsible header. */
export function FormSection({
  num,
  am,
  en,
  open,
  onToggle,
  children,
}: {
  num: number;
  am: string;
  en: string;
  open: boolean;
  onToggle: () => void;
  children: ReactNode;
}) {
  return (
    <div className="overflow-hidden rounded-xl border border-border bg-card">
      <button
        type="button"
        onClick={onToggle}
        className="flex w-full items-center justify-between gap-3 border-b border-border bg-muted/30 px-4 py-3 text-left hover:bg-muted/60"
      >
        <div className="flex items-center gap-3 min-w-0">
          <span className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary text-xs font-bold text-primary-foreground">
            {num}
          </span>
          <div className="min-w-0">
            <div className="font-amharic truncate text-sm font-semibold text-foreground">
              {am}
            </div>
            <div className="truncate text-xs text-muted-foreground">{en}</div>
          </div>
        </div>
        <span className="text-xs text-muted-foreground">{open ? "−" : "+"}</span>
      </button>
      {open && (
        <div className="grid grid-cols-1 gap-4 p-4 md:grid-cols-3">{children}</div>
      )}
    </div>
  );
}

/** Checkbox list — multi-select string array. */
export function CheckboxList({
  options,
  value,
  onChange,
}: {
  options: Array<{ value: string; labelAm: string; labelEn?: string }>;
  value: string[];
  onChange: (next: string[]) => void;
}) {
  const toggle = (v: string) => {
    onChange(value.includes(v) ? value.filter((x) => x !== v) : [...value, v]);
  };
  return (
    <div className="flex flex-wrap gap-2">
      {options.map((o) => {
        const active = value.includes(o.value);
        return (
          <button
            type="button"
            key={o.value}
            onClick={() => toggle(o.value)}
            className={cn(
              "rounded-md border px-2.5 py-1.5 text-xs transition",
              active
                ? "border-primary bg-primary/10 text-primary"
                : "border-input bg-background text-foreground hover:bg-accent",
            )}
          >
            <span className="font-amharic">{o.labelAm}</span>
            {o.labelEn && <span className="ml-1 text-muted-foreground">/ {o.labelEn}</span>}
          </button>
        );
      })}
    </div>
  );
}

export function Switch({
  checked,
  onChange,
  am,
  en,
}: {
  checked: boolean;
  onChange: (v: boolean) => void;
  am: string;
  en?: string;
}) {
  return (
    <label className="flex cursor-pointer items-center gap-2 rounded-md border border-input bg-background px-3 py-2">
      <input
        type="checkbox"
        checked={checked}
        onChange={(e) => onChange(e.target.checked)}
        className="h-4 w-4 rounded border-input text-primary focus:ring-primary"
      />
      <span className="font-amharic text-sm text-foreground">{am}</span>
      {en && <span className="text-xs text-muted-foreground">/ {en}</span>}
    </label>
  );
}
