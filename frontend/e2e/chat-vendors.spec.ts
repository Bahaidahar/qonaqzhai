import { test, expect } from "@playwright/test";
import { seedApprovedVendor, loginAs, uniqueEmail } from "./helpers";

// The chat planner recognises a category keyword and answers with real
// approved vendors from the catalog (vendors block).
test("chat recommends real vendors for a category request", async ({ page }) => {
  // guarantees at least one approved decor vendor exists in the catalog
  await seedApprovedVendor("ChatDecor", { category: "Decor & Florists" });

  await loginAs(page, uniqueEmail("chatv"), "password123", "ChatV", "customer");
  await page.goto("/");

  const ta = page.locator("textarea").first();
  await expect(ta).toBeVisible({ timeout: 10000 });
  await ta.fill("нужен декор на свадьбу");
  await page.locator('button[type="submit"]').click();

  // a vendors block with at least one bookable vendor card surfaces in chat
  const card = page.locator("main div.w-64").first();
  await expect(card).toBeVisible({ timeout: 20000 });
  await expect(card.getByRole("button", { name: /Book now/i })).toBeVisible();
});
