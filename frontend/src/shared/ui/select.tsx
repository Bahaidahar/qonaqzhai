"use client";

import {
  useEffect,
  useRef,
  useState,
  type KeyboardEvent,
  type ReactNode,
} from "react";
import { Check, ChevronDown } from "lucide-react";
import { cn } from "@/shared/lib/utils";

export interface SelectOption {
  value: string;
  label: ReactNode;
}

interface SelectProps {
  value: string;
  onChange: (value: string) => void;
  options: SelectOption[] | string[];
  placeholder?: string;
  className?: string;
  disabled?: boolean;
  "aria-label"?: string;
}

function normalize(opts: SelectOption[] | string[]): SelectOption[] {
  return opts.map((o) =>
    typeof o === "string" ? { value: o, label: o } : o
  );
}

export function Select({
  value,
  onChange,
  options,
  placeholder,
  className,
  disabled,
  "aria-label": ariaLabel,
}: SelectProps) {
  const opts = normalize(options);
  const [open, setOpen] = useState(false);
  const [focusIdx, setFocusIdx] = useState(() =>
    Math.max(
      0,
      opts.findIndex((o) => o.value === value)
    )
  );
  const rootRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLUListElement>(null);

  const selected = opts.find((o) => o.value === value);

  useEffect(() => {
    if (!open) return;
    function onClickOutside(e: MouseEvent) {
      if (!rootRef.current?.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", onClickOutside);
    return () => document.removeEventListener("mousedown", onClickOutside);
  }, [open]);

  useEffect(() => {
    if (open && listRef.current) {
      const el = listRef.current.children[focusIdx] as HTMLElement | undefined;
      el?.scrollIntoView({ block: "nearest" });
    }
  }, [open, focusIdx]);

  function pick(idx: number) {
    const opt = opts[idx];
    if (!opt) return;
    onChange(opt.value);
    setOpen(false);
  }

  function onKey(e: KeyboardEvent<HTMLButtonElement>) {
    if (disabled) return;
    if (!open) {
      if (e.key === "Enter" || e.key === " " || e.key === "ArrowDown") {
        e.preventDefault();
        setOpen(true);
      }
      return;
    }
    if (e.key === "Escape") {
      e.preventDefault();
      setOpen(false);
      return;
    }
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setFocusIdx((i) => Math.min(opts.length - 1, i + 1));
      return;
    }
    if (e.key === "ArrowUp") {
      e.preventDefault();
      setFocusIdx((i) => Math.max(0, i - 1));
      return;
    }
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      pick(focusIdx);
      return;
    }
    if (e.key === "Home") {
      e.preventDefault();
      setFocusIdx(0);
      return;
    }
    if (e.key === "End") {
      e.preventDefault();
      setFocusIdx(opts.length - 1);
    }
  }

  return (
    <div ref={rootRef} className={cn("relative", className)}>
      <button
        type="button"
        onClick={() => !disabled && setOpen((o) => !o)}
        onKeyDown={onKey}
        disabled={disabled}
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-label={ariaLabel}
        className={cn(
          "flex h-11 w-full items-center justify-between gap-2 rounded-xl border bg-[var(--color-input)]/60 px-4 text-sm transition",
          "focus:outline-none focus:ring-2 focus:ring-[var(--color-ring)]",
          open && "ring-2 ring-[var(--color-ring)]",
          disabled && "cursor-not-allowed opacity-50",
          !selected && "text-[var(--color-muted-foreground)]"
        )}
      >
        <span className="truncate">
          {selected?.label ?? placeholder ?? ""}
        </span>
        <ChevronDown
          className={cn(
            "h-4 w-4 shrink-0 text-[var(--color-muted-foreground)] transition-transform",
            open && "rotate-180"
          )}
        />
      </button>

      {open && (
        <ul
          ref={listRef}
          role="listbox"
          className="absolute z-50 mt-1 max-h-64 w-full overflow-auto rounded-xl border bg-[var(--color-card)] py-1 shadow-lg"
        >
          {opts.map((opt, idx) => {
            const active = opt.value === value;
            const focused = idx === focusIdx;
            return (
              <li
                key={opt.value}
                role="option"
                aria-selected={active}
                onMouseEnter={() => setFocusIdx(idx)}
                onClick={() => pick(idx)}
                className={cn(
                  "flex cursor-pointer items-center justify-between gap-2 px-3 py-2 text-sm",
                  focused && "bg-[var(--color-muted)]",
                  active && "text-[var(--color-primary)]"
                )}
              >
                <span className="truncate">{opt.label}</span>
                {active && <Check className="h-4 w-4 shrink-0" />}
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
