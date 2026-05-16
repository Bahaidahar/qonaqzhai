import { test, expect } from "@playwright/test";
import { adminLogin, loginAs, uniqueEmail } from "./helpers";

test("admin suspends user → user cannot login → reactivates", async ({
  browser,
}) => {
  const victimEmail = uniqueEmail("victim");

  // create victim
  const victimCtx = await browser.newContext();
  const victimPage = await victimCtx.newPage();
  await loginAs(victimPage, victimEmail, "password123", "Victim", "customer");
  await victimPage.goto("/");
  // confirm logged in (sidebar shows name)
  await expect(victimPage.getByText("Victim").first()).toBeVisible({
    timeout: 10000,
  });
  await victimCtx.close();

  // admin suspends
  const adminCtx = await browser.newContext();
  const adminPage = await adminCtx.newPage();
  await adminLogin(adminPage);
  await adminPage.goto("/admin/users");
  // find the victim row
  const row = adminPage.locator("tr", { hasText: victimEmail });
  await expect(row).toBeVisible({ timeout: 10000 });
  await row.getByRole("button", { name: /Suspend/i }).click();
  await expect(row.getByText("suspended").first()).toBeVisible({ timeout: 5000 });

  // victim tries to log in fresh → 403
  const reLoginRes = await fetch("http://localhost:8080/api/login", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ email: victimEmail, password: "password123" }),
  });
  expect(reLoginRes.status).toBe(403);

  // reactivate
  await row.getByRole("button", { name: /Activate/i }).click();
  await expect(row.getByText("active").first()).toBeVisible({ timeout: 5000 });

  const ok = await fetch("http://localhost:8080/api/login", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ email: victimEmail, password: "password123" }),
  });
  expect(ok.status).toBe(200);

  await adminCtx.close();
});
