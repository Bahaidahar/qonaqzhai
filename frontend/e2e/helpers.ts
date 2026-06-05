import type { Page } from "@playwright/test";

const BACKEND = "http://localhost:8080";

export interface TestUser {
  email: string;
  password: string;
  name: string;
}

let counter = 0;
export function uniqueEmail(prefix: string): string {
  counter += 1;
  return `${prefix}_${Date.now()}_${counter}@e2e.test`;
}

/** Signup user directly via backend, then plant token in browser storage. */
export async function loginAs(
  page: Page,
  email: string,
  password: string,
  name: string,
  role: "customer" | "vendor" | "admin"
): Promise<string> {
  // try signup first; if already exists, login
  let res = await fetch(`${BACKEND}/api/signup`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ email, password, name, role }),
  });
  if (res.status === 409) {
    res = await fetch(`${BACKEND}/api/login`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
  }
  if (!res.ok) {
    throw new Error(`login/signup failed for ${email}: ${res.status}`);
  }
  const body = (await res.json()) as { token: string };

  await page.addInitScript((token) => {
    window.localStorage.setItem("qonaqzhai_token", token);
    window.localStorage.setItem("qonaqzhai_locale", "en");
  }, body.token);
  return body.token;
}

/** Plant an existing token (and English locale) into a page's storage. */
export async function plantToken(page: Page, token: string): Promise<void> {
  await page.addInitScript((t) => {
    window.localStorage.setItem("qonaqzhai_token", t);
    window.localStorage.setItem("qonaqzhai_locale", "en");
  }, token);
}

export async function adminLogin(page: Page): Promise<void> {
  const res = await fetch(`${BACKEND}/api/login`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({
      email: "admin@qonaqzhai.kz",
      password: "admin12345",
    }),
  });
  if (!res.ok) throw new Error(`admin login failed: ${res.status}`);
  const body = (await res.json()) as { token: string };
  await page.addInitScript((token) => {
    window.localStorage.setItem("qonaqzhai_token", token);
    window.localStorage.setItem("qonaqzhai_locale", "en");
  }, body.token);
}

export async function approveVendorByEmail(email: string): Promise<void> {
  // Re-login as the vendor user to fetch their vendor.id via /api/me/vendor.
  // (Core has no admin-side vendor list endpoint exposed over HTTP.)
  const vLoginRes = await fetch(`${BACKEND}/api/login`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ email, password: "password123" }),
  });
  if (!vLoginRes.ok) throw new Error(`vendor login failed for ${email}`);
  const vAuth = (await vLoginRes.json()) as { token: string };

  const myRes = await fetch(`${BACKEND}/api/me/vendor`, {
    headers: { authorization: `Bearer ${vAuth.token}` },
  });
  if (!myRes.ok) throw new Error(`/api/me/vendor failed: ${myRes.status}`);
  const vendor = (await myRes.json()) as { id: string };

  const adminRes = await fetch(`${BACKEND}/api/login`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({
      email: "admin@qonaqzhai.kz",
      password: "admin12345",
    }),
  });
  const admin = (await adminRes.json()) as { token: string };

  const approveRes = await fetch(
    `${BACKEND}/api/admin/vendors/${vendor.id}/status`,
    {
      method: "PATCH",
      headers: {
        "content-type": "application/json",
        authorization: `Bearer ${admin.token}`,
      },
      body: JSON.stringify({ status: "approved" }),
    }
  );
  if (!approveRes.ok) {
    throw new Error(`approve failed: ${approveRes.status}`);
  }
}

// --- backend API helpers (drive flows directly through the gateway) ---

async function call<T>(
  path: string,
  init: RequestInit & { token?: string } = {}
): Promise<T> {
  const { token, ...rest } = init;
  const res = await fetch(`${BACKEND}${path}`, {
    ...rest,
    headers: {
      "content-type": "application/json",
      ...(token ? { authorization: `Bearer ${token}` } : {}),
      ...(rest.headers ?? {}),
    },
  });
  if (!res.ok) {
    const body = await res.text().catch(() => "");
    throw new Error(`${init.method ?? "GET"} ${path} -> ${res.status} ${body}`);
  }
  if (res.status === 204) return undefined as T;
  return (await res.json()) as T;
}

/** Sign up a user via backend and return its access token. */
export async function apiSignup(
  email: string,
  name: string,
  role: "customer" | "vendor" | "admin"
): Promise<string> {
  const body = await call<{ token: string }>("/api/signup", {
    method: "POST",
    body: JSON.stringify({ email, password: "password123", name, role }),
  });
  return body.token;
}

export async function apiAdminToken(): Promise<string> {
  const body = await call<{ token: string }>("/api/login", {
    method: "POST",
    body: JSON.stringify({
      email: "admin@qonaqzhai.kz",
      password: "admin12345",
    }),
  });
  return body.token;
}

export interface SeededVendor {
  email: string;
  token: string;
  vendorId: string;
  serviceId: string;
  name: string;
}

/**
 * Create a vendor user, upsert an approved profile with one service.
 * Returns identifiers needed to drive the customer side of a flow.
 */
export async function seedApprovedVendor(
  namePrefix: string,
  opts: { category?: string; price?: number } = {}
): Promise<SeededVendor> {
  const email = uniqueEmail(namePrefix.toLowerCase().replace(/\s+/g, ""));
  const name = `${namePrefix} ${Date.now()}`;
  const token = await apiSignup(email, name, "vendor");

  const vendor = await call<{ id: string }>("/api/me/vendor", {
    method: "POST",
    token,
    body: JSON.stringify({
      name,
      category: opts.category ?? "Venue",
      city: "Almaty",
      description: "E2E seeded vendor.",
      priceFrom: opts.price ?? 200000,
    }),
  });

  const service = await call<{ id: string }>("/api/me/vendor/services", {
    method: "POST",
    token,
    body: JSON.stringify({
      name: "Standard package",
      description: "Default offering.",
      price: opts.price ?? 200000,
      unit: "fixed",
    }),
  });

  const admin = await apiAdminToken();
  await call(`/api/admin/vendors/${vendor.id}/status`, {
    method: "PATCH",
    token: admin,
    body: JSON.stringify({ status: "approved" }),
  });

  return { email, token, vendorId: vendor.id, serviceId: service.id, name };
}

export async function apiCreateBooking(
  token: string,
  body: {
    vendorId: string;
    serviceId?: string;
    eventDate: string;
    guestCount: number;
    note?: string;
    amount?: number;
  }
): Promise<{ id: string; status: string; amount: number }> {
  return call("/api/bookings", {
    method: "POST",
    token,
    body: JSON.stringify(body),
  });
}

export async function apiSetBookingStatus(
  token: string,
  id: string,
  status: string
): Promise<{ id: string; status: string }> {
  return call(`/api/bookings/${id}`, {
    method: "PATCH",
    token,
    body: JSON.stringify({ status }),
  });
}

export async function apiAddCard(token: string): Promise<{ id: string }> {
  return call("/api/cards", {
    method: "POST",
    token,
    body: JSON.stringify({
      number: "4242 4242 4242 4242",
      expMonth: 11,
      expYear: 2030,
      holder: "E2E TESTER",
    }),
  });
}

export async function apiMyBookings(
  token: string
): Promise<{ items: Array<{ id: string; status: string }> | null }> {
  return call("/api/bookings", { token });
}
