"use client";

import { useState } from "react";
import { Store, Star, MapPin, ArrowRight, Check, CalendarPlus } from "lucide-react";
import { Card, CardContent } from "@/shared/ui/card";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { formatKZT } from "@/shared/lib/utils";
import { useI18n } from "@/shared/i18n/context";
import { api, ApiError } from "@/shared/api";
import type { VendorsBlock as VendorsBlockData } from "@/features/ai-chat/types";

type VendorItem = VendorsBlockData["items"][number];

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
            <VendorCard
              key={v.id}
              vendor={v}
              defaultGuests={data.guests ?? 0}
              defaultDate={data.eventDate ?? ""}
            />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

type BookState = "idle" | "form" | "booking" | "done" | "error";

function VendorCard({
  vendor: v,
  defaultGuests,
  defaultDate,
}: {
  vendor: VendorItem;
  defaultGuests: number;
  defaultDate: string;
}) {
  const { t } = useI18n();
  const [state, setState] = useState<BookState>("idle");
  const [date, setDate] = useState(defaultDate);
  const [guests, setGuests] = useState(
    defaultGuests > 0 ? String(defaultGuests) : ""
  );
  const [error, setError] = useState<string | null>(null);

  async function book(eventDate: string, guestCount: number) {
    setState("booking");
    setError(null);
    try {
      await api.createBooking({
        vendorId: v.id,
        eventDate,
        guestCount,
        amount: v.priceFrom > 0 ? v.priceFrom : undefined,
      });
      setState("done");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : t("auth_network_error"));
      setState("error");
    }
  }

  function onBookClick() {
    // book straight away when the prompt already gave us a date
    if (defaultDate) {
      void book(defaultDate, defaultGuests > 0 ? defaultGuests : 1);
      return;
    }
    setState("form");
  }

  function onConfirm() {
    if (!date) {
      setError(t("vendor_detail_date"));
      return;
    }
    void book(date, Number(guests) || 1);
  }

  return (
    <div className="group w-64 shrink-0 rounded-xl border bg-[var(--color-card)] transition hover:-translate-y-0.5 hover:border-[var(--color-primary)]/50 hover:shadow-md">
      <div className="aspect-[5/3] overflow-hidden rounded-t-xl">
        <div className="img-zoom h-full w-full bg-gradient-to-br from-[var(--color-primary)]/25 via-[var(--color-muted)] to-[var(--color-accent)]/40" />
      </div>
      <div className="space-y-2 p-3">
        <div className="flex items-start justify-between gap-2">
          <div>
            <Badge variant="outline" className="mb-1">
              {v.category}
            </Badge>
            <div className="text-sm font-semibold leading-tight">{v.name}</div>
          </div>
          <span className="flex shrink-0 items-center gap-0.5 text-xs">
            <Star className="h-3 w-3 fill-amber-400 text-amber-400" />
            <span className="font-medium">{(v.rating ?? 0).toFixed(1)}</span>
          </span>
        </div>
        <div className="flex items-center gap-1 text-[10px] text-[var(--color-muted-foreground)]">
          <MapPin className="h-2.5 w-2.5" />
          {v.city}
        </div>
        <div className="flex items-end justify-between border-t pt-2">
          <div className="text-[10px] text-[var(--color-muted-foreground)]">
            from{" "}
            <span className="font-medium text-[var(--color-foreground)]">
              {formatKZT(v.priceFrom)}
            </span>
          </div>
          <ArrowRight className="arrow-slide h-3.5 w-3.5 text-[var(--color-muted-foreground)]" />
        </div>

        {/* booking control */}
        {state === "done" ? (
          <div className="flex items-center gap-1.5 rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-2 py-1.5 text-[11px] font-medium text-emerald-700">
            <Check className="h-3.5 w-3.5" />
            {t("vendor_detail_sent_title")}
          </div>
        ) : state === "form" ? (
          <div className="space-y-1.5 pt-1">
            <input
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              aria-label={t("vendor_detail_date")}
              className="w-full rounded-md border bg-transparent px-2 py-1 text-xs"
            />
            <input
              type="number"
              value={guests}
              onChange={(e) => setGuests(e.target.value)}
              placeholder={t("vendor_detail_guests")}
              aria-label={t("vendor_detail_guests")}
              className="w-full rounded-md border bg-transparent px-2 py-1 text-xs"
            />
            {error && <div className="text-[10px] text-red-600">{error}</div>}
            <Button size="sm" className="w-full" onClick={onConfirm}>
              {t("vendor_detail_book_now")}
            </Button>
          </div>
        ) : (
          <Button
            size="sm"
            className="mt-1 w-full"
            disabled={state === "booking"}
            onClick={onBookClick}
          >
            <CalendarPlus className="h-3.5 w-3.5" />
            {state === "booking" ? "…" : t("vendor_detail_book_now")}
          </Button>
        )}
        {state === "error" && error && (
          <div className="text-[10px] text-red-600">{error}</div>
        )}
      </div>
    </div>
  );
}
