"use client";

import { useEffect, useRef, useState, type KeyboardEvent } from "react";
import { Send, Mic, MicOff, Sparkles } from "lucide-react";
import { Button } from "@/shared/ui/button";
import { cn } from "@/shared/lib/utils";
import { useI18n } from "@/shared/i18n/context";
import { useSpeechRecognition } from "./speech-recognition";

interface ChatInputProps {
  onSend: (text: string) => void;
  placeholder?: string;
  autoFocus?: boolean;
  size?: "md" | "lg";
  hint?: boolean;
}

export function ChatInput({
  onSend,
  placeholder,
  autoFocus,
  size = "md",
  hint,
}: ChatInputProps) {
  const [value, setValue] = useState("");
  const { t, locale } = useI18n();
  const ph = placeholder ?? t("input_placeholder_main");

  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  useEffect(() => {
    const el = textareaRef.current;
    if (!el) return;
    el.style.height = "auto";
    const max = size === "lg" ? 240 : 180;
    el.style.height = `${Math.min(el.scrollHeight, max)}px`;
  }, [value, size]);

  // Snapshot of input contents at the moment recording starts so partial
  // transcripts append cleanly without overwriting prior typing.
  const baseRef = useRef("");

  const speech = useSpeechRecognition(locale, (text, isFinal) => {
    const base = baseRef.current;
    const sep = base && !base.endsWith(" ") ? " " : "";
    setValue(`${base}${sep}${text}`.replace(/\s+/g, " "));
    if (isFinal) baseRef.current = `${base}${sep}${text}`;
  });

  useEffect(() => {
    if (speech.listening) baseRef.current = value;
    // intentional: only snapshot when transitioning into listening
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [speech.listening]);

  function submit() {
    const trimmed = value.trim();
    if (!trimmed) return;
    if (speech.listening) speech.stop();
    onSend(trimmed);
    setValue("");
    baseRef.current = "";
  }

  function onKey(e: KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      submit();
    }
  }

  return (
    <div className="w-full">
      <form
        onSubmit={(e) => {
          e.preventDefault();
          submit();
        }}
        className={cn(
          "relative flex items-end gap-2 rounded-2xl border bg-[var(--color-card)] shadow-lg transition focus-within:border-[var(--color-primary)]/40 focus-within:shadow-xl",
          size === "lg" ? "p-2.5" : "p-2",
          speech.listening &&
            "border-[var(--color-primary)] ring-2 ring-[var(--color-primary)]/30"
        )}
      >
        <textarea
          ref={textareaRef}
          autoFocus={autoFocus}
          value={value}
          onChange={(e) => setValue(e.target.value)}
          onKeyDown={onKey}
          placeholder={ph}
          rows={1}
          className={cn(
            "flex-1 resize-none overflow-y-auto break-words bg-transparent px-1 py-2.5 leading-relaxed placeholder:text-[var(--color-muted-foreground)] focus:outline-none",
            size === "lg" ? "text-base" : "text-sm"
          )}
        />
        {speech.supported && (
          <Button
            type="button"
            variant={speech.listening ? "primary" : "ghost"}
            size="icon"
            onClick={speech.toggle}
            aria-label={speech.listening ? "Stop dictation" : "Start dictation"}
            title={speech.listening ? "Stop dictation" : "Dictate"}
            className={cn("relative")}
          >
            {speech.listening ? (
              <MicOff className="h-4 w-4" />
            ) : (
              <Mic className="h-4 w-4" />
            )}
            {speech.listening && (
              <span
                className="absolute inset-0 animate-ping rounded-xl bg-[var(--color-primary)]/40"
                aria-hidden
              />
            )}
          </Button>
        )}
        <Button type="submit" size="icon" disabled={!value.trim()}>
          <Send className="h-4 w-4" />
        </Button>
      </form>
      {speech.error && (
        <p className="mt-2 text-center text-[10px] text-[var(--color-destructive)]">
          {speech.error}
        </p>
      )}
      {hint && !speech.error && (
        <p className="mt-2 flex items-center justify-center gap-1.5 text-[10px] text-[var(--color-muted-foreground)]">
          <Sparkles className="h-2.5 w-2.5" />
          {t("input_hint")}
        </p>
      )}
    </div>
  );
}
