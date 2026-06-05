import { test, expect } from "@playwright/test";
import { loginAs, uniqueEmail } from "./helpers";

// Customer-only payment card management on /cards.
test("customer adds two cards, switches default, deletes one", async ({
  page,
}) => {
  await loginAs(page, uniqueEmail("cards"), "password123", "Cards User", "customer");
  await page.goto("/cards");

  // empty state
  await expect(page.getByText(/No cards yet/i)).toBeVisible({ timeout: 10000 });

  // add first card (auto-becomes default)
  await page.getByRole("button", { name: /Add card/i }).click();
  await page.getByPlaceholder("4242 4242 4242 4242").fill("4242 4242 4242 4242");
  await page.getByPlaceholder("MM/YY").fill("11/30");
  await page.getByPlaceholder("NAME ON CARD").fill("FIRST CARD");
  await page.getByRole("button", { name: /Save card/i }).click();

  await expect(page.getByText("•••• 4242")).toBeVisible({ timeout: 10000 });
  await expect(page.getByText(/^Default$/i)).toBeVisible();

  // add second card
  await page.getByRole("button", { name: /Add card/i }).click();
  await page.getByPlaceholder("4242 4242 4242 4242").fill("5555 5555 5555 4444");
  await page.getByPlaceholder("MM/YY").fill("01/29");
  await page.getByPlaceholder("NAME ON CARD").fill("SECOND CARD");
  await page.getByRole("button", { name: /Save card/i }).click();
  await expect(page.getByText("•••• 4444")).toBeVisible({ timeout: 10000 });

  // promote the second card to default
  const secondRow = page.locator("li", { hasText: "•••• 4444" });
  await secondRow.getByRole("button", { name: /Set default/i }).click();
  await expect(secondRow.getByText(/^Default$/i)).toBeVisible({ timeout: 10000 });

  // delete the first card
  const firstRow = page.locator("li", { hasText: "•••• 4242" });
  await firstRow.getByRole("button", { name: /^Delete$/i }).click();
  await expect(page.getByText("•••• 4242")).toHaveCount(0, { timeout: 10000 });
  await expect(page.getByText("•••• 4444")).toBeVisible();
});
