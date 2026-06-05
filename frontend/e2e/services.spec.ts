import { test, expect } from "@playwright/test";
import { seedApprovedVendor, plantToken } from "./helpers";

// Vendor service management (ServicesManager on /vendor): create, toggle
// active state, and delete. Exercises /api/me/vendor/services CRUD.
test("vendor creates, toggles and deletes a service", async ({ page }) => {
  const vendor = await seedApprovedVendor("SvcVendor");

  // native confirm() on delete — auto-accept
  page.on("dialog", (d) => d.accept());

  await plantToken(page, vendor.token);
  await page.goto("/vendor");

  // seeded service is listed
  await expect(page.getByText("Standard package")).toBeVisible({
    timeout: 10000,
  });

  // create a new service
  await page.getByRole("button", { name: /New service/i }).click();
  const dialog = page.locator("div.fixed.inset-0.z-50");
  await dialog.getByPlaceholder(/Wedding photography/i).fill("Premium combo");
  await dialog.getByRole("spinbutton").fill("99000");
  await dialog.getByRole("button", { name: /^Save$/i }).click();

  const row = page.locator("li", { hasText: "Premium combo" });
  await expect(row).toBeVisible({ timeout: 10000 });

  // disable → inactive badge appears, then re-enable
  await row.getByTitle("Disable").click();
  await expect(row.getByText(/inactive/i)).toBeVisible({ timeout: 10000 });
  await row.getByTitle("Enable").click();
  await expect(row.getByText(/inactive/i)).toHaveCount(0, { timeout: 10000 });

  // delete it
  await row.getByRole("button").last().click();
  await expect(page.getByText("Premium combo")).toHaveCount(0, {
    timeout: 10000,
  });
});
