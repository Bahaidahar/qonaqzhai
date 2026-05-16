import { test, expect, type Page } from "@playwright/test";
import { adminLogin, approveVendorByEmail, loginAs, uniqueEmail } from "./helpers";

/**
 * Full QA sweep: 3 roles × 3 locales × 2 themes = 18 combinations.
 * For each: visit every page belonging to that role, assert no JS errors
 * and no broken text (no missing translation keys leaking through).
 */

const LOCALES = ["en", "ru", "kz"] as const;
const THEMES = ["light", "dark"] as const;

type Combo = {
  role: "customer" | "vendor" | "admin";
  locale: (typeof LOCALES)[number];
  theme: (typeof THEMES)[number];
};

const PAGES_BY_ROLE: Record<Combo["role"], string[]> = {
  customer: ["/", "/vendors", "/bookings", "/settings"],
  vendor: ["/vendor", "/vendor/bookings", "/settings"],
  admin: ["/admin", "/admin/users", "/settings"],
};

function presetStorage(page: Page, locale: string, theme: string) {
  return page.addInitScript(
    ({ l, th }) => {
      window.localStorage.setItem("qonaqzhai_locale", l);
      window.localStorage.setItem("qonaqzhai_theme", th);
    },
    { l: locale, th: theme }
  );
}

const combos: Combo[] = [];
for (const role of ["customer", "vendor", "admin"] as const) {
  for (const locale of LOCALES) {
    for (const theme of THEMES) {
      combos.push({ role, locale, theme });
    }
  }
}

// pre-seed one approved vendor so customer catalog has content
test.beforeAll(async () => {
  const email = "sweep_vendor@e2e.test";
  await fetch("http://localhost:8080/api/signup", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({
      email,
      password: "password123",
      name: "Sweep Vendor",
      role: "vendor",
    }),
  }).catch(() => null);
  const loginRes = await fetch("http://localhost:8080/api/login", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ email, password: "password123" }),
  });
  if (loginRes.ok) {
    const body = (await loginRes.json()) as { token: string };
    await fetch("http://localhost:8080/api/vendor", {
      method: "POST",
      headers: {
        "content-type": "application/json",
        authorization: `Bearer ${body.token}`,
      },
      body: JSON.stringify({
        name: "Sweep Studio",
        category: "Venue",
        city: "Almaty",
        priceFrom: 500000,
        description: "QA sweep fixture",
      }),
    });
    await approveVendorByEmail(email);
  }
});

for (const combo of combos) {
  test(`${combo.role} · ${combo.locale} · ${combo.theme}`, async ({
    page,
  }) => {
    const errors: string[] = [];
    page.on("pageerror", (err) => errors.push(err.message));
    page.on("console", (msg) => {
      if (msg.type() === "error") errors.push(msg.text());
    });

    await presetStorage(page, combo.locale, combo.theme);

    if (combo.role === "admin") {
      await adminLogin(page);
    } else {
      await loginAs(
        page,
        uniqueEmail(`sweep_${combo.role}`),
        "password123",
        "Sweeper",
        combo.role
      );
    }

    const pages = PAGES_BY_ROLE[combo.role];
    for (const path of pages) {
      await page.goto(path);
      // wait for any content beyond loading state
      await page.waitForLoadState("networkidle");
      // theme attribute applied
      const themeAttr = await page
        .locator("html")
        .getAttribute("data-theme");
      expect(themeAttr).not.toBeNull();
      // body has some text (not blank)
      const text = await page.locator("body").innerText();
      expect(text.length).toBeGreaterThan(10);
      // no untranslated keys leaking through (placeholder pattern like settings_xxx)
      expect(text).not.toMatch(/auth_|settings_|nav_|hero_|capability_/);
    }

    // ignore well-known noise: Next dev image warnings + favicon
    const meaningful = errors.filter(
      (e) =>
        !e.includes("favicon") &&
        !e.includes("Failed to load resource") &&
        !e.toLowerCase().includes("hydration") &&
        !e.includes("React DevTools")
    );
    expect(meaningful, `console/page errors: ${meaningful.join(" | ")}`).toEqual(
      []
    );
  });
}
