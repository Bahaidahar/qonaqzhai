import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/shared/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium",
  {
    variants: {
      variant: {
        default:
          "bg-[var(--color-secondary)] text-[var(--color-secondary-foreground)]",
        primary:
          "bg-[var(--color-primary)]/15 text-[var(--color-primary)]",
        success:
          "bg-[var(--color-success)]/15 text-[var(--color-success)]",
        warning:
          "bg-amber-500/15 text-amber-700 dark:text-amber-300",
        danger:
          "bg-[var(--color-destructive)]/15 text-[var(--color-destructive)]",
        outline: "border bg-transparent",
      },
    },
    defaultVariants: { variant: "default" },
  }
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <span className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}
