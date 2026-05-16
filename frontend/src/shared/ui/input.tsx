import * as React from "react";
import { cn } from "@/shared/lib/utils";

export type InputProps = React.InputHTMLAttributes<HTMLInputElement>;

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => (
    <input
      ref={ref}
      type={type}
      className={cn(
        "flex h-11 w-full rounded-xl border bg-[var(--color-input)]/60 px-4 py-2 text-sm",
        "placeholder:text-[var(--color-muted-foreground)]",
        "focus:outline-none focus:ring-2 focus:ring-[var(--color-ring)] focus:bg-[var(--color-card)]",
        "transition-colors disabled:cursor-not-allowed disabled:opacity-50",
        className
      )}
      {...props}
    />
  )
);
Input.displayName = "Input";
