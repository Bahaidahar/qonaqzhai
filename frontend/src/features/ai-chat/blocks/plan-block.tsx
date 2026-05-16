"use client";

import { Calendar, MapPin, Users, Wallet, Sparkles } from "lucide-react";
import { Card, CardContent } from "@/shared/ui/card";
import { Badge } from "@/shared/ui/badge";
import { formatKZT } from "@/shared/lib/utils";
import { useI18n } from "@/shared/i18n/context";
import type { PlanBlock as PlanBlockData } from "@/features/ai-chat/types";

export function PlanBlock({ data }: { data: PlanBlockData }) {
  const { t } = useI18n();
  return (
    <Card className="hover-lift overflow-hidden">
      <div className="flex items-center gap-2 border-b px-4 py-2.5">
        <Sparkles className="h-3.5 w-3.5 text-[var(--color-primary)]" />
        <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
          {t("block_event_plan")}
        </span>
        <Badge variant="primary" className="ml-auto">
          {t("block_draft")}
        </Badge>
      </div>
      <CardContent className="p-5">
        <h3 className="font-display text-2xl tracking-tight">{data.title}</h3>
        <div className="mt-1 text-xs text-[var(--color-muted-foreground)]">
          {data.eventType}
        </div>
        <dl className="mt-5 grid grid-cols-2 gap-3 sm:grid-cols-4">
          <Stat icon={Calendar} label={t("stat_date")} value={data.date} />
          <Stat icon={MapPin} label={t("stat_city")} value={data.city} />
          <Stat icon={Users} label={t("stat_guests")} value={String(data.guests)} />
          <Stat
            icon={Wallet}
            label={t("stat_budget")}
            value={formatKZT(data.budget)}
          />
        </dl>
      </CardContent>
    </Card>
  );
}

function Stat({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: string;
}) {
  return (
    <div className="rounded-xl bg-[var(--color-muted)]/40 p-3">
      <div className="flex items-center gap-1.5 text-[10px] uppercase tracking-wide text-[var(--color-muted-foreground)]">
        <Icon className="h-3 w-3" />
        {label}
      </div>
      <div className="mt-1 text-sm font-medium">{value}</div>
    </div>
  );
}
