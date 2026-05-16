"use client";

import { useCallback, useEffect, useState } from "react";
import { Plus, Pencil, Trash2, Power, PowerOff } from "lucide-react";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Select } from "@/shared/ui/select";
import { api, ApiError, type Service, type ServiceUnit } from "@/shared/api";
import { SERVICE_UNITS } from "@/entities/service/types";
import { formatKZT } from "@/shared/lib/utils";

interface Props {
  vendorId: string;
}

const UNIT_SUFFIX: Record<ServiceUnit, string> = {
  fixed: "",
  hour: "/hr",
  item: "/item",
  person: "/person",
  day: "/day",
};

export function ServicesManager({ vendorId: _vendorId }: Props) {
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
      setError(err instanceof ApiError ? err.message : "load failed");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  async function remove(id: string) {
    if (!confirm("Delete this service?")) return;
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
            / Services
          </h2>
          <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
            Услуги, которые вы предлагаете. Customer выбирает одну при бронировании.
          </p>
        </div>
        <Button onClick={() => setEditing("new")} size="sm">
          <Plus className="h-4 w-4" />
          New service
        </Button>
      </div>

      {error && (
        <div className="mt-3 rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
          {error}
        </div>
      )}

      {loading ? (
        <div className="mt-6 text-sm text-[var(--color-muted-foreground)]">Loading…</div>
      ) : items.length === 0 ? (
        <div className="mt-6 rounded-lg border border-dashed py-10 text-center text-sm text-[var(--color-muted-foreground)]">
          No services yet — add your first one.
        </div>
      ) : (
        <ul className="mt-6 space-y-2">
          {items.map((s) => (
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
                        inactive
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
                    <span className="font-mono text-xs text-[var(--color-muted-foreground)]">
                      {UNIT_SUFFIX[s.unit]}
                    </span>
                  </div>
                </div>
                <div className="flex gap-1">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => toggleActive(s)}
                    title={s.isActive ? "Disable" : "Enable"}
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
          ))}
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
  const [name, setName] = useState(initial?.name ?? "");
  const [description, setDescription] = useState(initial?.description ?? "");
  const [price, setPrice] = useState(String(initial?.price ?? "0"));
  const [unit, setUnit] = useState<ServiceUnit>(initial?.unit ?? "fixed");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

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
      setError(err instanceof ApiError ? err.message : "save failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="fixed inset-0 z-50 grid place-items-center bg-black/40 backdrop-blur-sm p-4">
      <div className="w-full max-w-md rounded-2xl border bg-[var(--color-card)] p-6 shadow-xl">
        <h2 className="font-display text-2xl tracking-[-0.045em]">
          {initial ? "Edit service" : "New service"}
        </h2>
        <div className="mt-4 space-y-3">
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Name</label>
            <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Wedding photography" />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={2}
              className="w-full resize-none rounded-xl border bg-[var(--color-input)]/60 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-ring)]"
            />
          </div>
          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="text-sm font-medium">Price (₸)</label>
              <Input
                type="number"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-medium">Unit</label>
              <Select
                value={unit}
                onChange={(v) => setUnit(v as ServiceUnit)}
                options={SERVICE_UNITS}
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
            Cancel
          </Button>
          <Button onClick={submit} disabled={busy || !name}>
            {busy ? "..." : "Save"}
          </Button>
        </div>
      </div>
    </div>
  );
}
