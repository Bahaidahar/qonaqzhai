"use client";

import { useEffect, useState } from "react";
import { API_BASE } from "@/shared/config/env";
import { getToken } from "@/shared/api";
import { useI18n } from "@/shared/i18n/context";

interface TimePoint {
  date: string;
  value: number;
}
interface CategoryCount {
  category: string;
  count: number;
}
interface FunnelStage {
  stage: string;
  count: number;
}

async function fetchAuthed<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { Authorization: `Bearer ${getToken() ?? ""}` },
  });
  if (!res.ok) throw new Error(`${path}: ${res.status}`);
  return (await res.json()) as T;
}

export function AdminCharts() {
  const { t } = useI18n();
  const [series, setSeries] = useState<TimePoint[]>([]);
  const [categories, setCategories] = useState<CategoryCount[]>([]);
  const [funnel, setFunnel] = useState<FunnelStage[]>([]);

  useEffect(() => {
    void (async () => {
      const [s, c, f] = await Promise.all([
        fetchAuthed<{ items: TimePoint[] }>("/api/admin/stats/bookings"),
        fetchAuthed<{ items: CategoryCount[] }>("/api/admin/stats/categories"),
        fetchAuthed<{ items: FunnelStage[] }>("/api/admin/stats/funnel"),
      ]);
      setSeries(s.items ?? []);
      setCategories(c.items ?? []);
      setFunnel(f.items ?? []);
    })();
  }, []);

  return (
    <div className="mt-8 grid gap-4 lg:grid-cols-3">
      <Card title={t("charts_bookings_per_day")}>
        <BarChart points={series.slice(-30).map((p) => ({ label: p.date.slice(5), value: Number(p.value) }))} />
      </Card>
      <Card title={t("charts_top_categories")}>
        <BarChart
          points={categories.slice(0, 6).map((c) => ({ label: c.category, value: c.count }))}
        />
      </Card>
      <Card title={t("charts_funnel")}>
        <FunnelChart stages={funnel} />
      </Card>
    </div>
  );
}

function Card({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="rounded-xl border bg-[var(--color-card)] p-4">
      <div className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
        {title}
      </div>
      <div className="mt-4">{children}</div>
    </div>
  );
}

function BarChart({ points }: { points: { label: string; value: number }[] }) {
  const { t } = useI18n();
  if (points.length === 0) {
    return (
      <div className="py-8 text-center text-xs text-[var(--color-muted-foreground)]">
        {t("charts_no_data")}
      </div>
    );
  }
  const max = Math.max(...points.map((p) => p.value), 1);
  return (
    <div className="space-y-1.5">
      {points.map((p) => (
        <div key={p.label} className="flex items-center gap-2 text-xs">
          <span className="w-20 shrink-0 truncate text-[var(--color-muted-foreground)]">
            {p.label}
          </span>
          <div className="relative h-4 flex-1 overflow-hidden rounded-sm bg-[var(--color-muted)]">
            <div
              className="absolute inset-y-0 left-0 bg-[var(--color-primary)] transition-all"
              style={{ width: `${(p.value / max) * 100}%` }}
            />
          </div>
          <span className="w-6 text-right font-mono">{p.value}</span>
        </div>
      ))}
    </div>
  );
}

function FunnelChart({ stages }: { stages: FunnelStage[] }) {
  const { t } = useI18n();
  if (stages.length === 0) {
    return (
      <div className="py-8 text-center text-xs text-[var(--color-muted-foreground)]">
        {t("charts_no_data")}
      </div>
    );
  }
  const max = Math.max(...stages.map((s) => s.count), 1);
  return (
    <div className="space-y-2">
      {stages.map((s) => (
        <div key={s.stage}>
          <div className="flex items-center justify-between text-xs">
            <span className="capitalize">{s.stage}</span>
            <span className="font-mono">{s.count}</span>
          </div>
          <div className="mt-1 h-2 overflow-hidden rounded-full bg-[var(--color-muted)]">
            <div
              className="h-full bg-[var(--color-primary)] transition-all"
              style={{ width: `${(s.count / max) * 100}%` }}
            />
          </div>
        </div>
      ))}
    </div>
  );
}
