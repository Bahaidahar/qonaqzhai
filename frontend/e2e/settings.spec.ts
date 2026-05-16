import { test, expect } from "@playwright/test";
import { loginAs, uniqueEmail } from "./helpers";

test("language switch updates UI", async ({ page }) => {
  await loginAs(page, uniqueEmail("l10n"), "password123", "L10n", "customer");
  await page.goto("/settings");

  // helpers force EN by default
  await expect(page.getByText(/^Chat$/)).toBeVisible();

  // switch to RU
  await page.locator('button[lang="ru"]').click();
  await expect(page.getByText(/^Чат$/).first()).toBeVisible();

  // KZ
  await page.locator('button[lang="kz"]').click();
  await expect(page.getByText("Мердігерлер")).toBeVisible();
});

test("theme switch applies data-theme attribute", async ({ page }) => {
  await loginAs(page, uniqueEmail("theme"), "password123", "T", "customer");
  await page.goto("/settings");

  // start light
  await page
    .locator('button[aria-pressed], button[aria-pressed="false"]')
    .filter({ hasText: /Light|Светлая|Ашық/ })
    .first()
    .click();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "light");

  // dark
  await page
    .locator("button")
    .filter({ hasText: /Dark|Тёмная|Қараңғы/ })
    .first()
    .click();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");

  // system — depends on env; just ensure attribute set to light or dark
  await page
    .locator("button")
    .filter({ hasText: /System|Системная|Жүйелік/ })
    .first()
    .click();
  await expect(page.locator("html")).toHaveAttribute("data-theme", /light|dark/);
});

test("theme persists across reload", async ({ page }) => {
  await loginAs(page, uniqueEmail("persist"), "password123", "P", "customer");
  await page.goto("/settings");
  await page
    .locator("button")
    .filter({ hasText: /^Dark$|Тёмная|Қараңғы/ })
    .first()
    .click();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
  await page.reload();
  await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
});

test("sign out clears state and shows auth screen", async ({ page }) => {
  await loginAs(page, uniqueEmail("out"), "password123", "Out", "customer");
  await page.goto("/settings");
  await page.getByRole("button", { name: /Sign out|Выйти|Шығу/ }).click();
  // back to auth screen
  await expect(page.getByPlaceholder("you@example.com")).toBeVisible({
    timeout: 5000,
  });
  // token removed
  const token = await page.evaluate(() =>
    window.localStorage.getItem("qonaqzhai_token")
  );
  expect(token).toBeNull();
});
