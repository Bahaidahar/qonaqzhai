import { test, expect } from "@playwright/test";
import { adminLogin, loginAs, uniqueEmail } from "./helpers";

test("admin redirected to /admin from /", async ({ page }) => {
  await adminLogin(page);
  await page.goto("/");
  await expect(page).toHaveURL(/\/admin$/, { timeout: 5000 });
});

test("vendor cannot access admin pages", async ({ page }) => {
  const email = uniqueEmail("v_block");
  await loginAs(page, email, "password123", "V", "vendor");
  await page.goto("/admin");
  await expect(page.getByText(/Access denied|denied/i)).toBeVisible();
});

test("customer cannot access vendor pages", async ({ page }) => {
  const email = uniqueEmail("c_block");
  await loginAs(page, email, "password123", "C", "customer");
  await page.goto("/vendor");
  await expect(page.getByText(/Access denied|denied/i)).toBeVisible();
});
