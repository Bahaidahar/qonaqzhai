import { test, expect } from "@playwright/test";

test("invalid token in storage → cleared, auth screen shown", async ({
  page,
}) => {
  await page.addInitScript(() => {
    window.localStorage.setItem("qonaqzhai_token", "garbage.invalid.token");
  });
  await page.goto("/");
  // app should detect invalid /me, clear token, show auth screen
  await expect(page.getByPlaceholder("you@example.com")).toBeVisible({
    timeout: 10000,
  });
  const token = await page.evaluate(() =>
    window.localStorage.getItem("qonaqzhai_token")
  );
  expect(token).toBeNull();
});

test("no token → auth screen", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByPlaceholder("you@example.com")).toBeVisible();
});
