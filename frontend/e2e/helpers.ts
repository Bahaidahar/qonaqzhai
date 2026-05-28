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
): Promise<void> {
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
