"use client";

import { Wallet } from "lucide-react";
import { Card, CardContent } from "@/shared/ui/card";
import { formatKZT } from "@/shared/lib/utils";
import { useI18n } from "@/shared/i18n/context";
import type { BudgetBlock as BudgetBlockData } from "@/features/ai-chat/types";

export function BudgetBlock({ data }: { data: BudgetBlockData }) {
  const { t } = useI18n();
  return (
    <Card className="hover-lift overflow-hidden">
      <div className="flex items-center gap-2 border-b px-4 py-2.5">
        <Wallet className="h-3.5 w-3.5 text-[var(--color-primary)]" />
        <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
          {t("block_budget")}
        </span>
        <span className="ml-auto font-mono text-xs font-semibold">
          {formatKZT(data.total)}
        </span>
      </div>
      <CardContent className="space-y-3 p-5">
        {data.categories.map((c) => (
          <div key={c.name}>
            <div className="flex items-center justify-between text-sm">
              <span className="font-medium">{c.name}</span>
              <span className="text-[var(--color-muted-foreground)]">
                {formatKZT(c.amount)}{" "}
                <span className="font-mono text-[10px]">· {c.pct}%</span>
              </span>
            </div>
            <div className="mt-1.5 h-1.5 overflow-hidden rounded-full bg-[var(--color-muted)]">
              <div
                className="h-full rounded-full bg-[var(--color-primary)] transition-all duration-700"
                style={{ width: `${c.pct}%` }}
              />
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
