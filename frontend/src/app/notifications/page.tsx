"use client";

import { useEffect, useState } from "react";
import { Bell } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { getToken } from "@/shared/api";
import { API_BASE } from "@/shared/config/env";
import { useI18n } from "@/shared/i18n/context";

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
  const { t } = useI18n();
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
          {items.map((n) => (
            <li
              key={n.id}
              className="rounded-xl border bg-[var(--color-card)] p-4"
            >
              <div className="flex items-center justify-between">
                <div className="font-medium">{n.title}</div>
                <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
                  {n.type}
                </span>
              </div>
              <p
                className="mt-1 text-sm text-[var(--color-muted-foreground)]"
                dangerouslySetInnerHTML={{ __html: n.body }}
              />
              <div className="mt-2 flex items-center gap-3 text-[10px] text-[var(--color-muted-foreground)]">
                <span>{new Date(n.createdAt).toLocaleString()}</span>
                <span className="rounded-full border px-1.5 py-0.5">{n.channel}</span>
                <span className="rounded-full border px-1.5 py-0.5">{n.status}</span>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
