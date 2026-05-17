"use client";

import { Suspense, useCallback, useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { ChatInput } from "@/features/ai-chat/chat-input";
import { ChatMessageView } from "@/features/ai-chat/chat-message";
import { AuthGate } from "@/features/auth/auth-gate";
import { useAuth } from "@/features/auth/context";
import { RedirectIfWrongRole } from "@/features/auth/role-redirect";
import { sendChat, userMessage } from "@/features/ai-chat/client";
import type { ChatMessage } from "@/features/ai-chat/types";
import { loadChat, notifyChatChanged } from "@/features/ai-chat/history";
import { useI18n } from "@/shared/i18n/context";
import { MessageSquare } from "lucide-react";

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
    loadedRef.current = chatIdFromUrl;
    setChatId(chatIdFromUrl);
    void loadChat(chatIdFromUrl).then((res) => {
      if (!res) return;
      // Only apply if URL still matches (user may have switched chats mid-load).
      if (loadedRef.current === chatIdFromUrl) {
        setMessages(res.messages);
      }
    });
  }, [chatIdFromUrl]);

  const send = useCallback(
    async (text: string) => {
      const currentId = chatId;
      // Append the user message optimistically.
      setMessages((prev) => [...prev, userMessage(text)]);
      setThinking(true);
      try {
        const { chatId: resolvedId, message } = await sendChat(text, currentId ?? undefined);
        if (resolvedId && resolvedId !== currentId) {
          setChatId(resolvedId);
          loadedRef.current = resolvedId;
          router.replace(`/?c=${resolvedId}`);
        }
        setMessages((prev) => [...prev, { ...message, streaming: true }]);
        notifyChatChanged();
      } finally {
        setThinking(false);
      }
    },
    [chatId, router],
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
  const { user } = useAuth();
  const firstName = (user?.name ?? user?.email ?? "").split(/[\s@]/)[0];
  return (
    <div className="relative flex h-full flex-1 flex-col items-center justify-center overflow-hidden px-6 py-10">
      <div className="pointer-events-none absolute inset-0 -z-10">
        <div className="absolute inset-0 grid-bg opacity-50" />
        <div className="glow-indigo absolute left-1/2 top-1/3 h-[500px] w-[700px] -translate-x-1/2 -translate-y-1/2 pulse-soft" />
      </div>

      <div className="w-full max-w-2xl">
        <div className="text-center">
          <h1 className="font-display text-6xl tracking-[-0.05em] sm:text-7xl">
            {t("hero_title_a")}
            <br />
            <span className="text-[var(--color-primary)]">
              {firstName ? `${firstName}?` : t("hero_title_b")}
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
