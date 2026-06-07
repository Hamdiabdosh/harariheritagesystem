import { cn } from "@/lib/utils";

const sizes = {
  sm: "h-8 w-8",
  md: "h-10 w-10",
  lg: "h-16 w-16",
} as const;

interface AppLogoProps {
  size?: keyof typeof sizes;
  className?: string;
}

export function AppLogo({ size = "md", className }: AppLogoProps) {
  return (
    <img
      src="/logo.jpeg"
      alt="Qirs Mezgeb"
      className={cn("shrink-0 object-contain", sizes[size], className)}
    />
  );
}
