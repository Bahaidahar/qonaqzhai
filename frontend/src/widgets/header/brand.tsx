import { cn } from "@/shared/lib/utils";

interface BrandProps {
  className?: string;
}

export function Wordmark({ className }: BrandProps) {
  return (
    <span
      className={cn(
        "font-display inline-flex items-center text-[15px] font-bold tracking-[-0.04em]",
        className
      )}
    >
      qonaqzhai
      <span className="ml-0.5 inline-block h-1.5 w-1.5 translate-y-1.5 rounded-full bg-[var(--color-primary)]" />
    </span>
  );
}

export function BrandRow({ className }: BrandProps) {
  return (
    <div className={cn("flex items-center", className)}>
      <Wordmark />
    </div>
  );
}
