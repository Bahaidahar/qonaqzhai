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
  const adminRes = await fetch(`${BACKEND}/api/login`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({
      email: "admin@qonaqzhai.kz",
      password: "admin12345",
    }),
  });
  const admin = (await adminRes.json()) as { token: string };

  const listRes = await fetch(`${BACKEND}/api/admin/users`, {
    headers: { authorization: `Bearer ${admin.token}` },
  });
  const list = (await listRes.json()) as { items: { id: string; email: string }[] };
  const user = list.items.find((u) => u.email === email);
  if (!user) throw new Error(`user ${email} not found`);

  // find vendor by user
  const vendorsRes = await fetch(`${BACKEND}/api/vendors?status=`, {
    headers: { authorization: `Bearer ${admin.token}` },
  });
  const vendors = (await vendorsRes.json()) as {
    items: { id: string; userId: string }[] | null;
  };
  const vendor = vendors.items?.find((v) => v.userId === user.id);
  if (!vendor) throw new Error(`vendor for ${email} not found`);

  await fetch(`${BACKEND}/api/admin/vendors/${vendor.id}`, {
    method: "PATCH",
    headers: {
      "content-type": "application/json",
      authorization: `Bearer ${admin.token}`,
    },
    body: JSON.stringify({ status: "approved" }),
  });
}
