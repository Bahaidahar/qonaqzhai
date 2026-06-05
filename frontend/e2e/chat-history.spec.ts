import { test, expect } from "@playwright/test";
import { loginAs, uniqueEmail } from "./helpers";

// AI chat is server-persisted: a sent message creates a chat that shows up in
// the sidebar "Recent" list and survives navigation / reload.
test("a sent chat is saved and appears in recent chats", async ({ page }) => {
  await loginAs(page, uniqueEmail("hist"), "password123", "History", "customer");
  await page.goto("/");

  const marker = `e2e recents ${Date.now()}`;
  const ta = page.locator("textarea").first();
  await expect(ta).toBeVisible({ timeout: 10000 });
  await ta.fill(marker);
  await page.locator('button[type="submit"]').click();

  // AI reply block renders → the turn was persisted
  await expect(
    page.getByText(/event_plan|budget|vendors|той_жоспары/i).first()
  ).toBeVisible({ timeout: 20000 });

  // sidebar "Recent" now lists this chat
  const sidebar = page.locator("aside");
  await expect(sidebar.getByText(marker)).toBeVisible({ timeout: 10000 });

  // start a fresh chat — transcript clears but the recent entry stays
  await page.getByRole("button", { name: /New chat/i }).click();
  await expect(sidebar.getByText(marker)).toBeVisible({ timeout: 10000 });

  // reopening it restores the conversation
  await sidebar.getByText(marker).click();
  await expect(page.locator("main").getByText(marker)).toBeVisible({
    timeout: 10000,
  });

  // survives a full reload (server-backed, not just in-memory)
  await page.reload();
  await expect(page.locator("aside").getByText(marker)).toBeVisible({
    timeout: 10000,
  });
});
