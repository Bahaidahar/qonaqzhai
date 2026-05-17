"use client";

import { useCallback, useEffect, useState } from "react";
import { getToken } from "@/shared/api";
import { API_BASE } from "@/shared/config/env";
import type { Block, ChatMessage } from "./types";

function adaptRawBlock(raw: { type: string; data: Record<string, unknown> }): Block | null {
  const data = raw.data ?? {};
  switch (raw.type) {
    case "plan":
      return {
        type: "plan",
        title: String(data.title ?? ""),
        eventType: String(data.eventType ?? ""),
        date: String(data.date ?? ""),
        city: String(data.city ?? ""),
        guests: Number(data.guests ?? 0),
        budget: Number(data.budget ?? 0),
      };
    case "budget":
      return {
        type: "budget",
        total: Number(data.total ?? 0),
        categories: Array.isArray(data.categories)
          ? (data.categories as Array<{ name: string; amount: number; pct: number }>)
          : [],
      };
    case "vendors":
      return {
        type: "vendors",
        query: String(data.query ?? ""),
        items: Array.isArray(data.items)
          ? (data.items as Block extends { type: "vendors"; items: infer I } ? I : never)
          : [],
      };
    default:
      return null;
  }
}

// Server-backed chat history. Each chat is owned by its user and persisted in
// backend `chats` + `chat_messages` tables.

export interface ChatSession {
  id: string;
  userId: string;
  title: string;
  createdAt: string;
  updatedAt: string;
}

interface ChatDetail extends ChatSession {
  messages: Array<{
    id: string;
    role: "user" | "ai";
    text: string;
    blocks?: Array<{ type: string; data: Record<string, unknown> }>;
    createdAt: string;
  }>;
}

const HISTORY_EVENT = "qonaqzhai:chat-changed";

export function notifyChatChanged(): void {
  if (typeof window === "undefined") return;
  window.dispatchEvent(new CustomEvent(HISTORY_EVENT));
}

async function authedFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  const token = getToken();
  if (token) headers.set("Authorization", `Bearer ${token}`);
  if (!headers.has("Content-Type") && init.body) {
    headers.set("Content-Type", "application/json");
  }
  const res = await fetch(`${API_BASE}${path}`, { ...init, headers });
  if (res.status === 204) return undefined as T;
  const body = await res.json().catch(() => null);
  if (!res.ok) {
    const message =
      (body && typeof body === "object" && "error" in body
        ? String((body as { error: unknown }).error)
        : null) ?? `request failed: ${res.status}`;
    throw new Error(message);
  }
  return body as T;
}

export async function listChats(): Promise<ChatSession[]> {
  const r = await authedFetch<{ items: ChatSession[] | null }>("/api/chats");
  return r.items ?? [];
}

export async function loadChat(
  id: string,
): Promise<{ session: ChatSession; messages: ChatMessage[] } | null> {
  try {
    const r = await authedFetch<ChatDetail>(`/api/chats/${id}`);
    let counter = 0;
    const messages: ChatMessage[] = r.messages.map((m) => {
      const blocks: Block[] = (m.blocks ?? [])
        .map(adaptRawBlock)
        .filter((b): b is Block => b !== null);
      return {
        id: m.id || `m_${++counter}`,
        role: m.role,
        text: m.text,
        blocks: blocks.length ? blocks : undefined,
      };
    });
    return {
      session: {
        id: r.id,
        userId: r.userId,
        title: r.title,
        createdAt: r.createdAt,
        updatedAt: r.updatedAt,
      },
      messages,
    };
  } catch {
    return null;
  }
}

export async function deleteChat(id: string): Promise<void> {
  await authedFetch<void>(`/api/chats/${id}`, { method: "DELETE" });
  notifyChatChanged();
}

export async function renameChat(id: string, title: string): Promise<void> {
  await authedFetch<void>(`/api/chats/${id}`, {
    method: "PATCH",
    body: JSON.stringify({ title }),
  });
  notifyChatChanged();
}

// Process-lifetime cache for the sidebar list. Survives route navigations so
// switching tabs doesn't re-trigger a loading spinner each time — we hydrate
// instantly from cache and revalidate in the background.
let cachedChats: ChatSession[] | null = null;
type Subscriber = (next: ChatSession[]) => void;
const subscribers = new Set<Subscriber>();

function setCache(next: ChatSession[]): void {
  cachedChats = next;
  for (const s of subscribers) s(next);
}

async function revalidate(): Promise<void> {
  if (!getToken()) {
    setCache([]);
    return;
  }
  try {
    const list = await listChats();
    setCache(list);
  } catch {
    // keep prior cache on transient error
  }
}

/** React hook: reactive list of chat sessions owned by current user. */
export function useChatHistory() {
  const [items, setItems] = useState<ChatSession[]>(cachedChats ?? []);
  const [loading, setLoading] = useState(cachedChats === null);

  const refresh = useCallback(async () => {
    await revalidate();
    setLoading(false);
  }, []);

  useEffect(() => {
    const sub: Subscriber = (next) => setItems(next);
    subscribers.add(sub);
    // First mount: kick off a fetch. Subsequent mounts read cache, then refresh
    // silently — no UI loading flicker between tabs.
    void revalidate().then(() => setLoading(false));
    function onCustom() {
      void revalidate();
    }
    window.addEventListener(HISTORY_EVENT, onCustom);
    window.addEventListener("qonaqzhai:user-changed", onCustom);
    return () => {
      subscribers.delete(sub);
      window.removeEventListener(HISTORY_EVENT, onCustom);
      window.removeEventListener("qonaqzhai:user-changed", onCustom);
    };
  }, []);

  return { items, loading, refresh };
}
