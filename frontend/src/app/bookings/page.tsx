"use client";

import { useCallback, useEffect, useState } from "react";
import { Calendar, Users, X, Star } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { Button } from "@/shared/ui/button";
import { useI18n } from "@/shared/i18n/context";
import { api, type Booking } from "@/shared/api";
import type { DictKey } from "@/shared/i18n/dict";
import { ReviewForm } from "@/features/reviews/review-form";

const STATUS_COLOR: Record<Booking["status"], string> = {
  pending: "border-amber-500/30 bg-amber-500/10 text-amber-700",
  accepted: "border-emerald-500/30 bg-emerald-500/10 text-emerald-700",
  paid: "border-emerald-500/30 bg-emerald-500/10 text-emerald-700",
  completed: "border-sky-500/30 bg-sky-500/10 text-sky-700",
  declined: "border-red-500/30 bg-red-500/10 text-red-700",
  cancelled: "border-slate-500/30 bg-slate-500/10 text-slate-600",
};

const STATUS_KEY: Record<Booking["status"], DictKey> = {
  pending: "bookings_status_pending",
  accepted: "bookings_status_accepted",
  paid: "bookings_status_accepted",
  completed: "bookings_status_accepted",
  declined: "bookings_status_declined",
  cancelled: "bookings_status_cancelled",
};

export default function BookingsPage() {
  return (
    <AuthGate allowedRoles={["customer"]}>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <BookingsList />
        </main>
      </div>
    </AuthGate>
  );
}

function BookingsList() {
  const { t } = useI18n();
  const [items, setItems] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [reviewingId, setReviewingId] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    const res = await api.bookings();
    setItems(res.items ?? []);
    setLoading(false);
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  async function cancel(id: string) {
    await api.updateBooking(id, "cancelled");
    await load();
  }

  return (
    <div className="mx-auto max-w-4xl px-6 py-10">
      <h1 className="font-display text-4xl tracking-[-0.045em]">
        {t("bookings_title")}
      </h1>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        {t("bookings_hint")}
      </p>

      {loading ? (
        <div className="mt-10 text-sm text-[var(--color-muted-foreground)]">
          {t("common_loading")}
        </div>
      ) : items.length === 0 ? (
        <div className="mt-8 rounded-xl border border-dashed py-16 text-center text-sm text-[var(--color-muted-foreground)]">
          {t("bookings_empty")}
        </div>
      ) : (
        <ul className="mt-8 space-y-3">
          {items.map((b) => (
            <li
              key={b.id}
              className="rounded-xl border bg-[var(--color-card)] p-5"
            >
              <div className="flex flex-wrap items-start justify-between gap-4">
                <div className="space-y-2">
                  <div className="flex flex-wrap items-center gap-4 text-sm text-[var(--color-muted-foreground)]">
                    <span className="inline-flex items-center gap-1.5">
                      <Calendar className="h-3.5 w-3.5" />
                      {b.eventDate}
                    </span>
                    <span className="inline-flex items-center gap-1.5">
                      <Users className="h-3.5 w-3.5" />
                      {b.guestCount}
                    </span>
                  </div>
                  {b.note && (
                    <p className="text-sm text-[var(--color-foreground)]">{b.note}</p>
                  )}
                  <span
                    className={`inline-block rounded-full border px-2 py-0.5 text-[10px] font-medium uppercase tracking-wide ${STATUS_COLOR[b.status]}`}
                  >
                    {t(STATUS_KEY[b.status])}
                  </span>
                </div>
                <div className="flex flex-col gap-2">
                  {b.status === "pending" && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => cancel(b.id)}
                    >
                      <X className="h-4 w-4" />
                      {t("bookings_btn_cancel")}
                    </Button>
                  )}
                  {b.status === "completed" && (
                    <Button
                      size="sm"
                      onClick={() =>
                        setReviewingId(reviewingId === b.id ? null : b.id)
                      }
                    >
                      <Star className="h-4 w-4" />
                      {reviewingId === b.id ? t("common_cancel") : t("reviews_title")}
                    </Button>
                  )}
                </div>
              </div>
              {reviewingId === b.id && (
                <div className="mt-4">
                  <ReviewForm
                    bookingId={b.id}
                    onSubmitted={() => setReviewingId(null)}
                  />
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
