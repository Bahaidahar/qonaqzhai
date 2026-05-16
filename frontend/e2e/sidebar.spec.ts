import { test, expect } from "@playwright/test";
import { loginAs, uniqueEmail } from "./helpers";

test("sidebar collapse and expand", async ({ page }) => {
  await loginAs(page, uniqueEmail("sb"), "password123", "Sidebar", "customer");
  await page.goto("/");

  // visible by default
  await expect(page.getByRole("link", { name: /Vendors|Подрядчики|Мердігерлер/ })).toBeVisible();

  // collapse
  await page.getByRole("button", { name: /Collapse sidebar/i }).click();
  await expect(page.getByRole("button", { name: /Open sidebar/i })).toBeVisible();

  // expand
  await page.getByRole("button", { name: /Open sidebar/i }).click();
  await expect(page.getByRole("link", { name: /Vendors|Подрядчики|Мердігерлер/ })).toBeVisible();
});

test("sidebar state persists across reload", async ({ page }) => {
  await loginAs(page, uniqueEmail("sb_persist"), "password123", "P", "customer");
  await page.goto("/");
  await page.getByRole("button", { name: /Collapse sidebar/i }).click();
  await page.reload();
  // floating expand button shows after reload (still collapsed)
  await expect(page.getByRole("button", { name: /Open sidebar/i })).toBeVisible();
});
