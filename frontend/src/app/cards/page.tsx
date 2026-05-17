"use client";

import { useCallback, useEffect, useState } from "react";
import { CreditCard, Trash2, Check } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { Button } from "@/shared/ui/button";
import { useI18n } from "@/shared/i18n/context";
import { api, type PaymentCard } from "@/shared/api";

export default function CardsPage() {
  return (
    <AuthGate allowedRoles={["customer"]}>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <CardsView />
        </main>
      </div>
    </AuthGate>
  );
}

function CardsView() {
  const { t } = useI18n();
  const [items, setItems] = useState<PaymentCard[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    const res = await api.listCards();
    setItems(res.items ?? []);
    setLoading(false);
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  async function remove(id: string) {
    await api.deleteCard(id);
    await load();
  }

  async function makeDefault(id: string) {
    await api.setDefaultCard(id);
    await load();
  }

  return (
    <div className="mx-auto max-w-3xl px-6 py-10">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-4xl tracking-[-0.045em]">
            {t("cards_title")}
          </h1>
          <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
            {t("cards_hint")}
          </p>
        </div>
        <Button onClick={() => setShowForm((v) => !v)}>
          <CreditCard className="h-4 w-4" />
          {showForm ? t("common_cancel") : t("cards_add")}
        </Button>
      </div>

      {showForm && (
        <CardForm
          onSaved={() => {
            setShowForm(false);
            void load();
          }}
        />
      )}

      {loading ? (
        <div className="mt-10 text-sm text-[var(--color-muted-foreground)]">
          {t("common_loading")}
        </div>
      ) : items.length === 0 ? (
        <div className="mt-8 rounded-xl border border-dashed py-16 text-center text-sm text-[var(--color-muted-foreground)]">
          {t("cards_empty")}
        </div>
      ) : (
        <ul className="mt-8 space-y-3">
          {items.map((c) => (
            <li
              key={c.id}
              className="flex items-center justify-between rounded-xl border bg-[var(--color-card)] p-5"
            >
              <div className="flex items-center gap-4">
                <div className="grid h-10 w-14 place-items-center rounded-md bg-[var(--color-muted)] text-xs font-semibold uppercase">
                  {c.brand}
                </div>
                <div>
                  <div className="font-mono text-base">•••• {c.last4}</div>
                  <div className="text-xs text-[var(--color-muted-foreground)]">
                    {String(c.expMonth).padStart(2, "0")}/{c.expYear}
                    {c.holder ? ` · ${c.holder}` : ""}
                  </div>
                </div>
                {c.isDefault && (
                  <span className="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2 py-0.5 text-[10px] font-medium uppercase tracking-wide text-emerald-700">
                    {t("cards_default")}
                  </span>
                )}
              </div>
              <div className="flex gap-2">
                {!c.isDefault && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => makeDefault(c.id)}
                  >
                    <Check className="h-4 w-4" />
                    {t("cards_make_default")}
                  </Button>
                )}
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => remove(c.id)}
                >
                  <Trash2 className="h-4 w-4" />
                  {t("common_delete")}
                </Button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

interface CardFormProps {
  onSaved: () => void;
}

function CardForm({ onSaved }: CardFormProps) {
  const { t } = useI18n();
  const [number, setNumber] = useState("");
  const [exp, setExp] = useState("");
  const [holder, setHolder] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  function formatNumber(raw: string): string {
    const digits = raw.replace(/\D/g, "").slice(0, 19);
    return digits.replace(/(.{4})/g, "$1 ").trim();
  }

  function formatExp(raw: string): string {
    const digits = raw.replace(/\D/g, "").slice(0, 4);
    if (digits.length <= 2) return digits;
    return `${digits.slice(0, 2)}/${digits.slice(2)}`;
  }

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    const [mm, yy] = exp.split("/");
    const expMonth = parseInt(mm, 10);
    const expYear = parseInt(yy, 10);
    if (!expMonth || !expYear) {
      setError(t("cards_err_exp"));
      return;
    }
    setSaving(true);
    try {
      await api.createCard({
        number,
        expMonth,
        expYear: expYear < 100 ? 2000 + expYear : expYear,
        holder,
      });
      onSaved();
    } catch (e) {
      setError(e instanceof Error ? e.message : "failed");
    } finally {
      setSaving(false);
    }
  }

  return (
    <form
      onSubmit={submit}
      className="mt-6 space-y-4 rounded-xl border bg-[var(--color-card)] p-5"
    >
      <div>
        <label className="text-xs uppercase tracking-wide text-[var(--color-muted-foreground)]">
          {t("cards_field_number")}
        </label>
        <input
          inputMode="numeric"
          autoComplete="cc-number"
          placeholder="4242 4242 4242 4242"
          value={number}
          onChange={(e) => setNumber(formatNumber(e.target.value))}
          className="mt-1 w-full rounded-md border bg-transparent px-3 py-2 font-mono text-sm"
        />
      </div>
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="text-xs uppercase tracking-wide text-[var(--color-muted-foreground)]">
            {t("cards_field_exp")}
          </label>
          <input
            inputMode="numeric"
            autoComplete="cc-exp"
            placeholder="MM/YY"
            value={exp}
            onChange={(e) => setExp(formatExp(e.target.value))}
            className="mt-1 w-full rounded-md border bg-transparent px-3 py-2 font-mono text-sm"
          />
        </div>
        <div>
          <label className="text-xs uppercase tracking-wide text-[var(--color-muted-foreground)]">
            {t("cards_field_holder")}
          </label>
          <input
            autoComplete="cc-name"
            placeholder="NAME ON CARD"
            value={holder}
            onChange={(e) => setHolder(e.target.value.toUpperCase())}
            className="mt-1 w-full rounded-md border bg-transparent px-3 py-2 text-sm"
          />
        </div>
      </div>
      {error && (
        <div className="text-sm text-red-600">{error}</div>
      )}
      <div className="flex justify-end">
        <Button type="submit" disabled={saving}>
          {saving ? t("common_loading") : t("cards_save")}
        </Button>
      </div>
      <p className="text-[11px] text-[var(--color-muted-foreground)]">
        {t("cards_disclaimer")}
      </p>
    </form>
  );
}
