"use client";

import { useCallback, useEffect, useState } from "react";
import type { ChatMessage } from "./types";

const STORAGE_KEY = "qonaqzhai_chat_history";
const MAX_CHATS = 30;

export interface ChatSession {
  id: string;
  title: string;
  createdAt: number;
  updatedAt: number;
  messages: ChatMessage[];
}

interface Index {
  id: string;
  title: string;
  createdAt: number;
  updatedAt: number;
}

function loadAll(): Record<string, ChatSession> {
  if (typeof window === "undefined") return {};
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (!raw) return {};
    return JSON.parse(raw) as Record<string, ChatSession>;
  } catch {
    return {};
  }
}

function saveAll(map: Record<string, ChatSession>): void {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(map));
}

export function newChatId(): string {
  return `c_${Date.now().toString(36)}_${Math.random().toString(36).slice(2, 6)}`;
}

function deriveTitle(messages: ChatMessage[]): string {
  const first = messages.find((m) => m.role === "user" && m.text);
  if (!first?.text) return "New chat";
  const t = first.text.trim().replace(/\s+/g, " ");
  return t.length > 40 ? `${t.slice(0, 40)}…` : t;
}

export function listChats(): Index[] {
  const all = loadAll();
  return Object.values(all)
    .map(({ id, title, createdAt, updatedAt }) => ({ id, title, createdAt, updatedAt }))
    .sort((a, b) => b.updatedAt - a.updatedAt);
}

export function loadChat(id: string): ChatSession | null {
  return loadAll()[id] ?? null;
}

export function saveChat(id: string, messages: ChatMessage[]): ChatSession {
  const all = loadAll();
  const existing = all[id];
  const now = Date.now();
  // strip transient `streaming` flag so reloads don't re-animate
  const sanitized: ChatMessage[] = messages.map((m) => {
    if (!m.streaming) return m;
    const { streaming: _omit, ...rest } = m;
    return rest;
  });
  const next: ChatSession = {
    id,
    title: deriveTitle(sanitized),
    createdAt: existing?.createdAt ?? now,
    updatedAt: now,
    messages: sanitized,
  };
  all[id] = next;
  // cap at MAX_CHATS by recency
  const entries = Object.values(all).sort((a, b) => b.updatedAt - a.updatedAt);
  if (entries.length > MAX_CHATS) {
    for (const stale of entries.slice(MAX_CHATS)) {
      delete all[stale.id];
    }
  }
  saveAll(all);
  return next;
}

export function deleteChat(id: string): void {
  const all = loadAll();
  delete all[id];
  saveAll(all);
}

/** Reactive hook: list of chat sessions for the current browser. */
export function useChatHistory() {
  const [items, setItems] = useState<Index[]>([]);

  const refresh = useCallback(() => {
    setItems(listChats());
  }, []);

  useEffect(() => {
    refresh();
    function onStorage(e: StorageEvent) {
      if (e.key === STORAGE_KEY) refresh();
    }
    function onCustom() {
      refresh();
    }
    window.addEventListener("storage", onStorage);
    window.addEventListener("qonaqzhai:chat-changed", onCustom);
    return () => {
      window.removeEventListener("storage", onStorage);
      window.removeEventListener("qonaqzhai:chat-changed", onCustom);
    };
  }, [refresh]);

  return { items, refresh };
}

export function notifyChatChanged(): void {
  if (typeof window === "undefined") return;
  window.dispatchEvent(new CustomEvent("qonaqzhai:chat-changed"));
}
