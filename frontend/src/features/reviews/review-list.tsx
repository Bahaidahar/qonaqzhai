"use client";

import { useEffect, useState } from "react";
import { Star } from "lucide-react";
import { api, type Review } from "@/shared/api";
import { useI18n } from "@/shared/i18n/context";

export function ReviewList({ vendorId }: { vendorId: string }) {
  const { t } = useI18n();
  const [items, setItems] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .reviewsForVendor(vendorId)
      .then((r) => setItems(r.items ?? []))
      .finally(() => setLoading(false));
  }, [vendorId]);

  if (loading) {
    return (
      <div className="text-sm text-[var(--color-muted-foreground)]">
        {t("reviews_loading")}
      </div>
    );
  }
  if (items.length === 0) {
    return (
      <div className="rounded-xl border border-dashed py-10 text-center text-sm text-[var(--color-muted-foreground)]">
        {t("reviews_empty")}
      </div>
    );
  }
  return (
    <ul className="space-y-3">
      {items.map((r) => (
        <li
          key={r.id}
          className="rounded-xl border bg-[var(--color-card)] p-4"
        >
          <Stars rating={r.rating} />
          {r.text && <p className="mt-2 text-sm">{r.text}</p>}
          <div className="mt-2 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
            {new Date(r.createdAt).toLocaleDateString()}
          </div>
        </li>
      ))}
    </ul>
  );
}

export function Stars({ rating }: { rating: number }) {
  return (
    <div className="flex items-center gap-0.5">
      {[1, 2, 3, 4, 5].map((n) => (
        <Star
          key={n}
          className={
            n <= rating
              ? "h-3.5 w-3.5 fill-amber-400 text-amber-400"
              : "h-3.5 w-3.5 text-[var(--color-muted-foreground)]"
          }
        />
      ))}
    </div>
  );
}
