"use client";

import { useCallback, useEffect, useState } from "react";
import Image from "next/image";
import { Check, ChevronUp, Eye, MapPin, X } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { RedirectIfWrongRole } from "@/features/auth/role-redirect";
import { Button } from "@/shared/ui/button";
import { useI18n } from "@/shared/i18n/context";
import { useLabels } from "@/shared/i18n/labels";
import { api, photoURL, type Service, type Vendor } from "@/shared/api";
import { formatKZT } from "@/shared/lib/utils";
import { AdminCharts } from "@/widgets/admin-dashboard/charts";

export default function AdminPage() {
  return (
    <AuthGate allowedRoles={["admin"]}>
      <RedirectIfWrongRole expected="admin" />
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <AdminDashboard />
        </main>
      </div>
    </AuthGate>
  );
}

function AdminDashboard() {
  const { t } = useI18n();
  const [vendors, setVendors] = useState<Vendor[]>([]);
  const [stats, setStats] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    const [vs, st] = await Promise.all([
      api.adminVendors(),
      api.adminStats(),
    ]);
    setVendors(vs.items ?? []);
    setStats(st);
    setLoading(false);
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  async function update(id: string, status: Vendor["status"]) {
    await api.adminUpdateVendor(id, status);
    await load();
  }

  const pending = vendors.filter((v) => v.status === "pending");
  const approved = vendors.filter((v) => v.status === "approved");
  const rejected = vendors.filter((v) => v.status === "rejected");

  return (
    <div className="mx-auto max-w-5xl px-6 py-10">
      <h1 className="font-display text-4xl tracking-[-0.045em]">
        {t("admin_title")}
      </h1>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        {t("admin_hint")}
      </p>

      <section className="mt-6 grid gap-3 sm:grid-cols-4">
        <Stat label={t("admin_stat_users")} value={stats.users ?? 0} />
        <Stat label={t("admin_stat_pending")} value={stats.vendors_pending ?? 0} />
        <Stat label={t("admin_stat_approved")} value={stats.vendors_approved ?? 0} />
        <Stat label={t("admin_stat_bookings")} value={stats.bookings_total ?? 0} />
      </section>

      <AdminCharts />

      {loading ? (
        <div className="mt-10 text-sm text-[var(--color-muted-foreground)]">
          {t("common_loading")}
        </div>
      ) : (
        <>
          <Group title={`${t("admin_group_pending")} (${pending.length})`}>
            {pending.length === 0 ? (
              <Empty />
            ) : (
              pending.map((v) => (
                <VendorRow
                  key={v.id}
                  vendor={v}
                  actions={
                    <>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => update(v.id, "rejected")}
                      >
                        <X className="h-4 w-4" />
                        {t("admin_btn_reject")}
                      </Button>
                      <Button size="sm" onClick={() => update(v.id, "approved")}>
                        <Check className="h-4 w-4" />
                        {t("admin_btn_approve")}
                      </Button>
                    </>
                  }
                />
              ))
            )}
          </Group>
          <Group title={`${t("admin_group_approved")} (${approved.length})`}>
            {approved.length === 0 ? (
              <Empty />
            ) : (
              approved.map((v) => (
                <VendorRow
                  key={v.id}
                  vendor={v}
                  actions={
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => update(v.id, "rejected")}
                    >
                      <X className="h-4 w-4" />
                      {t("admin_btn_suspend")}
                    </Button>
                  }
                />
              ))
            )}
          </Group>
          {rejected.length > 0 && (
            <Group title={`${t("admin_group_rejected")} (${rejected.length})`}>
              {rejected.map((v) => (
                <VendorRow
                  key={v.id}
                  vendor={v}
                  actions={
                    <Button size="sm" onClick={() => update(v.id, "approved")}>
                      <Check className="h-4 w-4" />
                      {t("admin_btn_reapprove")}
                    </Button>
                  }
                />
              ))}
            </Group>
          )}
        </>
      )}
    </div>
  );
}

function Stat({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border bg-[var(--color-card)] p-4">
      <div className="text-xs text-[var(--color-muted-foreground)]">{label}</div>
      <div className="mt-1 font-display text-2xl">{value}</div>
    </div>
  );
}

function Group({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section className="mt-8">
      <h2 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
        / {title}
      </h2>
      <ul className="mt-3 space-y-2">{children}</ul>
    </section>
  );
}

function Empty() {
  const { t } = useI18n();
  return (
    <li className="rounded-xl border border-dashed py-8 text-center text-xs text-[var(--color-muted-foreground)]">
      {t("admin_empty")}
    </li>
  );
}

function VendorRow({
  vendor: v,
  actions,
}: {
  vendor: Vendor;
  actions: React.ReactNode;
}) {
  const { t } = useI18n();
  const labels = useLabels();
  const cover = v.photoIds[0];
  const [open, setOpen] = useState(false);
  const [services, setServices] = useState<Service[] | null>(null);

  useEffect(() => {
    if (!open || services !== null) return;
    let cancelled = false;
    api
      .vendorServices(v.id)
      .then((r) => {
        if (!cancelled) setServices(r.items ?? []);
      })
      .catch(() => {
        if (!cancelled) setServices([]);
      });
    return () => {
      cancelled = true;
    };
  }, [open, services, v.id]);

  return (
    <li className="rounded-xl border bg-[var(--color-card)]">
      <div className="flex items-center gap-4 p-3">
        <div className="relative h-14 w-14 shrink-0 overflow-hidden rounded-lg">
          {cover ? (
            <Image
              src={photoURL(cover)}
              alt=""
              fill
              sizes="56px"
              className="object-cover"
              unoptimized
            />
          ) : (
            <div className="h-full w-full bg-[var(--color-muted)]" />
          )}
        </div>
        <div className="min-w-0 flex-1">
          <div className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
            {labels.category(v.category)}
          </div>
          <div className="truncate text-sm font-medium">{v.name}</div>
          <div className="mt-0.5 flex items-center gap-3 text-xs text-[var(--color-muted-foreground)]">
            <span className="inline-flex items-center gap-1">
              <MapPin className="h-3 w-3" />
              {labels.city(v.city)}
            </span>
            {v.priceFrom > 0 && (
              <span>
                {t("vendors_from")} {formatKZT(v.priceFrom)}
              </span>
            )}
          </div>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setOpen((o) => !o)}
            aria-expanded={open}
          >
            {open ? (
              <ChevronUp className="h-4 w-4" />
            ) : (
              <Eye className="h-4 w-4" />
            )}
            {open ? t("admin_btn_hide") : t("admin_btn_preview")}
          </Button>
          {actions}
        </div>
      </div>
      {open && (
        <div className="border-t px-3 pb-4 pt-3 sm:px-5">
          <VendorPreview vendor={v} services={services} />
        </div>
      )}
    </li>
  );
}

function VendorPreview({
  vendor: v,
  services,
}: {
  vendor: Vendor;
  services: Service[] | null;
}) {
  const { t } = useI18n();
  return (
    <div className="space-y-4">
      {v.photoIds.length === 0 ? (
        <p className="text-xs text-[var(--color-muted-foreground)]">
          {t("admin_preview_no_photos")}
        </p>
      ) : (
        <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
          {v.photoIds.map((pid) => (
            <div
              key={pid}
              className="relative aspect-[4/3] overflow-hidden rounded-lg"
            >
              <Image
                src={photoURL(pid)}
                alt=""
                fill
                sizes="(max-width: 640px) 50vw, 25vw"
                className="object-cover"
                unoptimized
              />
            </div>
          ))}
        </div>
      )}

      <div>
        <div className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
          {t("admin_preview_description")}
        </div>
        <p className="mt-1 whitespace-pre-line text-sm">
          {v.description?.trim()
            ? v.description
            : t("admin_preview_no_description")}
        </p>
      </div>

      <div>
        <div className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
          {t("admin_preview_services")}
        </div>
        {services === null ? (
          <p className="mt-1 text-xs text-[var(--color-muted-foreground)]">
            {t("common_loading")}
          </p>
        ) : services.length === 0 ? (
          <p className="mt-1 text-xs text-[var(--color-muted-foreground)]">
            {t("admin_preview_no_services")}
          </p>
        ) : (
          <ul className="mt-2 space-y-1 text-sm">
            {services.map((s) => (
              <li
                key={s.id}
                className="flex items-baseline justify-between gap-3 border-b border-dashed py-1 last:border-b-0"
              >
                <div className="min-w-0">
                  <div className="font-medium">{s.name}</div>
                  {s.description && (
                    <div className="text-xs text-[var(--color-muted-foreground)]">
                      {s.description}
                    </div>
                  )}
                </div>
                <div className="shrink-0 text-right">
                  <div className="font-medium">{formatKZT(s.price)}</div>
                  <div className="text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
                    {s.unit}
                  </div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

