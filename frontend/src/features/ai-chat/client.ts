import { getToken } from "@/shared/api";
import { API_BASE } from "@/shared/config/env";
import type { Block, ChatMessage } from "./types";

let counter = 1;
const id = () => `m_${Date.now()}_${counter++}`;

export function userMessage(text: string): ChatMessage {
  return { id: id(), role: "user", text };
}

interface ChatBlockRaw {
  type: string;
  data: Record<string, unknown>;
}

interface ChatReplyResponse {
  chatId?: string;
  reply: string;
  blocks?: ChatBlockRaw[];
}

function adaptBlock(b: ChatBlockRaw): Block | null {
  const data = b.data;
  switch (b.type) {
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

export interface SendResult {
  chatId: string;
  message: ChatMessage;
}

/**
 * Send a chat message. Server persists user message + AI reply, returns chatId.
 * On first message of a brand-new chat, chatId may be empty; server creates one and echoes it back.
 */
export async function sendChat(text: string, chatId?: string): Promise<SendResult> {
  const token = getToken();
  const res = await fetch(`${API_BASE}/api/chat`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify({ message: text, chatId: chatId ?? "" }),
  });
  const body: ChatReplyResponse | { error?: string } = await res.json().catch(() => ({}) as never);
  if (!res.ok) {
    const message =
      "error" in body && typeof body.error === "string" ? body.error : `chat failed: ${res.status}`;
    throw new Error(message);
  }
  const data = body as ChatReplyResponse;
  const blocks = (data.blocks ?? [])
    .map(adaptBlock)
    .filter((b): b is Block => b !== null);
  return {
    chatId: data.chatId ?? chatId ?? "",
    message: {
      id: id(),
      role: "ai",
      text: data.reply,
      blocks: blocks.length ? blocks : undefined,
    },
  };
}
