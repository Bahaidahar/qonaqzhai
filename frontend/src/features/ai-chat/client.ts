import { api, type ChatBlockRaw } from "@/shared/api";
import type { Block, ChatMessage } from "./types";

let counter = 1;
const id = () => `m_${Date.now()}_${counter++}`;

export function userMessage(text: string): ChatMessage {
  return { id: id(), role: "user", text };
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

export async function sendChat(text: string): Promise<ChatMessage> {
  const res = await api.chat(text);
  const blocks = (res.blocks ?? [])
    .map(adaptBlock)
    .filter((b): b is Block => b !== null);
  return {
    id: id(),
    role: "ai",
    text: res.reply,
    blocks: blocks.length ? blocks : undefined,
  };
}
