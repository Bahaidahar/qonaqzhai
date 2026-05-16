"use client";

import { useCallback, useEffect, useState } from "react";
import { Plus, Pencil, Trash2, Power, PowerOff } from "lucide-react";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Select } from "@/shared/ui/select";
import { api, ApiError, type Service, type ServiceUnit } from "@/shared/api";
import { SERVICE_UNITS } from "@/entities/service/types";
import { formatKZT } from "@/shared/lib/utils";
import { useI18n } from "@/shared/i18n/context";
import type { DictKey } from "@/shared/i18n/dict";

interface Props {
  vendorId: string;
}

const UNIT_SUFFIX_KEY: Record<ServiceUnit, DictKey | null> = {
  fixed: null,
  hour: "services_per_hour",
  item: "services_per_item",
  person: "services_per_person",
  day: "services_per_day",
};

const UNIT_LABEL_KEY: Record<ServiceUnit, DictKey> = {
  fixed: "services_unit_fixed",
  hour: "services_unit_hour",
  item: "services_unit_item",
  person: "services_unit_person",
  day: "services_unit_day",
};

export function ServicesManager({ vendorId: _vendorId }: Props) {
  const { t } = useI18n();
  const [items, setItems] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState<Service | "new" | null>(null);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const r = await api.myServices();
      setItems(r.items ?? []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : t("services_load_failed"));
    } finally {
      setLoading(false);
    }
  }, [t]);

  useEffect(() => {
    void load();
  }, [load]);

  async function remove(id: string) {
    if (!confirm(t("services_confirm_delete"))) return;
    await api.deleteService(id);
    await load();
  }

  async function toggleActive(s: Service) {
    await api.updateService(s.id, {
      name: s.name,
      description: s.description,
      price: s.price,
      unit: s.unit,
      isActive: !s.isActive,
    });
    await load();
  }

  return (
    <section className="mt-6 rounded-xl border bg-[var(--color-card)] p-6">
      <div className="flex items-start justify-between">
        <div>
          <h2 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
            / {t("services_title")}
          </h2>
          <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
            {t("services_hint")}
          </p>
        </div>
        <Button onClick={() => setEditing("new")} size="sm">
          <Plus className="h-4 w-4" />
          {t("services_btn_new")}
        </Button>
      </div>

      {error && (
        <div className="mt-3 rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
          {error}
        </div>
      )}

      {loading ? (
        <div className="mt-6 text-sm text-[var(--color-muted-foreground)]">
          {t("common_loading")}
        </div>
      ) : items.length === 0 ? (
        <div className="mt-6 rounded-lg border border-dashed py-10 text-center text-sm text-[var(--color-muted-foreground)]">
          {t("services_empty")}
        </div>
      ) : (
        <ul className="mt-6 space-y-2">
          {items.map((s) => {
            const suffixKey = UNIT_SUFFIX_KEY[s.unit];
            return (
              <li
                key={s.id}
                className={`rounded-lg border p-4 ${s.isActive ? "" : "opacity-60"}`}
              >
                <div className="flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <div className="flex items-center gap-2">
                      <h3 className="font-medium">{s.name}</h3>
                      {!s.isActive && (
                        <span className="rounded-full border px-2 py-0.5 text-[10px] uppercase">
                          {t("services_inactive")}
                        </span>
                      )}
                    </div>
                    {s.description && (
                      <p className="mt-1 text-sm text-[var(--color-muted-foreground)]">
                        {s.description}
                      </p>
                    )}
                    <div className="mt-2 text-sm font-medium">
                      {formatKZT(s.price)}
                      {suffixKey && (
                        <span className="font-mono text-xs text-[var(--color-muted-foreground)]">
                          {t(suffixKey)}
                        </span>
                      )}
                    </div>
                  </div>
                  <div className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => toggleActive(s)}
                      title={s.isActive ? t("services_disable") : t("services_enable")}
                    >
                      {s.isActive ? (
                        <Power className="h-4 w-4" />
                      ) : (
                        <PowerOff className="h-4 w-4" />
                      )}
                    </Button>
                    <Button variant="ghost" size="icon" onClick={() => setEditing(s)}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => remove(s.id)}
                    >
                      <Trash2 className="h-4 w-4 text-[var(--color-destructive)]" />
                    </Button>
                  </div>
                </div>
              </li>
            );
          })}
        </ul>
      )}

      {editing !== null && (
        <ServiceFormDialog
          initial={editing === "new" ? null : editing}
          onClose={() => setEditing(null)}
          onSaved={() => {
            setEditing(null);
            void load();
          }}
        />
      )}
    </section>
  );
}

function ServiceFormDialog({
  initial,
  onClose,
  onSaved,
}: {
  initial: Service | null;
  onClose: () => void;
  onSaved: () => void;
}) {
  const { t } = useI18n();
  const [name, setName] = useState(initial?.name ?? "");
  const [description, setDescription] = useState(initial?.description ?? "");
  const [price, setPrice] = useState(String(initial?.price ?? "0"));
  const [unit, setUnit] = useState<ServiceUnit>(initial?.unit ?? "fixed");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const unitOptions = SERVICE_UNITS.map((u) => ({
    value: u.value,
    label: t(UNIT_LABEL_KEY[u.value]),
  }));

  async function submit() {
    setBusy(true);
    setError(null);
    const body = {
      name,
      description,
      price: Number(price) || 0,
      unit,
    };
    try {
      if (initial) {
        await api.updateService(initial.id, body);
      } else {
        await api.createService(body);
      }
      onSaved();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : t("services_save_failed"));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-black/40 backdrop-blur-sm p-4">
      <div className="w-full max-w-md rounded-2xl border bg-[var(--color-card)] p-6 shadow-xl">
        <h2 className="font-display text-2xl tracking-[-0.045em]">
          {initial ? t("services_dialog_edit") : t("services_dialog_new")}
        </h2>
        <div className="mt-4 space-y-3">
          <div className="space-y-1.5">
            <label className="text-sm font-medium">{t("services_field_name")}</label>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t("services_field_name_ph")}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">{t("services_field_description")}</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={2}
              className="w-full resize-none rounded-xl border bg-[var(--color-input)]/60 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-ring)]"
            />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-sm font-medium">{t("services_field_price")}</label>
              <Input
                type="number"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-medium">{t("services_field_unit")}</label>
              <Select
                value={unit}
                onChange={(v) => setUnit(v as ServiceUnit)}
                options={unitOptions}
              />
            </div>
          </div>
          {error && (
            <div className="rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
              {error}
            </div>
          )}
        </div>
        <div className="mt-6 flex justify-end gap-2">
          <Button variant="outline" onClick={onClose} disabled={busy}>
            {t("common_cancel")}
          </Button>
          <Button onClick={submit} disabled={busy || !name}>
            {busy ? "..." : t("common_save")}
          </Button>
        </div>
      </div>
    </div>
  );
}
