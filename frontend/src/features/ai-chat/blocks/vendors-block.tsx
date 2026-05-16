"use client";

import { Store, Star, MapPin, ArrowRight } from "lucide-react";
import { Card, CardContent } from "@/shared/ui/card";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { formatKZT } from "@/shared/lib/utils";
import { useI18n } from "@/shared/i18n/context";
import type { VendorsBlock as VendorsBlockData } from "@/features/ai-chat/types";

export function VendorsBlock({ data }: { data: VendorsBlockData }) {
  const { t } = useI18n();
  return (
    <Card className="hover-lift overflow-hidden">
      <div className="flex items-center gap-2 border-b px-4 py-2.5">
        <Store className="h-3.5 w-3.5 text-[var(--color-primary)]" />
        <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
          {t("block_vendors")} · {data.query}
        </span>
        <span className="ml-auto font-mono text-xs font-semibold">
          {data.items.length} {t("vendor_matches")}
        </span>
      </div>
      <CardContent className="p-5">
        <div className="-mx-1 flex gap-3 overflow-x-auto px-1 pb-1">
          {data.items.map((v) => (
            <div
              key={v.id}
              className="group w-64 shrink-0 cursor-pointer rounded-xl border bg-[var(--color-card)] transition hover:-translate-y-0.5 hover:border-[var(--color-primary)]/50 hover:shadow-md"
            >
              <div className="aspect-[5/3] overflow-hidden rounded-t-xl">
                <div className="img-zoom h-full w-full bg-gradient-to-br from-[var(--color-primary)]/25 via-[var(--color-muted)] to-[var(--color-accent)]/40" />
              </div>
              <div className="space-y-2 p-3">
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <Badge variant="outline" className="mb-1">
                      {v.category}
                    </Badge>
                    <div className="text-sm font-semibold leading-tight">
                      {v.name}
                    </div>
                  </div>
                  <span className="flex shrink-0 items-center gap-0.5 text-xs">
                    <Star className="h-3 w-3 fill-amber-400 text-amber-400" />
                    <span className="font-medium">
                      {(v.rating ?? 0).toFixed(1)}
                    </span>
                  </span>
                </div>
                <div className="flex items-center gap-1 text-[10px] text-[var(--color-muted-foreground)]">
                  <MapPin className="h-2.5 w-2.5" />
                  {v.city}
                </div>
                {(v.tags ?? []).length > 0 && (
                  <div className="flex flex-wrap gap-1">
                    {(v.tags ?? []).slice(0, 2).map((t) => (
                      <span
                        key={t}
                        className="rounded-full bg-[var(--color-muted)] px-1.5 py-0.5 text-[9px] text-[var(--color-muted-foreground)]"
                      >
                        {t}
                      </span>
                    ))}
                  </div>
                )}
                <div className="flex items-end justify-between border-t pt-2">
                  <div className="text-[10px] text-[var(--color-muted-foreground)]">
                    from{" "}
                    <span className="font-medium text-[var(--color-foreground)]">
                      {formatKZT(v.priceFrom)}
                    </span>
                  </div>
                  <ArrowRight className="arrow-slide h-3.5 w-3.5 text-[var(--color-muted-foreground)]" />
                </div>
              </div>
            </div>
          ))}
        </div>
        <div className="mt-4 flex justify-end">
          <Button variant="ghost" size="sm">
            {t("vendor_view_all")}
            <ArrowRight className="h-3.5 w-3.5" />
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
