"use client";

import { useState } from "react";
import { Star, Send } from "lucide-react";
import { Button } from "@/shared/ui/button";
import { api, ApiError } from "@/shared/api";
import { useI18n } from "@/shared/i18n/context";

interface Props {
  bookingId: string;
  onSubmitted?: () => void;
}

export function ReviewForm({ bookingId, onSubmitted }: Props) {
  const { t } = useI18n();
  const [rating, setRating] = useState(5);
  const [text, setText] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [done, setDone] = useState(false);

  if (done) {
    return (
      <div className="rounded-xl border bg-[var(--color-card)] p-4 text-sm">
        {t("reviews_thanks")}
      </div>
    );
  }

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      await api.submitReview({ bookingId, rating, text });
      setDone(true);
      onSubmitted?.();
    } catch (err) {
      if (err instanceof ApiError) setError(err.message);
      else setError(t("auth_network_error"));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form
      onSubmit={submit}
      className="rounded-xl border bg-[var(--color-card)] p-4"
    >
      <div className="flex items-center gap-1">
        {[1, 2, 3, 4, 5].map((n) => (
          <button
            key={n}
            type="button"
            onClick={() => setRating(n)}
            aria-label={`Rate ${n}`}
          >
            <Star
              className={
                n <= rating
                  ? "h-5 w-5 fill-amber-400 text-amber-400"
                  : "h-5 w-5 text-[var(--color-muted-foreground)]"
              }
            />
          </button>
        ))}
      </div>
      <textarea
        value={text}
        onChange={(e) => setText(e.target.value)}
        rows={3}
        placeholder={t("reviews_share_ph")}
        className="mt-3 w-full resize-none rounded-lg border bg-[var(--color-input)]/60 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-ring)]"
      />
      {error && (
        <p className="mt-2 text-xs text-[var(--color-destructive)]">{error}</p>
      )}
      <Button type="submit" size="sm" className="mt-3" disabled={submitting}>
        <Send className="h-3.5 w-3.5" />
        {submitting ? "..." : t("reviews_submit")}
      </Button>
    </form>
  );
}
