import type { User, Role } from "@/entities/user/types";
import type {
  Vendor,
  Photo,
  VendorSearchParams,
  VendorSearchResult,
} from "@/entities/vendor/types";
import type { Booking, BookingStatus } from "@/entities/booking/types";
import type { Review } from "@/entities/review/types";
import {
  API_BASE,
  TOKEN_KEY,
  REFRESH_TOKEN_KEY,
} from "@/shared/config/env";

export type { Role, User as ApiUser } from "@/entities/user/types";
export type { Vendor, Photo } from "@/entities/vendor/types";
export type { Booking } from "@/entities/booking/types";
export type { Review } from "@/entities/review/types";

export interface AuthResponse {
  token: string;
  accessToken?: string;
  refreshToken?: string;
  user: User;
}

export interface ChatBlockRaw {
  type: string;
  data: Record<string, unknown>;
}

export interface ChatResponse {
  reply: string;
  blocks?: ChatBlockRaw[];
}

export class ApiError extends Error {
  status: number;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string | null): void {
  if (typeof window === "undefined") return;
  if (token) window.localStorage.setItem(TOKEN_KEY, token);
  else window.localStorage.removeItem(TOKEN_KEY);
}

export function getRefreshToken(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function setRefreshToken(token: string | null): void {
  if (typeof window === "undefined") return;
  if (token) window.localStorage.setItem(REFRESH_TOKEN_KEY, token);
  else window.localStorage.removeItem(REFRESH_TOKEN_KEY);
}

export function photoURL(id: string): string {
  return `${API_BASE}/api/photos/${id}`;
}

// Concurrent 401s share a single refresh attempt so we don't burn the
// single-use rotating refresh token from multiple callers at once.
let refreshing: Promise<boolean> | null = null;

async function refreshTokens(): Promise<boolean> {
  const rt = getRefreshToken();
  if (!rt) return false;
  try {
    const res = await fetch(`${API_BASE}/api/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refreshToken: rt }),
    });
    if (!res.ok) {
      setToken(null);
      setRefreshToken(null);
      return false;
    }
    const data = (await res.json()) as AuthResponse;
    persistAuth(data);
    return true;
  } catch {
    return false;
  }
}

async function request<T>(
  path: string,
  init: RequestInit & { auth?: boolean; raw?: boolean; _retry?: boolean } = {}
): Promise<T> {
  const headers = new Headers(init.headers);
  if (!init.raw) headers.set("Content-Type", "application/json");
  if (init.auth) {
    const token = getToken();
    if (!token) throw new ApiError(401, "no token");
    headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(`${API_BASE}${path}`, { ...init, headers });
  if (res.status === 204) return undefined as T;

  // Silent refresh-and-retry on 401 for authed requests, once.
  if (res.status === 401 && init.auth && !init._retry) {
    const ok = await (refreshing ??= refreshTokens().finally(() => {
      refreshing = null;
    }));
    if (ok) return request<T>(path, { ...init, _retry: true });
  }

  let body: unknown = null;
  try {
    body = await res.json();
  } catch {
    /* no body */
  }
  if (!res.ok) {
    const message =
      (body && typeof body === "object" && "error" in body
        ? String((body as { error: unknown }).error)
        : null) ?? `request failed: ${res.status}`;
    throw new ApiError(res.status, message);
  }
  return body as T;
}

function persistAuth(res: AuthResponse): void {
  setToken(res.accessToken ?? res.token);
  if (res.refreshToken) setRefreshToken(res.refreshToken);
}

function vendorSearchQuery(params: VendorSearchParams): string {
  const q = new URLSearchParams();
  if (params.q) q.set("q", params.q);
  if (params.category) q.set("category", params.category);
  if (params.city) q.set("city", params.city);
  if (params.priceMin != null) q.set("price_min", String(params.priceMin));
  if (params.priceMax != null) q.set("price_max", String(params.priceMax));
  if (params.ratingMin != null) q.set("rating_min", String(params.ratingMin));
  if (params.sort) q.set("sort", params.sort);
  if (params.page != null) q.set("page", String(params.page));
  if (params.limit != null) q.set("limit", String(params.limit));
  const qs = q.toString();
  return qs ? `?${qs}` : "";
}

export const api = {
  async signup(body: {
    email: string;
    password: string;
    name?: string;
    role: Role;
  }): Promise<AuthResponse> {
    const res = await request<AuthResponse>("/api/signup", {
      method: "POST",
      body: JSON.stringify(body),
    });
    persistAuth(res);
    return res;
  },
  async login(body: {
    email: string;
    password: string;
  }): Promise<AuthResponse> {
    const res = await request<AuthResponse>("/api/login", {
      method: "POST",
      body: JSON.stringify(body),
    });
    persistAuth(res);
    return res;
  },
  async refresh(): Promise<AuthResponse | null> {
    const rt = getRefreshToken();
    if (!rt) return null;
    try {
      const res = await request<AuthResponse>("/api/auth/refresh", {
        method: "POST",
        body: JSON.stringify({ refreshToken: rt }),
      });
      persistAuth(res);
      return res;
    } catch {
      setRefreshToken(null);
      return null;
    }
  },
  async logout(): Promise<void> {
    const rt = getRefreshToken();
    if (rt) {
      try {
        await request<void>("/api/auth/logout", {
          method: "POST",
          body: JSON.stringify({ refreshToken: rt }),
        });
      } catch {
        /* best effort */
      }
    }
    setToken(null);
    setRefreshToken(null);
  },
  forgotPassword(email: string): Promise<{ status: string }> {
    return request<{ status: string }>("/api/auth/forgot-password", {
      method: "POST",
      body: JSON.stringify({ email }),
    });
  },
  resetPassword(token: string, newPassword: string): Promise<{ status: string }> {
    return request<{ status: string }>("/api/auth/reset-password", {
      method: "POST",
      body: JSON.stringify({ token, newPassword }),
    });
  },
  me(): Promise<User> {
    return request<User>("/api/me", { auth: true });
  },
  chat(message: string): Promise<ChatResponse> {
    return request<ChatResponse>("/api/chat", {
      method: "POST",
      body: JSON.stringify({ message }),
      auth: true,
    });
  },

  // vendor (self)
  vendorMine(): Promise<Vendor> {
    return request<Vendor>("/api/vendor", { auth: true });
  },
  vendorUpsert(body: {
    name: string;
    category: string;
    city: string;
    description: string;
    priceFrom: number;
  }): Promise<Vendor> {
    return request<Vendor>("/api/vendor", {
      method: "POST",
      body: JSON.stringify(body),
      auth: true,
    });
  },
  async uploadPhoto(file: File): Promise<Photo> {
    const fd = new FormData();
    fd.append("photo", file);
    const token = getToken();
    if (!token) throw new ApiError(401, "no token");
    const res = await fetch(`${API_BASE}/api/vendor/photos`, {
      method: "POST",
      body: fd,
      headers: { Authorization: `Bearer ${token}` },
    });
    const body = await res.json().catch(() => null);
    if (!res.ok) {
      throw new ApiError(
        res.status,
        body && typeof body === "object" && "error" in body
          ? String((body as { error: unknown }).error)
          : `upload failed: ${res.status}`
      );
    }
    return body as Photo;
  },
  deletePhoto(id: string): Promise<void> {
    return request<void>(`/api/vendor/photos/${id}`, {
      method: "DELETE",
      auth: true,
    });
  },

  // public catalog
  vendors(params: VendorSearchParams = {}): Promise<VendorSearchResult> {
    return request<VendorSearchResult>(`/api/vendors${vendorSearchQuery(params)}`);
  },
  vendor(id: string): Promise<Vendor> {
    return request<Vendor>(`/api/vendors/${id}`);
  },

  // bookings
  createBooking(body: {
    vendorId: string;
    eventDate: string;
    guestCount: number;
    note?: string;
    amount?: number;
  }): Promise<Booking> {
    return request<Booking>("/api/bookings", {
      method: "POST",
      body: JSON.stringify(body),
      auth: true,
    });
  },
  bookings(): Promise<{ items: Booking[] | null }> {
    return request<{ items: Booking[] | null }>("/api/bookings", { auth: true });
  },
  updateBooking(id: string, status: BookingStatus): Promise<Booking> {
    return request<Booking>(`/api/bookings/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
      auth: true,
    });
  },

  // reviews
  submitReview(body: {
    bookingId: string;
    rating: number;
    text?: string;
  }): Promise<Review> {
    return request<Review>("/api/reviews", {
      method: "POST",
      body: JSON.stringify(body),
      auth: true,
    });
  },
  reviewsForVendor(vendorId: string): Promise<{ items: Review[] | null }> {
    return request<{ items: Review[] | null }>(`/api/vendors/${vendorId}/reviews`);
  },
  adminDeleteReview(id: string): Promise<void> {
    return request<void>(`/api/admin/reviews/${id}`, {
      method: "DELETE",
      auth: true,
    });
  },

  // admin
  adminUsers(): Promise<{ items: User[] | null }> {
    return request<{ items: User[] | null }>("/api/admin/users", {
      auth: true,
    });
  },
  adminUpdateUser(id: string, status: "active" | "suspended"): Promise<User> {
    return request<User>(`/api/admin/users/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
      auth: true,
    });
  },
  adminVendors(): Promise<VendorSearchResult> {
    return request<VendorSearchResult>("/api/vendors?status=", { auth: true });
  },
  adminUpdateVendor(id: string, status: Vendor["status"]): Promise<Vendor> {
    return request<Vendor>(`/api/admin/vendors/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
      auth: true,
    });
  },
  adminStats(): Promise<Record<string, number>> {
    return request<Record<string, number>>("/api/admin/stats", { auth: true });
  },
};
