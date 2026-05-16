"use client";

import Link from "next/link";
import { ArrowLeft, LogOut, Check, Sun, Moon, Monitor } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { useI18n } from "@/shared/i18n/context";
import { LOCALES, type Locale, type DictKey } from "@/shared/i18n/dict";
import { useTheme, type ThemePref } from "@/features/theme/context";
import { useAuth } from "@/features/auth/context";
import { cn } from "@/shared/lib/utils";

const THEMES: { value: ThemePref; icon: typeof Sun; labelKey: DictKey }[] = [
  { value: "light", icon: Sun, labelKey: "theme_light" },
  { value: "dark", icon: Moon, labelKey: "theme_dark" },
  { value: "system", icon: Monitor, labelKey: "theme_system" },
];

export default function SettingsPage() {
  return (
    <AuthGate>
      <SettingsInner />
    </AuthGate>
  );
}

function SettingsInner() {
  const { t, locale, setLocale } = useI18n();
  const { theme, setTheme } = useTheme();
  const { user, logout } = useAuth();

  return (
    <div className="flex h-screen">
      <ChatSidebar />
      <main className="flex-1 overflow-y-auto">
        <div className="mx-auto max-w-4xl px-6 py-10">
          <Link
            href="/"
            className="link-underline inline-flex items-center gap-1.5 text-xs text-[var(--color-muted-foreground)]"
          >
            <ArrowLeft className="h-3 w-3" />
            {t("settings_back")}
          </Link>

          <h1 className="mt-6 font-display text-5xl tracking-[-0.045em]">
            {t("settings_title")}
          </h1>

          {/* Account */}
          <Section title={t("settings_section_account")}>
            <Row label={t("settings_name")} value={user?.name ?? "—"} />
            <Row label={t("settings_email")} value={user?.email ?? "—"} />
            <Row
              label={t("settings_plan_label")}
              value={t("settings_plan_value")}
            />
          </Section>

          {/* Preferences */}
          <Section title={t("settings_section_pref")}>
            <div className="space-y-2">
              <div className="text-sm font-medium">{t("settings_language")}</div>
              <p className="text-xs text-[var(--color-muted-foreground)]">
                {t("settings_language_hint")}
              </p>
              <div className="mt-3 grid gap-2 sm:grid-cols-3">
                {LOCALES.map((l) => (
                  <LocaleCard
                    key={l.code}
                    code={l.code}
                    label={l.label}
                    short={l.short}
                    active={locale === l.code}
                    onClick={() => setLocale(l.code)}
                  />
                ))}
              </div>
            </div>

            <Divider />

            <div className="space-y-2">
              <div className="text-sm font-medium">
                {t("settings_appearance")}
              </div>
              <p className="text-xs text-[var(--color-muted-foreground)]">
                {t("settings_appearance_hint")}
              </p>
              <div className="mt-3 grid gap-2 sm:grid-cols-3">
                {THEMES.map((opt) => (
                  <ThemeCard
                    key={opt.value}
                    value={opt.value}
                    Icon={opt.icon}
                    label={t(opt.labelKey)}
                    active={theme === opt.value}
                    onClick={() => setTheme(opt.value)}
                  />
                ))}
              </div>
            </div>
          </Section>

          {/* About */}
          <Section title={t("settings_section_about")}>
            <Row label={t("settings_version")} value="0.1.0 · MVP" />
            <Row label={t("settings_model")} value="Gemini 2.5 Flash" />
            <Row label={t("settings_built_in")} value="Almaty, KZ" />
          </Section>

          {/* Sign out */}
          <div className="mt-10 border-t pt-6">
            <button
              onClick={logout}
              className="inline-flex items-center gap-2 rounded-lg border border-[var(--color-destructive)]/30 px-4 py-2 text-sm font-medium text-[var(--color-destructive)] transition-colors hover:bg-[var(--color-destructive)]/8"
            >
              <LogOut className="h-4 w-4" />
              {t("settings_signout")}
            </button>
          </div>
        </div>
      </main>
    </div>
  );
}

function Section({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section className="mt-10">
      <h2 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
        / {title}
      </h2>
      <div className="mt-4 space-y-4 rounded-xl border bg-[var(--color-card)] p-5">
        {children}
      </div>
    </section>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between border-b py-2 text-sm last:border-0 last:pb-0 first:pt-0">
      <span className="text-[var(--color-muted-foreground)]">{label}</span>
      <span className="font-medium">{value}</span>
    </div>
  );
}

function Divider() {
  return <div className="h-px bg-[var(--color-border)]" />;
}

function ThemeCard({
  value,
  Icon,
  label,
  active,
  onClick,
}: {
  value: ThemePref;
  Icon: typeof Sun;
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={cn(
        "group relative flex items-center gap-3 rounded-lg border px-3 py-2.5 text-left transition-all",
        active
          ? "border-[var(--color-primary)] bg-[var(--color-primary)]/8"
          : "hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-muted)]"
      )}
      aria-pressed={active}
      data-value={value}
    >
      <Icon
        className={cn(
          "h-4 w-4",
          active
            ? "text-[var(--color-primary)]"
            : "text-[var(--color-muted-foreground)]"
        )}
      />
      <span className="flex-1 text-sm font-medium">{label}</span>
      {active && <Check className="h-4 w-4 text-[var(--color-primary)]" />}
    </button>
  );
}

function LocaleCard({
  code,
  label,
  short,
  active,
  onClick,
}: {
  code: Locale;
  label: string;
  short: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={cn(
        "group relative flex items-center justify-between rounded-lg border px-3 py-2.5 text-left transition-all",
        active
          ? "border-[var(--color-primary)] bg-[var(--color-primary)]/8"
          : "hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-muted)]"
      )}
      aria-pressed={active}
      lang={code}
    >
      <div className="min-w-0 flex-1">
        <div className="font-mono text-[10px] font-semibold uppercase tracking-wide text-[var(--color-muted-foreground)]">
          {short}
        </div>
        <div className="mt-0.5 text-sm font-medium">{label}</div>
      </div>
      {active && (
        <Check className="h-4 w-4 shrink-0 text-[var(--color-primary)]" />
      )}
    </button>
  );
}
