"use client";

import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { ThreadList } from "@/features/messaging/thread-list";

export default function ThreadsPage() {
  return (
    <AuthGate>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex-1 overflow-y-auto">
          <div className="mx-auto max-w-3xl px-6 py-10">
            <h1 className="font-display text-4xl tracking-[-0.045em]">Messages</h1>
            <p className="mt-2 text-sm text-[var(--color-muted-foreground)]">
              Direct conversations with vendors / customers, scoped to each accepted booking.
            </p>
            <div className="mt-8">
              <ThreadList />
            </div>
          </div>
        </main>
      </div>
    </AuthGate>
  );
}
