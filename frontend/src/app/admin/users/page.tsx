"use client";

import { useCallback, useEffect, useState } from "react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { Button } from "@/shared/ui/button";
import { useI18n } from "@/shared/i18n/context";
import { useLabels } from "@/shared/i18n/labels";
import { api, type ApiUser } from "@/shared/api";

const ROLE_COLOR: Record<string, string> = {
  customer: "bg-blue-500/10 text-blue-700 border-blue-500/30",
  vendor: "bg-purple-500/10 text-purple-700 border-purple-500/30",
  admin: "bg-red-500/10 text-red-700 border-red-500/30",
};

export default function AdminUsersPage() {
  return (
    <AuthGate allowedRoles={["admin"]}>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <UsersList />
        </main>
      </div>
    </AuthGate>
  );
}

function UsersList() {
  const { t } = useI18n();
  const labels = useLabels();
  const [users, setUsers] = useState<ApiUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<"all" | "customer" | "vendor" | "admin">(
    "all"
  );

  const load = useCallback(async () => {
    setLoading(true);
    const res = await api.adminUsers();
    setUsers(res.items ?? []);
    setLoading(false);
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  async function toggle(u: ApiUser) {
    await api.adminUpdateUser(
      u.id,
      u.status === "active" ? "suspended" : "active"
    );
    await load();
  }

  const filtered = filter === "all" ? users : users.filter((u) => u.role === filter);

  return (
    <div className="mx-auto max-w-5xl px-6 py-10">
      <h1 className="font-display text-4xl tracking-[-0.045em]">
        {t("users_title")}
      </h1>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        {t("users_hint")}
      </p>

      <div className="mt-6 flex gap-1.5">
        {(["all", "customer", "vendor", "admin"] as const).map((r) => (
          <button
            key={r}
            onClick={() => setFilter(r)}
            className={
              filter === r
                ? "rounded-full border border-[var(--color-primary)] bg-[var(--color-primary)]/10 px-3 py-1.5 text-xs font-medium text-[var(--color-primary)]"
                : "chip-hover rounded-full border bg-[var(--color-card)] px-3 py-1.5 text-xs font-medium"
            }
          >
            {r === "all" ? t("users_filter_all") : labels.role(r)}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="mt-10 text-sm text-[var(--color-muted-foreground)]">
          {t("common_loading")}
        </div>
      ) : (
        <div className="mt-6 overflow-hidden rounded-xl border bg-[var(--color-card)]">
          <table className="w-full text-sm">
            <thead className="border-b bg-[var(--color-muted)]/50 text-xs uppercase tracking-wider text-[var(--color-muted-foreground)]">
              <tr>
                <th className="px-4 py-3 text-left font-medium">{t("users_col_name")}</th>
                <th className="px-4 py-3 text-left font-medium">{t("users_col_email")}</th>
                <th className="px-4 py-3 text-left font-medium">{t("users_col_role")}</th>
                <th className="px-4 py-3 text-left font-medium">{t("users_col_status")}</th>
                <th className="px-4 py-3 text-right font-medium">{t("users_col_actions")}</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((u) => (
                <tr key={u.id} className="border-b last:border-0">
                  <td className="px-4 py-3 font-medium">{u.name}</td>
                  <td className="px-4 py-3 text-[var(--color-muted-foreground)]">
                    {u.email}
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`inline-block rounded-full border px-2 py-0.5 text-[10px] font-medium uppercase tracking-wide ${ROLE_COLOR[u.role]}`}
                    >
                      {labels.role(u.role)}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-xs">{labels.userStatus(u.status)}</td>
                  <td className="px-4 py-3 text-right">
                    {u.role !== "admin" && (
                      <Button
                        size="sm"
                        variant={u.status === "active" ? "outline" : "primary"}
                        onClick={() => toggle(u)}
                      >
                        {u.status === "active"
                          ? t("users_btn_suspend")
                          : t("users_btn_activate")}
                      </Button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
