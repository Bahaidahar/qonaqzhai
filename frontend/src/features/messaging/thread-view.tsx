"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { Send } from "lucide-react";
import { Button } from "@/shared/ui/button";
import { api, type BookingThread, type ThreadMessage } from "@/shared/api";
import { useAuth } from "@/features/auth/context";

interface Props {
  threadId: string;
}

/** Render a thread + send input, polls every 4s for new messages. */
export function ThreadView({ threadId }: Props) {
  const { user } = useAuth();
  const [thread, setThread] = useState<BookingThread | null>(null);
  const [messages, setMessages] = useState<ThreadMessage[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [draft, setDraft] = useState("");
  const [sending, setSending] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);

  const load = useCallback(async () => {
    try {
      const r = await api.getThread(threadId);
      setThread(r.thread);
      setMessages(r.messages ?? []);
    } catch (err) {
      setError(err instanceof Error ? err.message : "load failed");
    } finally {
      setLoading(false);
    }
  }, [threadId]);

  useEffect(() => {
    void load();
    const t = window.setInterval(() => void load(), 4000);
    return () => window.clearInterval(t);
  }, [load]);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  async function send(e: React.FormEvent) {
    e.preventDefault();
    const text = draft.trim();
    if (!text || !thread) return;
    setSending(true);
    try {
      const m = await api.sendThreadMessage(thread.id, text);
      setMessages((prev) => [...prev, m]);
      setDraft("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "send failed");
    } finally {
      setSending(false);
    }
  }

  if (loading) {
    return <div className="p-6 text-sm text-[var(--color-muted-foreground)]">Loading…</div>;
  }
  if (error || !thread) {
    return <div className="p-6 text-sm text-[var(--color-destructive)]">{error ?? "thread not found"}</div>;
  }

  const myId = user?.id;
  const peerLabel = thread.customerId === myId ? "Vendor" : "Customer";

  return (
    <div className="flex h-full flex-col">
      <header className="flex h-12 shrink-0 items-center gap-3 border-b bg-[var(--color-card)] px-5">
        <span className="text-sm font-medium">{peerLabel}</span>
        <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
          booking {thread.bookingId.slice(0, 8)}
        </span>
      </header>

      <div className="flex-1 overflow-y-auto px-4 py-6">
        <div className="mx-auto flex max-w-3xl flex-col gap-3">
          {messages.length === 0 ? (
            <div className="rounded-xl border border-dashed py-10 text-center text-sm text-[var(--color-muted-foreground)]">
              No messages yet. Say hi 👋
            </div>
          ) : (
            messages.map((m) => {
              const mine = m.senderId === myId;
              return (
                <div key={m.id} className={`flex ${mine ? "justify-end" : "justify-start"}`}>
                  <div
                    className={`max-w-[75%] rounded-2xl px-4 py-2.5 text-sm ${
                      mine
                        ? "rounded-tr-md bg-[var(--color-primary)] text-[var(--color-primary-foreground)]"
                        : "rounded-tl-md border bg-[var(--color-card)]"
                    }`}
                  >
                    {m.text}
                    <div
                      className={`mt-1 text-[10px] ${
                        mine
                          ? "text-[var(--color-primary-foreground)]/70"
                          : "text-[var(--color-muted-foreground)]"
                      }`}
                    >
                      {new Date(m.createdAt).toLocaleTimeString()}
                    </div>
                  </div>
                </div>
              );
            })
          )}
          <div ref={bottomRef} />
        </div>
      </div>

      <form
        onSubmit={send}
        className="border-t bg-[var(--color-background)]/80 backdrop-blur-md"
      >
        <div className="mx-auto flex max-w-3xl items-end gap-2 px-4 py-3">
          <textarea
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                void send(e as unknown as React.FormEvent);
              }
            }}
            rows={1}
            placeholder="Type a message…"
            className="flex-1 resize-none rounded-xl border bg-[var(--color-card)] px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[var(--color-ring)]"
          />
          <Button type="submit" size="icon" disabled={!draft.trim() || sending}>
            <Send className="h-4 w-4" />
          </Button>
        </div>
      </form>
    </div>
  );
}
