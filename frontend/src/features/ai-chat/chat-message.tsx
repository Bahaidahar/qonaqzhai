"use client";

import { useEffect, useState } from "react";
import { cn } from "@/shared/lib/utils";
import type { ChatMessage } from "@/features/ai-chat/types";
import { BlockRenderer } from "./block-renderer";

interface ChatMessageProps {
  message: ChatMessage;
  onChipClick?: (chip: string) => void;
}

/** Characters revealed per animation tick. Tune for desired feel. */
const STREAM_CHARS_PER_TICK = 2;
const STREAM_TICK_MS = 16;

export function ChatMessageView({ message, onChipClick }: ChatMessageProps) {
  if (message.role === "user") {
    return (
      <div className="flex justify-end">
        <div className="max-w-[80%] rounded-2xl rounded-tr-md bg-[var(--color-primary)] px-4 py-2.5 text-sm text-[var(--color-primary-foreground)]">
          {message.text}
        </div>
      </div>
    );
  }

  return (
    <div className="flex justify-start">
      <span className="mr-3 mt-1 grid h-7 w-7 shrink-0 place-items-center rounded-md bg-[var(--color-primary)] text-[10px] font-semibold text-[var(--color-primary-foreground)]">
        Q
      </span>
      <div className="flex min-w-0 max-w-[85%] flex-col gap-3">
        {message.text && <AIText text={message.text} streaming={!!message.streaming} />}
        {message.blocks?.map((block, i) => (
          <BlockRenderer key={i} block={block} />
        ))}
        {message.chips && message.chips.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {message.chips.map((chip) => (
              <button
                key={chip}
                onClick={() => onChipClick?.(chip)}
                className={cn(
                  "chip-hover rounded-full border bg-[var(--color-card)] px-3 py-1.5 text-xs font-medium"
                )}
              >
                {chip}
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function AIText({ text, streaming }: { text: string; streaming: boolean }) {
  const [shown, setShown] = useState(streaming ? 0 : text.length);

  useEffect(() => {
    if (!streaming) {
      setShown(text.length);
      return;
    }
    setShown(0);
    let i = 0;
    const interval = window.setInterval(() => {
      i = Math.min(i + STREAM_CHARS_PER_TICK, text.length);
      setShown(i);
      if (i >= text.length) {
        window.clearInterval(interval);
      }
    }, STREAM_TICK_MS);
    return () => window.clearInterval(interval);
  }, [text, streaming]);

  const done = shown >= text.length;
  return (
    <div className="text-sm leading-relaxed text-[var(--color-foreground)]">
      {text.slice(0, shown)}
      {!done && (
        <span
          className="ml-0.5 inline-block h-3.5 w-1 translate-y-0.5 animate-pulse bg-[var(--color-primary)]"
          aria-hidden
        />
      )}
    </div>
  );
}
