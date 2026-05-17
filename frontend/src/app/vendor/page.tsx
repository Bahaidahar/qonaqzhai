"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import Image from "next/image";
import { Upload, Trash2, Save, CheckCircle2, AlertCircle } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { RedirectIfWrongRole } from "@/features/auth/role-redirect";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Select } from "@/shared/ui/select";
import { useI18n } from "@/shared/i18n/context";
import { useLabels } from "@/shared/i18n/labels";
import { api, ApiError, photoURL, type Vendor } from "@/shared/api";
import { formatKZT } from "@/shared/lib/utils";
import type { DictKey } from "@/shared/i18n/dict";
import { ServicesManager } from "@/features/services/services-manager";

const CATEGORIES = [
  "Venue",
  "Catering",
  "Music & DJ",
  "Photo & Video",
  "Decor & Florists",
  "Cakes",
  "Other",
];

// MVP locked to Almaty; multi-city support will be re-enabled at expansion time.
const FIXED_CITY = "Almaty";

export default function VendorDashboardPage() {
  return (
    <AuthGate allowedRoles={["vendor"]}>
      <RedirectIfWrongRole expected="vendor" />
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <VendorEditor />
        </main>
      </div>
    </AuthGate>
  );
}

function VendorEditor() {
  const { t } = useI18n();
  const labels = useLabels();
  const [vendor, setVendor] = useState<Vendor | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [name, setName] = useState("");
  const [category, setCategory] = useState(CATEGORIES[0]);
  // City is locked to Almaty (MVP). Kept as state so the upsert payload stays consistent.
  const [city] = useState(FIXED_CITY);
  const [description, setDescription] = useState("");
  const [priceFrom, setPriceFrom] = useState("0");

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const v = await api.vendorMine();
      setVendor(v);
      setName(v.name);
      setCategory(v.category);
      // City stays Almaty regardless of stored value.
      setDescription(v.description);
      setPriceFrom(String(v.priceFrom));
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        setVendor(null);
      }
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  async function save() {
    setSaving(true);
    setError(null);
    try {
      const v = await api.vendorUpsert({
        name,
        category,
        city,
        description,
        priceFrom: Number(priceFrom) || 0,
      });
      setVendor(v);
    } catch (err) {
      if (err instanceof ApiError) setError(err.message);
    } finally {
      setSaving(false);
    }
  }

  if (loading) {
    return (
      <div className="p-10 text-sm text-[var(--color-muted-foreground)]">
        {t("common_loading")}
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl px-6 py-10">
      <h1 className="font-display text-4xl tracking-[-0.045em]">
        {t("vendor_profile_title")}
      </h1>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        {t("vendor_profile_hint")}
      </p>

      {vendor && (
        <div className="mt-6">
          <StatusBadge status={vendor.status} />
        </div>
      )}

      <section className="mt-8 rounded-xl border bg-[var(--color-card)] p-6">
        <h2 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
          / {t("vendor_profile_section_basics")}
        </h2>
        <div className="mt-4 grid gap-4 sm:grid-cols-2">
          <Field label={t("vendor_profile_field_name")}>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t("vendor_profile_field_name_ph")}
            />
          </Field>
          <Field label={t("vendor_profile_field_category")}>
            <Select
              value={category}
              onChange={setCategory}
              options={CATEGORIES.map((v) => ({ value: v, label: labels.category(v) }))}
            />
          </Field>
          <Field label={t("vendor_profile_field_city")}>
            <div className="flex h-11 w-full items-center rounded-xl border bg-[var(--color-muted)]/40 px-4 text-sm text-[var(--color-muted-foreground)]">
              {labels.city(FIXED_CITY)}
            </div>
          </Field>
          <Field label={t("vendor_profile_field_price")}>
            <Input
              type="number"
              value={priceFrom}
              onChange={(e) => setPriceFrom(e.target.value)}
              placeholder="500000"
            />
          </Field>
          <Field label={t("vendor_profile_field_desc")} full>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder={t("vendor_profile_field_desc_ph")}
              rows={4}
              className="w-full resize-none rounded-xl border bg-[var(--color-input)]/60 px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-ring)]"
            />
          </Field>
        </div>

        {error && (
          <div className="mt-4 rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
            {error}
          </div>
        )}

        <div className="mt-6 flex items-center justify-between">
          <div className="text-xs text-[var(--color-muted-foreground)]">
            {priceFrom
              ? `${t("vendor_profile_listed_from")} ${formatKZT(Number(priceFrom))}`
              : ""}
          </div>
          <Button
            onClick={save}
            disabled={saving || !name || !category || !city}
          >
            <Save className="h-4 w-4" />
            {saving ? "..." : t("common_save")}
          </Button>
        </div>
      </section>

      {vendor && <PhotosManager vendor={vendor} onReload={load} />}
      {vendor && <ServicesManager vendorId={vendor.id} />}
    </div>
  );
}

function Field({
  label,
  children,
  full,
}: {
  label: string;
  children: React.ReactNode;
  full?: boolean;
}) {
  return (
    <div className={full ? "sm:col-span-2 space-y-1.5" : "space-y-1.5"}>
      <label className="text-sm font-medium">{label}</label>
      {children}
    </div>
  );
}


const STATUS_KEYS: Record<Vendor["status"], DictKey> = {
  pending: "vendor_profile_status_pending",
  approved: "vendor_profile_status_approved",
  rejected: "vendor_profile_status_rejected",
};

function StatusBadge({ status }: { status: Vendor["status"] }) {
  const { t } = useI18n();
  const map = {
    pending: {
      icon: AlertCircle,
      cls: "border-amber-500/30 bg-amber-500/10 text-amber-700 dark:text-amber-400",
    },
    approved: {
      icon: CheckCircle2,
      cls: "border-emerald-500/30 bg-emerald-500/10 text-emerald-700 dark:text-emerald-400",
    },
    rejected: {
      icon: AlertCircle,
      cls: "border-red-500/30 bg-red-500/10 text-red-700 dark:text-red-400",
    },
  } as const;
  const m = map[status];
  return (
    <div className={`flex items-start gap-2 rounded-lg border px-3 py-2 text-xs ${m.cls}`}>
      <m.icon className="mt-0.5 h-3.5 w-3.5 shrink-0" />
      <span>{t(STATUS_KEYS[status])}</span>
    </div>
  );
}

function PhotosManager({
  vendor,
  onReload,
}: {
  vendor: Vendor;
  onReload: () => void;
}) {
  const { t } = useI18n();
  const fileRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onPick(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setUploading(true);
    setError(null);
    try {
      await api.uploadPhoto(file);
      await onReload();
    } catch (err) {
      if (err instanceof ApiError) setError(err.message);
    } finally {
      setUploading(false);
      if (fileRef.current) fileRef.current.value = "";
    }
  }

  async function onDelete(id: string) {
    await api.deletePhoto(id);
    await onReload();
  }

  return (
    <section className="mt-6 rounded-xl border bg-[var(--color-card)] p-6">
      <div className="flex items-start justify-between">
        <div>
          <h2 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
            / {t("vendor_profile_section_photos")}
          </h2>
          <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
            {t("vendor_profile_photos_hint")}
          </p>
        </div>
        <Button
          onClick={() => fileRef.current?.click()}
          disabled={uploading}
          variant="outline"
          size="sm"
        >
          <Upload className="h-4 w-4" />
          {uploading ? "..." : t("vendor_profile_btn_upload")}
        </Button>
        <input
          ref={fileRef}
          type="file"
          accept="image/*"
          className="hidden"
          onChange={onPick}
        />
      </div>

      {error && (
        <div className="mt-3 rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
          {error}
        </div>
      )}

      {vendor.photoIds.length === 0 ? (
        <div className="mt-6 rounded-lg border border-dashed py-10 text-center text-sm text-[var(--color-muted-foreground)]">
          {t("vendor_profile_no_photos")}
        </div>
      ) : (
        <div className="mt-6 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
          {vendor.photoIds.map((id) => (
            <div
              key={id}
              className="group relative aspect-square overflow-hidden rounded-lg border"
            >
              <Image
                src={photoURL(id)}
                alt=""
                fill
                sizes="(max-width: 768px) 50vw, 25vw"
                className="object-cover"
                unoptimized
              />
              <button
                onClick={() => onDelete(id)}
                className="absolute right-2 top-2 grid h-7 w-7 place-items-center rounded-md bg-[var(--color-card)]/80 text-[var(--color-destructive)] opacity-0 backdrop-blur transition group-hover:opacity-100"
                aria-label={t("vendor_profile_btn_delete")}
              >
                <Trash2 className="h-3.5 w-3.5" />
              </button>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}
