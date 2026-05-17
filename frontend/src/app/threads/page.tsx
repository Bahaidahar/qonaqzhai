"use client";

import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { ThreadList } from "@/features/messaging/thread-list";
import { useI18n } from "@/shared/i18n/context";

export default function ThreadsPage() {
  return (
    <AuthGate>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <ThreadsContent />
        </main>
      </div>
    </AuthGate>
  );
}

function ThreadsContent() {
  const { t } = useI18n();
  return (
    <div className="mx-auto max-w-3xl px-6 py-10">
      <h1 className="font-display text-4xl tracking-[-0.045em]">
        {t("threads_title")}
      </h1>
      <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
        {t("threads_hint")}
      </p>
      <div className="mt-8">
        <ThreadList />
      </div>
    </div>
  );
}
