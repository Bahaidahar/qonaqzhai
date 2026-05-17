"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Calendar, MessageSquare, Users, User } from "lucide-react";
import { api, type ThreadSummary } from "@/shared/api";
import { useI18n } from "@/shared/i18n/context";
import { formatKZT } from "@/shared/lib/utils";

export function ThreadList() {
  const { t, locale } = useI18n();
  const [items, setItems] = useState<ThreadSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .listThreads()
      .then((r) => setItems(r.items ?? []))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="text-sm text-[var(--color-muted-foreground)]">
        {t("common_loading")}
      </div>
    );
  }
  if (items.length === 0) {
    return (
      <div className="rounded-xl border border-dashed py-10 text-center text-sm text-[var(--color-muted-foreground)]">
        <MessageSquare className="mx-auto h-6 w-6 opacity-40" />
        <p className="mt-3">{t("threads_empty")}</p>
      </div>
    );
  }

  return (
    <ul className="space-y-2">
      {items.map((s) => (
        <li key={s.thread.id}>
          <Link
            href={`/threads/${s.thread.id}`}
            className="block rounded-xl border bg-[var(--color-card)] p-4 transition hover:border-[var(--color-primary)]/40"
          >
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0">
                <div className="truncate font-semibold">
                  {s.vendorName || t("threads_unknown_vendor")}
                </div>
                <div className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-[var(--color-muted-foreground)]">
                  <span className="inline-flex items-center gap-1">
                    <User className="h-3 w-3" />
                    {s.counterpart || t("threads_unknown_user")}
                  </span>
                  {s.eventDate && (
                    <span className="inline-flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      {formatDate(s.eventDate, locale)}
                    </span>
                  )}
                  {s.guestCount > 0 && (
                    <span className="inline-flex items-center gap-1">
                      <Users className="h-3 w-3" />
                      {s.guestCount}
                    </span>
                  )}
                  {s.amount > 0 && (
                    <span className="font-medium text-[var(--color-foreground)]">
                      {formatKZT(s.amount)}
                    </span>
                  )}
                  {s.status && (
                    <span className="rounded-full border px-2 py-0.5 text-[10px] uppercase tracking-widest">
                      {t((`booking_status_${s.status}` as never)) || s.status}
                    </span>
                  )}
                </div>
              </div>
              <span className="shrink-0 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
                {new Date(s.thread.updatedAt).toLocaleString(localeTag(locale))}
              </span>
            </div>
          </Link>
        </li>
      ))}
    </ul>
  );
}

function localeTag(locale: string): string {
  if (locale === "ru") return "ru-RU";
  if (locale === "kz") return "kk-KZ";
  return "en-US";
}

function formatDate(iso: string, locale: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleDateString(localeTag(locale), {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}
