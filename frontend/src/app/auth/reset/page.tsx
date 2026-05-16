"use client";

import { Suspense, useEffect, useState } from "react";
import Link from "next/link";
import { useSearchParams, useRouter } from "next/navigation";
import { CheckCircle2 } from "lucide-react";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { api, ApiError } from "@/shared/api";
import { Wordmark } from "@/widgets/header/brand";
import { useI18n } from "@/shared/i18n/context";

export default function ResetPasswordPage() {
  return (
    <Suspense fallback={null}>
      <ResetInner />
    </Suspense>
  );
}

function ResetInner() {
  const { t } = useI18n();
  const router = useRouter();
  const params = useSearchParams();
  const [token, setToken] = useState("");
  const [pw, setPw] = useState("");
  const [busy, setBusy] = useState(false);
  const [done, setDone] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const tk = params.get("token");
    if (tk) setToken(tk);
  }, [params]);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await api.resetPassword(token.trim(), pw);
      setDone(true);
      window.setTimeout(() => router.replace("/"), 1500);
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
          {t("auth_reset_title")}
        </h1>

        {done ? (
          <div className="mt-8 rounded-xl border bg-[var(--color-card)] p-6 text-center">
            <CheckCircle2 className="mx-auto h-10 w-10 text-[oklch(0.58_0.14_145)]" />
            <h2 className="mt-3 font-display text-xl">{t("auth_reset_done_title")}</h2>
            <p className="mt-2 text-xs text-[var(--color-muted-foreground)]">
              {t("auth_reset_done_hint")}
            </p>
          </div>
        ) : (
          <form onSubmit={submit} className="mt-8 space-y-4">
            <div className="space-y-1.5">
              <label className="text-sm font-medium">{t("auth_reset_token")}</label>
              <Input
                value={token}
                onChange={(e) => setToken(e.target.value)}
                placeholder={t("auth_reset_token_ph")}
                required
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-medium">{t("auth_reset_new_password")}</label>
              <Input
                type="password"
                value={pw}
                onChange={(e) => setPw(e.target.value)}
                placeholder="••••••••"
                minLength={8}
                required
              />
            </div>
            {error && (
              <div className="rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
                {error}
              </div>
            )}
            <Button type="submit" className="w-full" size="lg" disabled={busy}>
              {busy ? "..." : t("auth_reset_btn")}
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
