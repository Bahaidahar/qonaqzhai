"use client";

import { useEffect, useState } from "react";
import { Bell } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { getToken } from "@/shared/api";
import { API_BASE } from "@/shared/config/env";
import { useI18n } from "@/shared/i18n/context";
import type { DictKey } from "@/shared/i18n/dict";

const NOTIF_TITLE_KEY: Record<string, DictKey> = {
  "signup.welcome": "notif_signup_welcome_title",
  "auth.password_reset": "notif_password_reset_title",
  "booking.created": "notif_booking_created_title",
  "booking.accepted": "notif_booking_accepted_title",
  "booking.declined": "notif_booking_declined_title",
  "booking.paid": "notif_booking_paid_title",
  "vendor.approved": "notif_vendor_approved_title",
  "vendor.rejected": "notif_vendor_rejected_title",
  "thread.message": "notif_thread_message_title",
};

const NOTIF_BODY_KEY: Record<string, DictKey> = {
  "signup.welcome": "notif_signup_welcome_body",
  "auth.password_reset": "notif_password_reset_body",
  "booking.created": "notif_booking_created_body",
  "booking.accepted": "notif_booking_accepted_body",
  "booking.declined": "notif_booking_declined_body",
  "booking.paid": "notif_booking_paid_body",
  "vendor.approved": "notif_vendor_approved_body",
  "vendor.rejected": "notif_vendor_rejected_body",
};

const CHANNEL_KEY: Record<string, DictKey> = {
  email: "notif_channel_email",
  push: "notif_channel_push",
};

const STATUS_KEY: Record<string, DictKey> = {
  pending: "notif_status_pending",
  sent: "notif_status_sent",
  failed: "notif_status_failed",
};

interface Notification {
  id: string;
  userId: string;
  type: string;
  channel: string;
  title: string;
  body: string;
  status: string;
  createdAt: string;
}

export default function NotificationsPage() {
  return (
    <AuthGate>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <Inbox />
        </main>
      </div>
    </AuthGate>
  );
}

function Inbox() {
  const { t, locale } = useI18n();
  const tr = (key: DictKey | undefined, fallback: string) =>
    key ? t(key) : fallback;
  const [items, setItems] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API_BASE}/api/notifications`, {
      headers: { Authorization: `Bearer ${getToken() ?? ""}` },
    })
      .then((r) => r.json())
      .then((d: { items: Notification[] }) => setItems(d.items ?? []))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="mx-auto max-w-3xl px-6 py-10">
      <h1 className="font-display text-4xl tracking-[-0.045em]">
        {t("notifications_title")}
      </h1>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        {t("notifications_hint")}
      </p>

      {loading ? (
        <div className="mt-10 text-sm text-[var(--color-muted-foreground)]">
          {t("common_loading")}
        </div>
      ) : items.length === 0 ? (
        <div className="mt-10 rounded-xl border border-dashed py-16 text-center text-sm text-[var(--color-muted-foreground)]">
          <Bell className="mx-auto h-6 w-6 opacity-40" />
          <p className="mt-3">{t("notifications_empty")}</p>
        </div>
      ) : (
        <ul className="mt-6 space-y-3">
          {items.map((n) => {
            const title = tr(NOTIF_TITLE_KEY[n.type], n.title);
            const bodyKey = NOTIF_BODY_KEY[n.type];
            const body = bodyKey ? t(bodyKey) : n.body;
            const isHtml = !bodyKey && /<[^>]+>/.test(n.body);
            const channel = tr(CHANNEL_KEY[n.channel], n.channel);
            const status = tr(STATUS_KEY[n.status], n.status);
            const localeTag =
              locale === "ru" ? "ru-RU" : locale === "kz" ? "kk-KZ" : "en-US";
            return (
              <li
                key={n.id}
                className="rounded-xl border bg-[var(--color-card)] p-4"
              >
                <div className="flex items-center justify-between">
                  <div className="font-medium">{title}</div>
                  <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
                    {n.type}
                  </span>
                </div>
                {isHtml ? (
                  <p
                    className="mt-1 text-sm text-[var(--color-muted-foreground)]"
                    dangerouslySetInnerHTML={{ __html: body }}
                  />
                ) : (
                  <p className="mt-1 text-sm text-[var(--color-muted-foreground)]">
                    {body}
                  </p>
                )}
                <div className="mt-2 flex items-center gap-3 text-[10px] text-[var(--color-muted-foreground)]">
                  <span>{new Date(n.createdAt).toLocaleString(localeTag)}</span>
                  <span className="rounded-full border px-1.5 py-0.5">{channel}</span>
                  <span className="rounded-full border px-1.5 py-0.5">{status}</span>
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
