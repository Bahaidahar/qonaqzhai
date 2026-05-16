"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import {
  api,
  getToken,
  getRefreshToken,
  setToken,
  setRefreshToken,
  type ApiUser,
  type Role,
} from "@/shared/api";

interface AuthContextValue {
  user: ApiUser | null;
  loading: boolean;
  signup: (
    email: string,
    password: string,
    role: Role,
    name?: string
  ) => Promise<void>;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<ApiUser | null>(null);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    // No tokens at all → guest.
    if (!getToken() && !getRefreshToken()) {
      setUser(null);
      setLoading(false);
      return;
    }
    // If only refresh token survived (e.g. access expired), pre-roll it.
    if (!getToken() && getRefreshToken()) {
      const r = await api.refresh();
      if (!r) {
        setUser(null);
        setLoading(false);
        return;
      }
    }
    try {
      const u = await api.me();
      setUser(u);
    } catch {
      setToken(null);
      setRefreshToken(null);
      setUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  async function signup(
    email: string,
    password: string,
    role: Role,
    name?: string
  ) {
    // api.signup already persists access + refresh tokens via shared/api.
    const res = await api.signup({ email, password, name, role });
    setUser(res.user);
  }

  async function login(email: string, password: string) {
    const res = await api.login({ email, password });
    setUser(res.user);
  }

  async function logout() {
    // api.logout revokes the refresh token server-side, then clears local storage.
    await api.logout();
    setUser(null);
  }

  // Persist current user id in localStorage so feature modules (e.g. chat history)
  // can scope their storage to the active account without a hard dependency on this context.
  useEffect(() => {
    if (typeof window === "undefined") return;
    if (user) {
      window.localStorage.setItem("qonaqzhai_current_user", user.id);
    } else {
      window.localStorage.removeItem("qonaqzhai_current_user");
    }
    window.dispatchEvent(new CustomEvent("qonaqzhai:user-changed"));
  }, [user]);

  return (
    <AuthContext.Provider
      value={{ user, loading, signup, login, logout, refresh }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
