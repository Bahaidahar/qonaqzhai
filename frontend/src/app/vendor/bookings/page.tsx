"use client";

import { useCallback, useEffect, useState } from "react";
import { Check, X, Calendar, Users, MessageSquare } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { Button } from "@/shared/ui/button";
import { useI18n } from "@/shared/i18n/context";
import { useLabels } from "@/shared/i18n/labels";
import { api, type Booking } from "@/shared/api";

export default function VendorBookingsPage() {
  return (
    <AuthGate allowedRoles={["vendor"]}>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <BookingsInbox />
        </main>
      </div>
    </AuthGate>
  );
}

function BookingsInbox() {
  const { t } = useI18n();
  const labels = useLabels();
  const [items, setItems] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    const res = await api.bookings();
    setItems(res.items ?? []);
    setLoading(false);
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  async function update(id: string, status: "accepted" | "declined") {
    await api.updateBooking(id, status);
    await load();
  }

  return (
    <div className="mx-auto max-w-4xl px-6 py-10">
      <h1 className="font-display text-4xl tracking-[-0.045em]">
        {t("vendor_bookings_title")}
      </h1>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        {t("vendor_bookings_hint")}
      </p>

      {loading ? (
        <div className="mt-10 text-sm text-[var(--color-muted-foreground)]">
          {t("common_loading")}
        </div>
      ) : items.length === 0 ? (
        <div className="mt-8 rounded-xl border border-dashed py-16 text-center text-sm text-[var(--color-muted-foreground)]">
          {t("vendor_bookings_empty")}
        </div>
      ) : (
        <ul className="mt-8 space-y-3">
          {items.map((b) => (
            <li
              key={b.id}
              className="rounded-xl border bg-[var(--color-card)] p-5"
            >
              <div className="flex flex-wrap items-start justify-between gap-4">
                <div className="space-y-1">
                  <div className="flex flex-wrap items-center gap-4 text-sm text-[var(--color-muted-foreground)]">
                    <span className="inline-flex items-center gap-1.5">
                      <Calendar className="h-3.5 w-3.5" />
                      {b.eventDate}
                    </span>
                    <span className="inline-flex items-center gap-1.5">
                      <Users className="h-3.5 w-3.5" />
                      {b.guestCount} {t("vendor_bookings_guests")}
                    </span>
                  </div>
                  {b.note && (
                    <p className="flex items-start gap-1.5 pt-2 text-sm">
                      <MessageSquare className="mt-0.5 h-3.5 w-3.5 shrink-0 text-[var(--color-muted-foreground)]" />
                      {b.note}
                    </p>
                  )}
                  <div className="pt-2 font-mono text-[10px] uppercase tracking-wide text-[var(--color-muted-foreground)]">
                    {labels.bookingStatus(b.status)}
                  </div>
                </div>

                {b.status === "pending" && (
                  <div className="flex gap-2">
                    <Button
                      onClick={() => update(b.id, "declined")}
                      variant="outline"
                      size="sm"
                    >
                      <X className="h-4 w-4" />
                      {t("vendor_bookings_btn_decline")}
                    </Button>
                    <Button onClick={() => update(b.id, "accepted")} size="sm">
                      <Check className="h-4 w-4" />
                      {t("vendor_bookings_btn_accept")}
                    </Button>
                  </div>
                )}
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
