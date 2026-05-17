"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { MessageSquare } from "lucide-react";
import { api, type BookingThread } from "@/shared/api";

export function ThreadList() {
  const [items, setItems] = useState<BookingThread[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .listThreads()
      .then((r) => setItems(r.items ?? []))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="text-sm text-[var(--color-muted-foreground)]">Loading…</div>;
  }
  if (items.length === 0) {
    return (
      <div className="rounded-xl border border-dashed py-10 text-center text-sm text-[var(--color-muted-foreground)]">
        <MessageSquare className="mx-auto h-6 w-6 opacity-40" />
        <p className="mt-3">
          No active conversations. Threads open automatically when a vendor accepts a booking.
        </p>
      </div>
    );
  }

  return (
    <ul className="space-y-2">
      {items.map((t) => (
        <li key={t.id}>
          <Link
            href={`/threads/${t.id}`}
            className="block rounded-xl border bg-[var(--color-card)] p-4 transition hover:border-[var(--color-primary)]/40"
          >
            <div className="flex items-center justify-between">
              <span className="font-medium">Booking {t.bookingId.slice(0, 8)}</span>
              <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
                {new Date(t.updatedAt).toLocaleString()}
              </span>
            </div>
          </Link>
        </li>
      ))}
    </ul>
  );
}
