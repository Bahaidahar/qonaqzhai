"use client";

import { use, useEffect, useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { ArrowLeft, MapPin, Calendar, Users, Send, CheckCircle2, CreditCard } from "lucide-react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { useI18n } from "@/shared/i18n/context";
import { useLabels } from "@/shared/i18n/labels";
import { api, photoURL, ApiError, type Vendor, type Booking, type Service } from "@/shared/api";
import { formatKZT } from "@/shared/lib/utils";
import { ReviewList, Stars } from "@/features/reviews/review-list";
import { ServicesList } from "@/features/services/services-list";

interface PageProps {
  params: Promise<{ id: string }>;
}

export default function VendorDetailPage({ params }: PageProps) {
  return (
    <AuthGate allowedRoles={["customer"]}>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <VendorDetailInner params={params} />
        </main>
      </div>
    </AuthGate>
  );
}

function VendorDetailInner({ params }: PageProps) {
  const { t } = useI18n();
  const labels = useLabels();
  const { id } = use(params);
  const [vendor, setVendor] = useState<Vendor | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .vendor(id)
      .then(setVendor)
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) {
    return (
      <div className="p-10 text-sm text-[var(--color-muted-foreground)]">
        {t("common_loading")}
      </div>
    );
  }
  if (!vendor) {
    return (
      <div className="p-10 text-sm text-[var(--color-muted-foreground)]">
        {t("vendor_detail_not_found")}
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-5xl px-6 py-10">
      <Link
        href="/vendors"
        className="link-underline inline-flex items-center gap-1.5 text-xs text-[var(--color-muted-foreground)]"
      >
        <ArrowLeft className="h-3 w-3" />
        {t("vendor_detail_back")}
      </Link>

      <div className="mt-6 grid grid-cols-1 gap-3 sm:grid-cols-3">
        {vendor.photoIds.length > 0 ? (
          <>
            <div className="relative aspect-[16/9] overflow-hidden rounded-2xl sm:col-span-2 sm:aspect-[16/10]">
              <Image
                src={photoURL(vendor.photoIds[0])}
                alt={vendor.name}
                fill
                sizes="66vw"
                className="object-cover"
                unoptimized
              />
            </div>
            <div className="grid grid-cols-2 gap-3 sm:grid-cols-1">
              {vendor.photoIds.slice(1, 3).map((pid) => (
                <div
                  key={pid}
                  className="relative aspect-square overflow-hidden rounded-2xl"
                >
                  <Image
                    src={photoURL(pid)}
                    alt=""
                    fill
                    sizes="33vw"
                    className="object-cover"
                    unoptimized
                  />
                </div>
              ))}
            </div>
          </>
        ) : (
          <div className="aspect-[16/9] rounded-2xl bg-gradient-to-br from-[var(--color-muted)] to-[var(--color-accent)] sm:col-span-3" />
        )}
      </div>

      <div className="mt-8 grid gap-8 lg:grid-cols-[1fr_360px]">
        <div>
          <div className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
            {labels.category(vendor.category)}
          </div>
          <h1 className="mt-3 font-display text-5xl tracking-[-0.045em]">
            {vendor.name}
          </h1>
          <div className="mt-3 flex flex-wrap items-center gap-4 text-sm text-[var(--color-muted-foreground)]">
            <span className="inline-flex items-center gap-1">
              <MapPin className="h-4 w-4" />
              {labels.city(vendor.city)}
            </span>
            {vendor.ratingCount > 0 && (
              <span className="inline-flex items-center gap-2">
                <Stars rating={Math.round(vendor.ratingAvg)} />
                <span className="font-medium text-[var(--color-foreground)]">
                  {vendor.ratingAvg.toFixed(1)}
                </span>
                <span>({vendor.ratingCount})</span>
              </span>
            )}
          </div>
          {vendor.description && (
            <p className="mt-6 text-sm leading-relaxed text-[var(--color-foreground)]">
              {vendor.description}
            </p>
          )}

          <section className="mt-12">
            <h2 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
              / {t("services_title")}
            </h2>
            <div className="mt-3">
              <ServicesList vendorId={vendor.id} showFallbackEmpty />
            </div>
          </section>

          <section className="mt-12">
            <h2 className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
              / {t("reviews_title")}
            </h2>
            <div className="mt-3">
              <ReviewList vendorId={vendor.id} />
            </div>
          </section>
        </div>

        <BookingPanel vendor={vendor} />
      </div>
    </div>
  );
}

function BookingPanel({ vendor }: { vendor: Vendor }) {
  const { t } = useI18n();
  const [eventDate, setEventDate] = useState("");
  const [guestCount, setGuestCount] = useState("");
  const [note, setNote] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [booking, setBooking] = useState<Booking | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [paying, setPaying] = useState(false);
  const [selectedService, setSelectedService] = useState<Service | null>(null);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const amountFromService =
        selectedService && selectedService.unit === "person"
          ? selectedService.price * (Number(guestCount) || 1)
          : selectedService?.price;
      const b = await api.createBooking({
        vendorId: vendor.id,
        serviceId: selectedService?.id,
        eventDate,
        guestCount: Number(guestCount) || 0,
        note,
        amount: amountFromService ?? (vendor.priceFrom > 0 ? vendor.priceFrom : undefined),
      });
      setBooking(b);
    } catch (err) {
      if (err instanceof ApiError) setError(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  async function pay() {
    if (!booking) return;
    setPaying(true);
    try {
      // Backend route: POST /api/bookings/{id}/pay returns redirectUrl
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080"}/api/bookings/${booking.id}/pay`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${window.localStorage.getItem("qonaqzhai_token") ?? ""}`,
          },
        }
      );
      if (!res.ok) {
        const body = await res.json().catch(() => null);
        throw new Error((body as { error?: string })?.error ?? `${res.status}`);
      }
      const data = (await res.json()) as { redirectUrl: string };
      window.location.href = data.redirectUrl;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Payment failed");
    } finally {
      setPaying(false);
    }
  }

  if (booking) {
    return (
      <aside className="rounded-xl border bg-[var(--color-card)] p-6">
        <CheckCircle2 className="mx-auto h-10 w-10 text-[oklch(0.58_0.14_145)]" />
        <h3 className="mt-3 text-center font-display text-xl">
          {t("vendor_detail_sent_title")}
        </h3>
        <p className="mt-2 text-center text-xs text-[var(--color-muted-foreground)]">
          {t("vendor_detail_sent_hint")}
        </p>
        {booking.amount > 0 && (
          <Button
            onClick={pay}
            disabled={paying}
            className="mt-4 w-full"
            size="lg"
            variant="primary"
          >
            <CreditCard className="h-4 w-4" />
            {paying ? "..." : `${t("booking_pay")} ${formatKZT(booking.amount)}`}
          </Button>
        )}
        {error && (
          <p className="mt-2 text-center text-xs text-[var(--color-destructive)]">
            {error}
          </p>
        )}
        <Link href="/bookings">
          <Button className="mt-2 w-full" size="sm" variant="outline">
            {t("vendor_detail_view_bookings")}
          </Button>
        </Link>
      </aside>
    );
  }

  return (
    <aside className="rounded-xl border bg-[var(--color-card)] p-6">
      <div className="text-xs text-[var(--color-muted-foreground)]">
        {t("vendor_detail_starting")}
      </div>
      <div className="mt-1 font-display text-3xl">
        {vendor.priceFrom > 0 ? formatKZT(vendor.priceFrom) : "—"}
      </div>

      <div className="mt-5 space-y-2">
        <div className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
          {t("services_pick")}
        </div>
        <ServicesList
          vendorId={vendor.id}
          selectedId={selectedService?.id}
          onSelect={setSelectedService}
        />
      </div>

      <form onSubmit={submit} className="mt-6 space-y-3">
        <div className="space-y-1.5">
          <label className="flex items-center gap-1.5 text-xs font-medium">
            <Calendar className="h-3 w-3" /> {t("vendor_detail_date")}
          </label>
          <Input
            type="date"
            value={eventDate}
            onChange={(e) => setEventDate(e.target.value)}
            required
          />
        </div>
        <div className="space-y-1.5">
          <label className="flex items-center gap-1.5 text-xs font-medium">
            <Users className="h-3 w-3" /> {t("vendor_detail_guests")}
          </label>
          <Input
            type="number"
            value={guestCount}
            onChange={(e) => setGuestCount(e.target.value)}
            placeholder="150"
          />
        </div>
        <div className="space-y-1.5">
          <label className="text-xs font-medium">{t("vendor_detail_note")}</label>
          <textarea
            value={note}
            onChange={(e) => setNote(e.target.value)}
            placeholder={t("vendor_detail_note_ph")}
            rows={3}
            className="w-full resize-none rounded-xl border bg-[var(--color-input)]/60 px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-ring)]"
          />
        </div>

        {error && (
          <div className="rounded-lg border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/8 px-3 py-2 text-xs text-[var(--color-destructive)]">
            {error}
          </div>
        )}

        <Button
          type="submit"
          className="w-full"
          size="lg"
          disabled={submitting || !eventDate}
        >
          <Send className="h-4 w-4" />
          {submitting ? "..." : t("vendor_detail_btn")}
        </Button>
      </form>
    </aside>
  );
}
