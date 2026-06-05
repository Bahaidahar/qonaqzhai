import { test, expect } from "@playwright/test";
import { seedApprovedVendor, loginAs, uniqueEmail } from "./helpers";

// Conversational refinement: a first category ask returns a plan + vendors;
// a "show me others" follow-up returns a fresh batch that excludes the vendors
// already shown in the same chat.
test("chat shows other vendor options on a follow-up request", async ({
  page,
}) => {
  // enough cake vendors so a second, distinct batch exists
  for (let i = 0; i < 4; i++) {
    await seedApprovedVendor(`MoreCake${i}`, { category: "Cakes" });
  }

  await loginAs(page, uniqueEmail("morec"), "password123", "More", "customer");
  await page.goto("/");

  const nameSel = "main div.w-64 .font-semibold";

  // first ask → plan + vendor cards
  const ta = page.locator("textarea").first();
  await ta.fill("нужен торт на свадьбу 100 гостей");
  await page.locator('button[type="submit"]').click();
  await expect(page.locator("main div.w-64").first()).toBeVisible({
    timeout: 20000,
  });
  const firstBatch = await page.locator(nameSel).allTextContents();
  expect(firstBatch.length).toBeGreaterThan(0);

  // follow-up → "other options", a new batch (text varies by planner, so we
  // assert structurally: a fresh vendor appears that wasn't in the first batch)
  await ta.fill("не нравится, покажи другие варианты");
  await page.locator('button[type="submit"]').click();

  // at least one vendor in the latest batch was not in the first one
  await expect
    .poll(
      async () => {
        const all = await page.locator(nameSel).allTextContents();
        return all.some((n) => !firstBatch.includes(n));
      },
      { timeout: 15000 }
    )
    .toBe(true);
});
