"use client";

import { Suspense, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { ChatInput } from "@/features/ai-chat/chat-input";
import { ChatMessageView } from "@/features/ai-chat/chat-message";
import { AuthGate } from "@/features/auth/auth-gate";
import { RedirectIfWrongRole } from "@/features/auth/role-redirect";
import { sendChat, userMessage } from "@/features/ai-chat/client";
import type { ChatMessage } from "@/features/ai-chat/types";
import {
  loadChat,
  newChatId,
  notifyChatChanged,
  saveChat,
} from "@/features/ai-chat/history";
import { useI18n } from "@/shared/i18n/context";
import {
  Calendar,
  Users,
  Wallet,
  Store,
  Sparkles,
  MessageSquare,
} from "lucide-react";
import type { DictKey } from "@/shared/i18n/dict";

const SUGGESTIONS: { icon: typeof Sparkles; key: DictKey }[] = [
  { icon: Sparkles, key: "suggestion_toi" },
  { icon: Wallet, key: "suggestion_corporate" },
  { icon: Store, key: "suggestion_photographer" },
  { icon: MessageSquare, key: "suggestion_invitation" },
];

const CAPABILITIES: { icon: typeof Calendar; key: DictKey }[] = [
  { icon: Calendar, key: "capability_plan" },
  { icon: Store, key: "capability_vendors" },
  { icon: Wallet, key: "capability_budget" },
  { icon: Users, key: "capability_guests" },
];

export default function HomePage() {
  return (
    <AuthGate>
      <RedirectIfWrongRole expected="customer" />
      <Suspense fallback={null}>
        <ChatApp />
      </Suspense>
    </AuthGate>
  );
}

function ChatApp() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const chatIdFromUrl = searchParams.get("c");

  const [chatId, setChatId] = useState<string | null>(chatIdFromUrl);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [thinking, setThinking] = useState(false);
  const loadedRef = useRef<string | null>(null);

  useEffect(() => {
    if (!chatIdFromUrl) {
      setChatId(null);
      setMessages([]);
      loadedRef.current = null;
      return;
    }
    if (loadedRef.current === chatIdFromUrl) return;
    const session = loadChat(chatIdFromUrl);
    setChatId(chatIdFromUrl);
    setMessages(session?.messages ?? []);
    loadedRef.current = chatIdFromUrl;
  }, [chatIdFromUrl]);

  const send = useCallback(
    async (text: string) => {
      let id = chatId;
      if (!id) {
        id = newChatId();
        setChatId(id);
        loadedRef.current = id;
        router.replace(`/?c=${id}`);
      }
      setMessages((prev) => {
        const next = [...prev, userMessage(text)];
        saveChat(id!, next);
        notifyChatChanged();
        return next;
      });
      setThinking(true);
      try {
        const reply = await sendChat(text);
        setMessages((prev) => {
          const next = [...prev, { ...reply, streaming: true }];
          saveChat(id!, next); // saveChat strips `streaming` before persist
          notifyChatChanged();
          return next;
        });
      } finally {
        setThinking(false);
      }
    },
    [chatId, router]
  );

  const isEmpty = messages.length === 0;

  return (
    <div className="flex h-screen">
      <ChatSidebar />
      <main className="flex h-screen flex-1 flex-col overflow-hidden">
        {chatId && !isEmpty && <ChatHeader chatId={chatId} />}
        {isEmpty ? (
          <CenteredGreeting onSend={send} />
        ) : (
          <Conversation messages={messages} thinking={thinking} onSend={send} />
        )}
      </main>
    </div>
  );
}

function ChatHeader({ chatId }: { chatId: string }) {
  const { t } = useI18n();
  return (
    <header className="flex h-12 shrink-0 items-center gap-3 border-b bg-[var(--color-card)] px-5">
      <MessageSquare className="h-4 w-4 text-[var(--color-primary)]" />
      <span className="text-sm font-medium">{t("chat_header_label")}</span>
      <span className="font-mono text-[10px] uppercase tracking-widest text-[var(--color-muted-foreground)]">
        {chatId}
      </span>
    </header>
  );
}

function CenteredGreeting({ onSend }: { onSend: (t: string) => void }) {
  const { t } = useI18n();
  const suggestions = useMemo(() => SUGGESTIONS, []);
  return (
    <div className="relative flex h-full flex-1 flex-col items-center justify-center overflow-hidden px-6 py-10">
      <div className="pointer-events-none absolute inset-0 -z-10">
        <div className="absolute inset-0 grid-bg opacity-50" />
        <div className="glow-indigo absolute left-1/2 top-1/3 h-[500px] w-[700px] -translate-x-1/2 -translate-y-1/2 pulse-soft" />
      </div>

      <div className="w-full max-w-2xl">
        <div className="text-center">
          <div className="mb-6 inline-flex items-center gap-2 rounded-full border bg-[var(--color-card)] px-3 py-1 text-xs font-medium">
            <span className="relative flex h-1.5 w-1.5">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-[var(--color-primary)] opacity-75" />
              <span className="relative inline-flex h-1.5 w-1.5 rounded-full bg-[var(--color-primary)]" />
            </span>
            <span>{t("status_online")}</span>
          </div>

          <h1 className="font-display text-6xl tracking-[-0.05em] sm:text-7xl">
            {t("hero_title_a")}
            <br />
            <span className="text-[var(--color-primary)]">
              {t("hero_title_b")}
            </span>
          </h1>
        </div>

        <div className="mt-12">
          <ChatInput
            onSend={onSend}
            placeholder={t("input_placeholder_main")}
            size="lg"
            autoFocus
            hint
          />
        </div>

        <div className="mt-10">
          <div className="grid gap-2 sm:grid-cols-2">
            {suggestions.map((s) => {
              const label = t(s.key);
              return (
                <button
                  key={s.key}
                  onClick={() => onSend(label)}
                  className="chip-hover group flex items-center gap-3 rounded-lg border bg-[var(--color-card)] px-4 py-3 text-left text-sm"
                >
                  <s.icon className="icon-pop h-4 w-4 text-[var(--color-primary)]" />
                  <span className="flex-1">{label}</span>
                </button>
              );
            })}
          </div>
        </div>

        <div className="mt-10 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-xs text-[var(--color-muted-foreground)]">
          {CAPABILITIES.map((c) => (
            <span key={c.key} className="inline-flex items-center gap-1.5">
              <c.icon className="h-3 w-3" />
              {t(c.key)}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}

function Conversation({
  messages,
  thinking,
  onSend,
}: {
  messages: ChatMessage[];
  thinking: boolean;
  onSend: (t: string) => void;
}) {
  const { t } = useI18n();
  return (
    <>
      <div className="flex-1 overflow-y-auto">
        <div className="mx-auto max-w-4xl space-y-6 px-4 py-8 sm:px-6">
          {messages.map((m) => (
            <ChatMessageView key={m.id} message={m} onChipClick={onSend} />
          ))}
          {thinking && <ThinkingIndicator />}
        </div>
      </div>
      <div className="bg-[var(--color-background)]/80 backdrop-blur-md">
        <div className="mx-auto max-w-4xl px-4 py-4 sm:px-6">
          <ChatInput
            onSend={onSend}
            placeholder={t("input_placeholder_reply")}
            hint
          />
        </div>
      </div>
    </>
  );
}

function ThinkingIndicator() {
  const { t } = useI18n();
  return (
    <div className="flex items-center gap-3">
      <span className="grid h-7 w-7 place-items-center rounded-md bg-[var(--color-primary)] text-[10px] font-bold text-[var(--color-primary-foreground)]">
        Q
      </span>
      <span className="flex items-center gap-1.5 text-xs text-[var(--color-muted-foreground)]">
        <span>{t("thinking")}</span>
        <span className="flex gap-1">
          <span className="h-1 w-1 animate-bounce rounded-full bg-current [animation-delay:-0.3s]" />
          <span className="h-1 w-1 animate-bounce rounded-full bg-current [animation-delay:-0.15s]" />
          <span className="h-1 w-1 animate-bounce rounded-full bg-current" />
        </span>
      </span>
    </div>
  );
}
