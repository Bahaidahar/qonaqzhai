"use client";

import { useState, type ReactNode } from "react";
import Link from "next/link";
import { useAuth } from "@/features/auth/context";
import { useI18n } from "@/shared/i18n/context";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { ApiError, type Role } from "@/shared/api";
import { cn } from "@/shared/lib/utils";
import { Sparkles, Store, ShieldCheck } from "lucide-react";

const DEMO_ACCOUNTS: {
  role: Role;
  email: string;
  password: string;
  icon: typeof Sparkles;
}[] = [
  {
    role: "customer",
    email: "customer1@demo.kz",
    password: "demo12345",
    icon: Sparkles,
  },
  {
    role: "vendor",
    email: "vendor1@demo.kz",
    password: "demo12345",
    icon: Store,
  },
  {
    role: "admin",
    email: "admin@qonaqzhai.kz",
    password: "admin12345",
    icon: ShieldCheck,
  },
];

interface AuthGateProps {
  children: ReactNode;
  allowedRoles?: Role[];
}

export function AuthGate({ children, allowedRoles }: AuthGateProps) {
  const { user, loading } = useAuth();
  const { t } = useI18n();

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center text-sm text-[var(--color-muted-foreground)]">
        {t("auth_loading")}
      </div>
    );
  }
  if (!user) return <AuthScreen />;
  if (allowedRoles && !allowedRoles.includes(user.role)) {
    return (
      <div className="flex h-screen items-center justify-center px-6 text-center">
        <div>
          <h1 className="font-display text-3xl">Access denied</h1>
          <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
            Your role ({user.role}) cannot view this page.
          </p>
        </div>
      </div>
    );
  }
  return <>{children}</>;
}

function AuthScreen() {
  const { t } = useI18n();
  const { signup, login } = useAuth();
  const [mode, setMode] = useState<"signin" | "signup">("signin");
  const [role, setRole] = useState<Role>("customer");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setBusy(true);
    try {
      if (mode === "signup") {
        await signup(email, password, role, name || undefined);
      } else {
        await login(email, password);
      }
    } catch (err) {
      if (err instanceof ApiError) setError(err.message);
      else setError("Network error");
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
        <h1 className="text-center font-display text-4xl tracking-[-0.045em]">
          {t("auth_welcome")}
        </h1>

        <div className="mt-8 grid grid-cols-2 gap-1 rounded-xl border bg-[var(--color-muted)] p-1">
          {(["signin", "signup"] as const).map((m) => (
            <button
              key={m}
              onClick={() => setMode(m)}
              className={cn(
                "rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                mode === m
                  ? "bg-[var(--color-card)] text-[var(--color-foreground)] shadow-sm"
                  : "text-[var(--color-muted-foreground)] hover:text-[var(--color-foreground)]"
              )}
            >
              {m === "signin" ? t("auth_signin_tab") : t("auth_signup_tab")}
            </button>
          ))}
        </div>

        <form onSubmit={onSubmit} className="mt-6 space-y-4">
          {mode === "signup" && (
            <div className="space-y-2">
              <div className="text-sm font-medium">
                {t("auth_role_question")}
              </div>
              <div className="grid grid-cols-2 gap-2">
                <RoleCard
                  icon={Sparkles}
                  label={t("auth_role_customer")}
                  hint={t("auth_role_customer_hint")}
                  active={role === "customer"}
                  onClick={() => setRole("customer")}
                />
                <RoleCard
                  icon={Store}
                  label={t("auth_role_vendor")}
                  hint={t("auth_role_vendor_hint")}
                  active={role === "vendor"}
                  onClick={() => setRole("vendor")}
                />
              </div>
            </div>
          )}

          {mode === "signup" && (
            <div className="space-y-1.5">
              <label className="text-sm font-medium">{t("auth_name")}</label>
              <Input
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Aigerim"
                autoComplete="name"
              />
            </div>
          )}
          <div className="space-y-1.5">
            <label className="text-sm font-medium">{t("auth_email")}</label>
            <Input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              autoComplete="email"
              required
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-sm font-medium">{t("auth_password")}</label>
            <Input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              autoComplete={
                mode === "signin" ? "current-password" : "new-password"
              }
              required
            />
          </div>

          {error && (
            <div className="rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
              {error}
            </div>
          )}

          <Button type="submit" className="w-full" size="lg" disabled={busy}>
            {busy
              ? "..."
              : mode === "signin"
                ? t("auth_signin_btn")
                : t("auth_signup_btn")}
          </Button>

          {mode === "signin" && (
            <Link
              href="/auth/forgot"
              className="link-underline block text-center text-xs text-[var(--color-muted-foreground)]"
            >
              {t("auth_forgot_link")}
            </Link>
          )}
        </form>

        <p className="mt-6 text-center text-xs text-[var(--color-muted-foreground)]">
          {mode === "signin" ? t("auth_no_account") : t("auth_has_account")}{" "}
          <button
            type="button"
            onClick={() => setMode(mode === "signin" ? "signup" : "signin")}
            className="font-medium text-[var(--color-foreground)] hover:underline"
          >
            {mode === "signin" ? t("auth_signup_tab") : t("auth_signin_tab")}
          </button>
        </p>

        {mode === "signin" && (
          <div className="mt-8">
            <div className="mb-3 flex items-center gap-3">
              <span className="h-px flex-1 bg-[var(--color-border)]" />
              <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
                {t("auth_demo_title")}
              </span>
              <span className="h-px flex-1 bg-[var(--color-border)]" />
            </div>
            <div className="grid grid-cols-3 gap-2">
              {DEMO_ACCOUNTS.map((d) => (
                <button
                  key={d.role}
                  type="button"
                  disabled={busy}
                  onClick={async () => {
                    setError(null);
                    setBusy(true);
                    try {
                      await login(d.email, d.password);
                    } catch (err) {
                      if (err instanceof ApiError) setError(err.message);
                      else setError("Network error");
                    } finally {
                      setBusy(false);
                    }
                  }}
                  className="chip-hover flex flex-col items-center gap-1.5 rounded-xl border bg-[var(--color-card)] px-3 py-3 text-xs font-medium disabled:opacity-50"
                >
                  <d.icon className="h-4 w-4 text-[var(--color-primary)]" />
                  <span className="capitalize">{d.role}</span>
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function RoleCard({
  icon: Icon,
  label,
  hint,
  active,
  onClick,
}: {
  icon: typeof Sparkles;
  label: string;
  hint: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "rounded-xl border p-3 text-left transition-all",
        active
          ? "border-[var(--color-primary)] bg-[var(--color-primary)]/8"
          : "hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-muted)]"
      )}
    >
      <Icon
        className={cn(
          "h-4 w-4",
          active
            ? "text-[var(--color-primary)]"
            : "text-[var(--color-muted-foreground)]"
        )}
      />
      <div className="mt-2 text-sm font-medium">{label}</div>
      <div className="mt-1 text-[10px] leading-tight text-[var(--color-muted-foreground)]">
        {hint}
      </div>
    </button>
  );
}
