import { test, expect } from "@playwright/test";
import {
  seedApprovedVendor,
  loginAs,
  apiMyBookings,
  uniqueEmail,
} from "./helpers";

// A vendor card in chat can book the vendor directly. When the prompt carries
// a date the booking is one click; otherwise an inline date/guests form opens.
// (Many decor vendors may match, so we act on the first card rather than a
// specific seeded one, and verify the booking landed in the backend.)

test("chat card books a vendor directly when the prompt has a date", async ({
  page,
}) => {
  await seedApprovedVendor("BookDecor", { category: "Decor & Florists" });
  const token = await loginAs(
    page,
    uniqueEmail("bookc"),
    "password123",
    "Booker",
    "customer"
  );
  await page.goto("/");

  const ta = page.locator("textarea").first();
  await ta.fill("декор на 80 гостей 15 сентября");
  await page.locator('button[type="submit"]').click();

  const card = page.locator("div.w-64").first();
  await expect(card).toBeVisible({ timeout: 20000 });
  await card.getByRole("button", { name: /Book now/i }).click();

  await expect(card.getByText(/Request sent/i)).toBeVisible({ timeout: 10000 });

  // backend has a pending booking for this customer
  const bookings = await apiMyBookings(token);
  expect(bookings.items?.length ?? 0).toBeGreaterThan(0);
  expect(bookings.items?.[0]?.status).toBe("pending");
});

test("chat card opens a date form when the prompt has no date", async ({
  page,
}) => {
  await seedApprovedVendor("FormDecor", { category: "Decor & Florists" });
  const token = await loginAs(
    page,
    uniqueEmail("formc"),
    "password123",
    "Former",
    "customer"
  );
  await page.goto("/");

  const ta = page.locator("textarea").first();
  await ta.fill("нужен декор для свадьбы");
  await page.locator('button[type="submit"]').click();

  const card = page.locator("div.w-64").first();
  await expect(card).toBeVisible({ timeout: 20000 });
  await card.getByRole("button", { name: /Book now/i }).click();

  // inline form appears; fill date + guests then confirm
  const dateInput = card.locator('input[type="date"]');
  await expect(dateInput).toBeVisible({ timeout: 10000 });
  await dateInput.fill("2026-09-20");
  await card.locator('input[type="number"]').fill("60");
  await card.getByRole("button", { name: /Book now/i }).click();

  await expect(card.getByText(/Request sent/i)).toBeVisible({ timeout: 10000 });

  const bookings = await apiMyBookings(token);
  expect(bookings.items?.length ?? 0).toBeGreaterThan(0);
});
