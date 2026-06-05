import { test, expect } from "@playwright/test";
import { seedApprovedVendor, loginAs, uniqueEmail } from "./helpers";

// Public catalog browse on /vendors: text search, price ceiling, category
// pill filter. Uses a freshly seeded, uniquely-named Catering vendor.
test("customer filters the catalog by search, price and category", async ({
  page,
}) => {
  const vendor = await seedApprovedVendor("CatVendor", {
    category: "Catering",
    price: 333000,
  });

  await loginAs(page, uniqueEmail("browse"), "password123", "Browser", "customer");
  await page.goto("/vendors");

  const search = page.getByPlaceholder(/Search/i);
  const card = page.getByText(vendor.name);

  // text search finds the vendor
  await search.fill(vendor.name);
  await expect(card).toBeVisible({ timeout: 10000 });

  // price ceiling below its priceFrom hides it, raising it brings it back
  const priceMax = page.getByPlaceholder(/Max price/i);
  await priceMax.fill("100000");
  await expect(card).toHaveCount(0, { timeout: 10000 });
  await priceMax.fill("");
  await expect(card).toBeVisible({ timeout: 10000 });

  // category pills: a different category excludes it, its own includes it
  await search.fill("");
  await page.getByRole("button", { name: "Venue", exact: true }).click();
  await expect(card).toHaveCount(0, { timeout: 10000 });
  await page.getByRole("button", { name: "Catering", exact: true }).click();
  await expect(card).toBeVisible({ timeout: 10000 });
});
