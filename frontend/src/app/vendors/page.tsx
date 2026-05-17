"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { Search, MapPin, ArrowRight, Star } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { Input } from "@/shared/ui/input";
import { Select } from "@/shared/ui/select";
import { useI18n } from "@/shared/i18n/context";
import { useLabels } from "@/shared/i18n/labels";
import { api, photoURL, type Vendor } from "@/shared/api";
import { formatKZT } from "@/shared/lib/utils";

const CATEGORIES = ["All", "Venue", "Catering", "Music & DJ", "Photo & Video", "Decor & Florists", "Cakes"];
// MVP locked to Almaty — city filter is implicit.
const FIXED_CITY = "Almaty";

type SortKey = "newest" | "price_asc" | "price_desc" | "rating_desc";

export default function VendorsPage() {
  return (
    <AuthGate allowedRoles={["customer"]}>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <VendorsCatalog />
        </main>
      </div>
    </AuthGate>
  );
}

function VendorsCatalog() {
  const { t } = useI18n();
  const [vendors, setVendors] = useState<Vendor[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [category, setCategory] = useState("All");
  const [sort, setSort] = useState<SortKey>("newest");
  const [priceMax, setPriceMax] = useState<string>("");
  const [search, setSearch] = useState("");

  // Debounce: free-text query and price filter trigger after a short pause.
  useEffect(() => {
    const timer = window.setTimeout(() => {
      setLoading(true);
      api
        .vendors({
          q: search || undefined,
          category: category === "All" ? undefined : category,
          city: FIXED_CITY,
          priceMax: priceMax ? Number(priceMax) : undefined,
          sort,
          limit: 24,
        })
        .then((r) => {
          setVendors(r.items ?? []);
          setTotal(r.total ?? 0);
        })
        .finally(() => setLoading(false));
    }, 250);
    return () => window.clearTimeout(timer);
  }, [search, category, sort, priceMax]);

  return (
    <div className="mx-auto max-w-6xl px-6 py-10">
      <h1 className="font-display text-4xl tracking-[-0.045em]">
        {t("vendors_title")}
      </h1>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        {t("vendors_hint")} · <span className="font-mono">{total}</span>{" "}
        {t("vendors_total_suffix")}
      </p>

      <div className="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-[var(--color-muted-foreground)]" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder={t("vendors_search_ph")}
            className="pl-9"
          />
        </div>
        <Input
          type="number"
          value={priceMax}
          onChange={(e) => setPriceMax(e.target.value)}
          placeholder={t("vendors_price_max_ph")}
          className="sm:max-w-40"
        />
        <div className="sm:max-w-48">
          <Select
            value={sort}
            onChange={(v) => setSort(v as SortKey)}
            options={[
              { value: "newest", label: t("vendors_sort_newest") },
              { value: "price_asc", label: t("vendors_sort_price_asc") },
              { value: "price_desc", label: t("vendors_sort_price_desc") },
              { value: "rating_desc", label: t("vendors_sort_rating_desc") },
            ]}
            aria-label={t("vendors_sort_aria")}
          />
        </div>
      </div>

      <div className="mt-4 flex flex-col gap-3">
        <Pills value={category} onChange={setCategory} options={CATEGORIES} kind="category" />
      </div>

      {loading ? (
        <div className="mt-10 text-sm text-[var(--color-muted-foreground)]">
          {t("common_loading")}
        </div>
      ) : vendors.length === 0 ? (
        <div className="mt-10 rounded-xl border border-dashed py-16 text-center text-sm text-[var(--color-muted-foreground)]">
          {t("vendors_empty")}
        </div>
      ) : (
        <div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {vendors.map((v) => (
            <VendorCard key={v.id} vendor={v} />
          ))}
        </div>
      )}
    </div>
  );
}

function Pills({
  value,
  onChange,
  options,
  kind,
}: {
  value: string;
  onChange: (v: string) => void;
  options: string[];
  kind: "category" | "city";
}) {
  const { t } = useI18n();
  const labels = useLabels();
  const labelFor = (v: string) => {
    if (v === "All") return t("vendors_filter_all");
    return kind === "category" ? labels.category(v) : labels.city(v);
  };
  return (
    <div className="flex flex-wrap gap-1.5">
      {options.map((o) => (
        <button
          key={o}
          onClick={() => onChange(o)}
          className={
            value === o
              ? "rounded-full border border-[var(--color-primary)] bg-[var(--color-primary)]/10 px-3 py-1.5 text-xs font-medium text-[var(--color-primary)]"
              : "chip-hover rounded-full border bg-[var(--color-card)] px-3 py-1.5 text-xs font-medium"
          }
        >
          {labelFor(o)}
        </button>
      ))}
    </div>
  );
}

function VendorCard({ vendor: v }: { vendor: Vendor }) {
  const { t } = useI18n();
  const labels = useLabels();
  const cover = v.photoIds[0];
  return (
    <Link href={`/vendors/${v.id}`} className="group">
      <div className="hover-lift overflow-hidden rounded-xl border bg-[var(--color-card)]">
        <div className="relative aspect-[4/3] overflow-hidden">
          {cover ? (
            <Image
              src={photoURL(cover)}
              alt={v.name}
              fill
              sizes="(max-width: 768px) 100vw, 33vw"
              className="img-zoom object-cover"
              unoptimized
            />
          ) : (
            <div className="img-zoom h-full w-full bg-gradient-to-br from-[var(--color-muted)] to-[var(--color-accent)]" />
          )}
        </div>
        <div className="p-4">
          <div className="flex items-start justify-between gap-2">
            <div className="min-w-0">
              <div className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
                {labels.category(v.category)}
              </div>
              <h3 className="mt-1 truncate font-semibold">{v.name}</h3>
              <div className="mt-1 flex items-center gap-1 text-xs text-[var(--color-muted-foreground)]">
                <MapPin className="h-3 w-3" />
                {labels.city(v.city)}
              </div>
              {v.ratingCount > 0 && (
                <div className="mt-1 flex items-center gap-1 text-xs">
                  <Star className="h-3 w-3 fill-amber-400 text-amber-400" />
                  <span className="font-medium">{v.ratingAvg.toFixed(1)}</span>
                  <span className="text-[var(--color-muted-foreground)]">
                    ({v.ratingCount})
                  </span>
                </div>
              )}
            </div>
            <ArrowRight className="arrow-slide h-4 w-4 text-[var(--color-muted-foreground)]" />
          </div>
          {v.priceFrom > 0 && (
            <div className="mt-3 border-t pt-3 text-sm">
              <span className="text-[var(--color-muted-foreground)]">
                {t("vendors_from")}{" "}
              </span>
              <span className="font-medium">{formatKZT(v.priceFrom)}</span>
            </div>
          )}
        </div>
      </div>
    </Link>
  );
}
