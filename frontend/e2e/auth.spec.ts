import { test, expect } from "@playwright/test";
import { uniqueEmail } from "./helpers";

test("unauthenticated visitor sees auth screen", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByPlaceholder("you@example.com")).toBeVisible();
});

test("customer can sign up via UI", async ({ page }) => {
  await page.goto("/");
  await page
    .getByRole("button", { name: /Sign up|Регистрация|Тіркелу/ })
    .first()
    .click();
  // pick customer role
  await page
    .locator('button:has-text("Plan an event"), button:has-text("Планировать"), button:has-text("Жоспарлау")')
    .first()
    .click();

  const email = uniqueEmail("cust");
  await page.getByPlaceholder("Aigerim").fill("E2E Customer");
  await page.getByPlaceholder("you@example.com").fill(email);
  await page.getByPlaceholder("••••••••").fill("password123");
  await page
    .getByRole("button", { name: /Create account|Создать|Аккаунт жасау/ })
    .click();

  // landed in chat — sidebar visible with user
  await expect(page.getByText("E2E Customer").first()).toBeVisible({
    timeout: 10000,
  });
});

test("vendor signup redirects to vendor dashboard", async ({ page }) => {
  await page.goto("/");
  await page
    .getByRole("button", { name: /Sign up|Регистрация|Тіркелу/ })
    .first()
    .click();
  await page
    .locator('button:has-text("Offer services"), button:has-text("Оказывать"), button:has-text("Қызмет")')
    .first()
    .click();

  const email = uniqueEmail("vend");
  await page.getByPlaceholder("Aigerim").fill("E2E Vendor");
  await page.getByPlaceholder("you@example.com").fill(email);
  await page.getByPlaceholder("••••••••").fill("password123");
  await page
    .getByRole("button", { name: /Create account|Создать|Аккаунт жасау/ })
    .click();

  await expect(page).toHaveURL(/\/vendor/, { timeout: 10000 });
});

test("bad credentials show error", async ({ page }) => {
  await page.goto("/");
  await page.getByPlaceholder("you@example.com").fill("nope@example.com");
  await page.getByPlaceholder("••••••••").fill("password123");
  await page.getByRole("button", { name: /^Sign in$|^Войти$|^Кіру$/ }).click();
  await expect(page.getByText(/invalid credentials/i)).toBeVisible();
});
