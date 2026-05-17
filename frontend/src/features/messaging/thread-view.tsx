"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { Send } from "lucide-react";
import { Button } from "@/shared/ui/button";
import {
  api,
  getToken,
  type BookingThread,
  type ThreadMessage,
} from "@/shared/api";
import { API_BASE } from "@/shared/config/env";
import { useAuth } from "@/features/auth/context";

interface Props {
  threadId: string;
}

function wsUrl(token: string): string {
  const base = API_BASE.replace(/^http/, "ws");
  return `${base}/api/ws?token=${encodeURIComponent(token)}`;
}

/** Realtime thread chat over a WebSocket. Falls back to REST POST on socket down. */
export function ThreadView({ threadId }: Props) {
  const { user } = useAuth();
  const [thread, setThread] = useState<BookingThread | null>(null);
  const [messages, setMessages] = useState<ThreadMessage[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [draft, setDraft] = useState("");
  const [sending, setSending] = useState(false);
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  const loadHistory = useCallback(async () => {
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
    void loadHistory();
  }, [loadHistory]);

  useEffect(() => {
    const token = getToken();
    if (!token) return;
    let cancelled = false;
    let retry: number | null = null;

    function connect() {
      const ws = new WebSocket(wsUrl(token!));
      wsRef.current = ws;
      ws.onopen = () => setConnected(true);
      ws.onmessage = (ev) => {
        try {
          const env = JSON.parse(ev.data) as { op: string; data: ThreadMessage };
          if (env.op !== "thread.message") return;
          if (env.data.threadId !== threadId) return;
          setMessages((prev) => {
            if (prev.some((m) => m.id === env.data.id)) return prev;
            return [...prev, env.data];
          });
        } catch {
          /* ignore */
        }
      };
      ws.onerror = () => setConnected(false);
      ws.onclose = () => {
        setConnected(false);
        if (cancelled) return;
        retry = window.setTimeout(connect, 2000);
      };
    }
    connect();

    return () => {
      cancelled = true;
      if (retry) window.clearTimeout(retry);
      wsRef.current?.close();
      wsRef.current = null;
    };
  }, [threadId]);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  async function send(e: React.FormEvent) {
    e.preventDefault();
    const text = draft.trim();
    if (!text || !thread) return;
    setSending(true);
    try {
      const ws = wsRef.current;
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ op: "message", data: { threadId: thread.id, text } }));
      } else {
        const m = await api.sendThreadMessage(thread.id, text);
        setMessages((prev) => [...prev, m]);
      }
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
    return (
      <div className="p-6 text-sm text-[var(--color-destructive)]">
        {error ?? "thread not found"}
      </div>
    );
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
        <span
          className={`ml-auto inline-flex items-center gap-1.5 text-[10px] ${
            connected ? "text-emerald-500" : "text-[var(--color-muted-foreground)]"
          }`}
        >
          <span
            className={`h-1.5 w-1.5 rounded-full ${
              connected ? "bg-emerald-500" : "bg-[var(--color-muted-foreground)]"
            }`}
          />
          {connected ? "live" : "offline"}
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
