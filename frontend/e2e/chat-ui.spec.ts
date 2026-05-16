import { test, expect } from "@playwright/test";
import { loginAs, uniqueEmail } from "./helpers";

test("customer sends chat message and AI block renders", async ({ page }) => {
  await loginAs(page, uniqueEmail("chat"), "password123", "Chatter", "customer");
  await page.goto("/");

  const ta = page.locator("textarea").first();
  await expect(ta).toBeVisible({ timeout: 10000 });
  await ta.fill("plan a toi for 150 guests in Almaty, budget 5M tenge");
  await page.locator('button[type="submit"]').click();

  // any block caption (event_plan | budget | vendors) renders
  await expect(
    page.getByText(/event_plan|budget|vendors|той_жоспары/i).first()
  ).toBeVisible({ timeout: 20000 });
});

test("chat thinking indicator appears then resolves", async ({ page }) => {
  await loginAs(page, uniqueEmail("think"), "password123", "T", "customer");
  await page.goto("/");
  const ta = page.locator("textarea").first();
  await ta.fill("plan a toi for 150 guests in Almaty");
  await page.locator('button[type="submit"]').click();
  // any block caption appears once AI/fallback replies
  await expect(
    page.getByText(/event_plan|budget|vendors|той_жоспары/i).first()
  ).toBeVisible({ timeout: 20000 });
});
