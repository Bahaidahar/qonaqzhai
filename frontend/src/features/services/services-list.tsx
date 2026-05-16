"use client";

import { useEffect, useState } from "react";
import { Check } from "lucide-react";
import { api, type Service, type ServiceUnit } from "@/shared/api";
import { formatKZT } from "@/shared/lib/utils";

const UNIT_SUFFIX: Record<ServiceUnit, string> = {
  fixed: "",
  hour: " / hr",
  item: " / item",
  person: " / person",
  day: " / day",
};

interface Props {
  vendorId: string;
  selectedId?: string | null;
  onSelect?: (s: Service | null) => void;
  showFallbackEmpty?: boolean;
}

/** Customer-facing list of vendor services with optional selection. */
export function ServicesList({ vendorId, selectedId, onSelect, showFallbackEmpty }: Props) {
  const [items, setItems] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let alive = true;
    api
      .vendorServices(vendorId)
      .then((r) => {
        if (!alive) return;
        setItems((r.items ?? []).filter((s) => s.isActive));
      })
      .finally(() => {
        if (alive) setLoading(false);
      });
    return () => {
      alive = false;
    };
  }, [vendorId]);

  if (loading) {
    return (
      <div className="text-sm text-[var(--color-muted-foreground)]">
        Loading services…
      </div>
    );
  }

  if (items.length === 0) {
    if (!showFallbackEmpty) return null;
    return (
      <div className="rounded-xl border border-dashed py-6 text-center text-sm text-[var(--color-muted-foreground)]">
        Vendor has not published any services.
      </div>
    );
  }

  return (
    <ul className="space-y-2">
      {items.map((s) => {
        const active = s.id === selectedId;
        const Tag = onSelect ? "button" : "div";
        return (
          <li key={s.id}>
            <Tag
              type={onSelect ? "button" : undefined}
              onClick={onSelect ? () => onSelect(active ? null : s) : undefined}
              className={`flex w-full items-start justify-between gap-3 rounded-xl border p-4 text-left transition ${
                onSelect ? "hover:bg-[var(--color-muted)]" : ""
              } ${active ? "border-[var(--color-primary)] bg-[var(--color-primary)]/8" : ""}`}
            >
              <div className="min-w-0">
                <h4 className="font-medium">{s.name}</h4>
                {s.description && (
                  <p className="mt-1 text-xs text-[var(--color-muted-foreground)]">
                    {s.description}
                  </p>
                )}
              </div>
              <div className="flex shrink-0 items-center gap-2">
                <span className="text-sm font-semibold">
                  {formatKZT(s.price)}
                  <span className="font-mono text-[10px] text-[var(--color-muted-foreground)]">
                    {UNIT_SUFFIX[s.unit]}
                  </span>
                </span>
                {active && <Check className="h-4 w-4 text-[var(--color-primary)]" />}
              </div>
            </Tag>
          </li>
        );
      })}
    </ul>
  );
}
