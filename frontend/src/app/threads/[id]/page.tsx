"use client";

import { use } from "react";
import { ChatSidebar } from "@/widgets/chat-sidebar/chat-sidebar";
import { AuthGate } from "@/features/auth/auth-gate";
import { ThreadView } from "@/features/messaging/thread-view";

interface PageProps {
  params: Promise<{ id: string }>;
}

export default function ThreadPage({ params }: PageProps) {
  const { id } = use(params);
  return (
    <AuthGate>
      <div className="flex h-screen">
        <ChatSidebar />
        <main className="flex h-screen flex-1 flex-col overflow-hidden">
          <ThreadView threadId={id} />
        </main>
      </div>
    </AuthGate>
  );
}
