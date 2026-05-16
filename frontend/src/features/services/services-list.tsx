"use client";

import { useEffect, useState } from "react";
import { Check } from "lucide-react";
import { api, type Service, type ServiceUnit } from "@/shared/api";
import { formatKZT } from "@/shared/lib/utils";
import { useI18n } from "@/shared/i18n/context";
import type { DictKey } from "@/shared/i18n/dict";

const UNIT_SUFFIX_KEY: Record<ServiceUnit, DictKey | null> = {
  fixed: null,
  hour: "services_per_hour",
  item: "services_per_item",
  person: "services_per_person",
  day: "services_per_day",
};

interface Props {
  vendorId: string;
  selectedId?: string | null;
  onSelect?: (s: Service | null) => void;
  showFallbackEmpty?: boolean;
}

/** Customer-facing list of vendor services with optional selection. */
export function ServicesList({ vendorId, selectedId, onSelect, showFallbackEmpty }: Props) {
  const { t } = useI18n();
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
        {t("services_loading")}
      </div>
    );
  }

  if (items.length === 0) {
    if (!showFallbackEmpty) return null;
    return (
      <div className="rounded-xl border border-dashed py-6 text-center text-sm text-[var(--color-muted-foreground)]">
        {t("services_empty_public")}
      </div>
    );
  }

  return (
    <ul className="space-y-2">
      {items.map((s) => {
        const active = s.id === selectedId;
        const Tag = onSelect ? "button" : "div";
        const suffixKey = UNIT_SUFFIX_KEY[s.unit];
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
                  {suffixKey && (
                    <span className="font-mono text-[10px] text-[var(--color-muted-foreground)]">
                      {t(suffixKey)}
                    </span>
                  )}
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
