"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import type { Locale } from "@/shared/i18n/dict";

type Listener = (text: string, isFinal: boolean) => void;

interface SpeechAlternative {
  transcript: string;
}
interface SpeechResultItem {
  isFinal: boolean;
  0: SpeechAlternative;
}
interface SpeechResultList {
  length: number;
  [index: number]: SpeechResultItem;
}
interface SpeechRecognitionEventLike {
  resultIndex: number;
  results: SpeechResultList;
}
interface SpeechRecognitionErrorEventLike {
  error: string;
}
interface SpeechRecognitionInstance {
  lang: string;
  interimResults: boolean;
  continuous: boolean;
  onresult: ((e: SpeechRecognitionEventLike) => void) | null;
  onerror: ((e: SpeechRecognitionErrorEventLike) => void) | null;
  onend: (() => void) | null;
  start: () => void;
  stop: () => void;
  abort: () => void;
}

type SpeechRecognitionCtor = new () => SpeechRecognitionInstance;

function getCtor(): SpeechRecognitionCtor | null {
  if (typeof window === "undefined") return null;
  const w = window as unknown as {
    SpeechRecognition?: SpeechRecognitionCtor;
    webkitSpeechRecognition?: SpeechRecognitionCtor;
  };
  return w.SpeechRecognition ?? w.webkitSpeechRecognition ?? null;
}

const LOCALE_TO_BCP47: Record<Locale, string> = {
  kz: "kk-KZ",
  ru: "ru-RU",
  en: "en-US",
};

/**
 * Wraps the browser SpeechRecognition API.
 * Returns null on unsupported browsers — caller should hide the mic button.
 */
export function useSpeechRecognition(locale: Locale, listener: Listener) {
  const [supported, setSupported] = useState(false);
  const [listening, setListening] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const recognitionRef = useRef<SpeechRecognitionInstance | null>(null);
  const listenerRef = useRef(listener);
  listenerRef.current = listener;

  useEffect(() => {
    setSupported(getCtor() !== null);
  }, []);

  const start = useCallback(() => {
    const Ctor = getCtor();
    if (!Ctor) {
      setError("Speech recognition is not supported in this browser");
      return;
    }
    setError(null);
    const r = new Ctor();
    r.lang = LOCALE_TO_BCP47[locale] ?? "en-US";
    r.interimResults = true;
    r.continuous = false;
    r.onresult = (e) => {
      let interim = "";
      let final = "";
      for (let i = e.resultIndex; i < e.results.length; i++) {
        const item = e.results[i];
        const txt = item[0]?.transcript ?? "";
        if (item.isFinal) final += txt;
        else interim += txt;
      }
      const combined = (final + interim).trim();
      listenerRef.current(combined, !!final && !interim);
    };
    r.onerror = (e) => {
      setError(e.error || "speech error");
      setListening(false);
    };
    r.onend = () => {
      setListening(false);
      recognitionRef.current = null;
    };
    recognitionRef.current = r;
    try {
      r.start();
      setListening(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "start failed");
      setListening(false);
    }
  }, [locale]);

  const stop = useCallback(() => {
    recognitionRef.current?.stop();
  }, []);

  const toggle = useCallback(() => {
    if (listening) stop();
    else start();
  }, [listening, start, stop]);

  useEffect(() => () => recognitionRef.current?.abort(), []);

  return { supported, listening, error, start, stop, toggle };
}
