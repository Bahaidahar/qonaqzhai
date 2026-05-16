"use client";

import { useState } from "react";
import Link from "next/link";
import { Mail, ArrowLeft } from "lucide-react";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { useI18n } from "@/shared/i18n/context";
import { api, ApiError } from "@/shared/api";
import { Wordmark } from "@/widgets/header/brand";

export default function ForgotPasswordPage() {
  const { t } = useI18n();
  const [email, setEmail] = useState("");
  const [busy, setBusy] = useState(false);
  const [done, setDone] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await api.forgotPassword(email.trim());
      setDone(true);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : t("auth_network_error"));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden px-6 py-10">
      <div className="pointer-events-none absolute inset-0 -z-10">
        <div className="absolute inset-0 grid-bg opacity-50" />
        <div className="glow-indigo absolute left-1/2 top-1/2 h-[500px] w-[700px] -translate-x-1/2 -translate-y-1/2 pulse-soft" />
      </div>

      <div className="w-full max-w-md">
        <Wordmark className="mb-6 block text-center" />
        <h1 className="text-center font-display text-3xl tracking-[-0.045em]">
          {t("auth_forgot_title")}
        </h1>
        <p className="mt-2 text-center text-sm text-[var(--color-muted-foreground)]">
          {t("auth_forgot_hint")}
        </p>

        {done ? (
          <div className="mt-8 rounded-xl border bg-[var(--color-card)] p-6 text-center">
            <Mail className="mx-auto h-8 w-8 text-[var(--color-primary)]" />
            <h2 className="mt-3 font-display text-xl">{t("auth_forgot_done_title")}</h2>
            <p className="mt-2 text-xs text-[var(--color-muted-foreground)]">
              {t("auth_forgot_done_hint")}
            </p>
            <Link href="/auth/reset">
              <Button className="mt-4 w-full" size="sm">
                {t("auth_forgot_enter_token")}
              </Button>
            </Link>
            <Link
              href="/"
              className="link-underline mt-4 inline-flex items-center gap-1.5 text-xs text-[var(--color-muted-foreground)]"
            >
              <ArrowLeft className="h-3 w-3" />
              {t("auth_forgot_back")}
            </Link>
          </div>
        ) : (
          <form onSubmit={submit} className="mt-8 space-y-4">
            <div className="space-y-1.5">
              <label className="text-sm font-medium">{t("auth_email")}</label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com"
                required
              />
            </div>
            {error && (
              <div className="rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
                {error}
              </div>
            )}
            <Button type="submit" className="w-full" size="lg" disabled={busy}>
              {busy ? "..." : t("auth_forgot_btn_send")}
            </Button>
            <Link
              href="/"
              className="link-underline mx-auto block w-fit text-xs text-[var(--color-muted-foreground)]"
            >
              {t("auth_forgot_back")}
            </Link>
          </form>
        )}
      </div>
    </div>
  );
}
