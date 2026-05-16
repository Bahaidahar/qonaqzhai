"use client";

import Link from "next/link";
import { useRouter, useSearchParams, usePathname } from "next/navigation";
import {
  Bell,
  MessageSquare,
  PanelLeftClose,
  PanelLeftOpen,
  Plus,
  Store,
  ClipboardList,
  User,
  Users,
  ShieldCheck,
  Trash2,
  type LucideIcon,
} from "lucide-react";
import { BrandRow } from "@/widgets/header/brand";
import { cn } from "@/shared/lib/utils";
import { useI18n } from "@/shared/i18n/context";
import { useSidebar } from "./sidebar-context";
import { useAuth } from "@/features/auth/context";
import type { DictKey } from "@/shared/i18n/dict";
import { deleteChat, useChatHistory } from "@/features/ai-chat/history";

interface NavItem {
  href: string;
  icon: LucideIcon;
  labelKey: DictKey;
}

const NAV_BY_ROLE: Record<"customer" | "vendor" | "admin", NavItem[]> = {
  customer: [
    { href: "/", icon: MessageSquare, labelKey: "nav_chat" },
    { href: "/vendors", icon: Store, labelKey: "nav_vendors" },
    { href: "/bookings", icon: ClipboardList, labelKey: "nav_bookings" },
    { href: "/notifications", icon: Bell, labelKey: "nav_notifications" },
  ],
  vendor: [
    { href: "/vendor", icon: User, labelKey: "nav_vendor_profile" },
    { href: "/vendor/bookings", icon: ClipboardList, labelKey: "nav_bookings" },
    { href: "/notifications", icon: Bell, labelKey: "nav_notifications" },
  ],
  admin: [
    { href: "/admin", icon: ShieldCheck, labelKey: "nav_admin_vendors" },
    { href: "/admin/users", icon: Users, labelKey: "nav_admin_users" },
  ],
};

function initials(name: string): string {
  return name
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((s) => s[0]?.toUpperCase() ?? "")
    .join("");
}

export function ChatSidebar() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const router = useRouter();
  const { t } = useI18n();
  const { open, toggle } = useSidebar();
  const { user } = useAuth();
  const { items: chats, refresh } = useChatHistory();
  const currentChatId = searchParams.get("c");

  if (!open) {
    return (
      <button
        onClick={toggle}
        className="fixed left-3 top-3 z-30 grid h-9 w-9 place-items-center rounded-lg border bg-[var(--color-card)] text-[var(--color-muted-foreground)] shadow-sm transition hover:text-[var(--color-foreground)]"
        aria-label="Open sidebar"
      >
        <PanelLeftOpen className="h-4 w-4" />
      </button>
    );
  }

  const nav = user ? NAV_BY_ROLE[user.role] : [];
  const showHistory = user?.role === "customer";

  function startNewChat() {
    // Drop the chat query param — first message creates a chat on the server.
    router.push("/");
  }

  async function removeChat(id: string) {
    try {
      await deleteChat(id);
    } catch {
      // ignore; refresh below will reconcile
    }
    await refresh();
    if (id === currentChatId) {
      router.push("/");
    }
  }

  return (
    <aside className="hidden h-screen w-72 shrink-0 flex-col border-r bg-[var(--color-card)] md:flex">
      <div className="flex h-16 items-center justify-between border-b px-5">
        <Link href="/">
          <BrandRow />
        </Link>
        <button
          onClick={toggle}
          className="grid h-7 w-7 place-items-center rounded-md text-[var(--color-muted-foreground)] transition hover:bg-[var(--color-muted)] hover:text-[var(--color-foreground)]"
          aria-label="Collapse sidebar"
        >
          <PanelLeftClose className="h-4 w-4" />
        </button>
      </div>

      <nav className="flex-1 space-y-0.5 overflow-y-auto p-3">
        {showHistory && (
          <button
            onClick={startNewChat}
            className="mb-3 flex w-full items-center justify-center gap-2 rounded-lg border bg-[var(--color-card)] px-3 py-2 text-sm font-medium transition hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-muted)]"
          >
            <Plus className="h-4 w-4" />
            {t("sidebar_new_chat")}
          </button>
        )}

        {nav.map((item) => {
          const active =
            pathname === item.href ||
            (item.href !== "/" && pathname.startsWith(item.href));
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-all",
                active
                  ? "bg-[var(--color-primary)]/10 text-[var(--color-primary)]"
                  : "text-[var(--color-muted-foreground)] hover:translate-x-0.5 hover:bg-[var(--color-muted)] hover:text-[var(--color-foreground)]"
              )}
            >
              <item.icon className="h-4 w-4" />
              {t(item.labelKey)}
            </Link>
          );
        })}

        {showHistory && chats.length > 0 && (
          <div className="mt-5">
            <div className="px-3 pb-2 font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
              {t("sidebar_recent")}
            </div>
            <ul className="space-y-0.5">
              {chats.map((c) => {
                const active = c.id === currentChatId;
                return (
                  <li key={c.id} className="group relative">
                    <Link
                      href={`/?c=${c.id}`}
                      className={cn(
                        "block truncate rounded-lg px-3 py-2 pr-9 text-sm transition-colors",
                        active
                          ? "bg-[var(--color-primary)]/10 text-[var(--color-primary)]"
                          : "text-[var(--color-muted-foreground)] hover:bg-[var(--color-muted)] hover:text-[var(--color-foreground)]"
                      )}
                      title={c.title}
                    >
                      {c.title}
                    </Link>
                    <button
                      onClick={(e) => {
                        e.preventDefault();
                        removeChat(c.id);
                      }}
                      className="absolute right-2 top-1/2 -translate-y-1/2 grid h-6 w-6 -translate-y-1/2 place-items-center rounded-md text-[var(--color-muted-foreground)] opacity-0 transition group-hover:opacity-100 hover:bg-[var(--color-muted)] hover:text-[var(--color-destructive)]"
                      aria-label="Delete chat"
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </li>
                );
              })}
            </ul>
          </div>
        )}
      </nav>

      <div className="border-t p-2">
        <Link
          href="/settings"
          className={cn(
            "flex w-full items-center gap-3 rounded-lg p-2 transition-colors",
            pathname === "/settings"
              ? "bg-[var(--color-muted)]"
              : "hover:bg-[var(--color-muted)]"
          )}
        >
          <div className="grid h-8 w-8 place-items-center rounded-md bg-[var(--color-muted)] text-xs font-semibold text-[var(--color-foreground)]">
            {user ? initials(user.name) : "—"}
          </div>
          <div className="min-w-0 flex-1 text-left">
            <div className="truncate text-sm font-medium">
              {user?.name ?? "—"}
            </div>
            <div className="truncate text-[10px] uppercase tracking-wide text-[var(--color-muted-foreground)]">
              {user?.role ?? ""}
            </div>
          </div>
        </Link>
      </div>
    </aside>
  );
}
