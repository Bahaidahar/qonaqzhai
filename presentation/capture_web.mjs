// Headless Playwright capture of the running web app. Mirrors the screenshots
// the diploma deck embeds — login first, then walk each customer route.
import { chromium } from "/Users/bahtiyarelik/Developer/diploma/frontend/node_modules/@playwright/test/index.mjs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const OUT = path.join(__dirname, "screens");
const APP = process.env.APP_URL ?? "http://localhost:3000";
const API = process.env.API_URL ?? "http://localhost:8080";

async function loginToken(email, password) {
  const r = await fetch(`${API}/api/login`, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!r.ok) throw new Error(`login ${email}: ${r.status}`);
  return (await r.json()).token;
}

async function capture(page, route, name) {
  await page.goto(`${APP}${route}`, { waitUntil: "networkidle" });
  await page.waitForTimeout(800);
  await page.screenshot({ path: path.join(OUT, `web-${name}.png`), fullPage: false });
  console.log(`✓ web-${name}.png`);
}

const SHOTS = [
  ["/", "customer-hero"],
  ["/vendors", "customer-catalog"],
  ["/bookings", "customer-bookings"],
  ["/settings", "settings"],
  ["/notifications", "notifications"],
];

const VENDOR_SHOTS = [
  ["/vendor", "vendor-profile"],
  ["/vendor/bookings", "vendor-inbox"],
];

async function run() {
  const browser = await chromium.launch();
  const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 } });
  const page = await ctx.newPage();

  // customer
  const custToken = await loginToken("customer1@demo.kz", "demo12345");
  await ctx.addInitScript((t) => {
    localStorage.setItem("qonaqzhai_token", t);
    localStorage.setItem("qonaqzhai_locale", "en");
    localStorage.setItem("qonaqzhai_theme", "light");
  }, custToken);
  for (const [route, name] of SHOTS) await capture(page, route, name);

  // vendor detail (use a known vendor)
  const vendors = await (await fetch(`${API}/api/vendors?limit=1`)).json();
  if (vendors.items?.[0]) {
    await capture(page, `/vendors/${vendors.items[0].id}`, "customer-vendor-detail");
  }

  // vendor
  await ctx.clearCookies();
  await page.context().clearPermissions();
  const vendToken = await loginToken("vendor1@demo.kz", "demo12345");
  await ctx.addInitScript((t) => {
    localStorage.setItem("qonaqzhai_token", t);
    localStorage.setItem("qonaqzhai_locale", "en");
    localStorage.setItem("qonaqzhai_theme", "light");
  }, vendToken);
  for (const [route, name] of VENDOR_SHOTS) await capture(page, route, name);

  await browser.close();
}

run().catch((e) => {
  console.error(e);
  process.exit(1);
});
