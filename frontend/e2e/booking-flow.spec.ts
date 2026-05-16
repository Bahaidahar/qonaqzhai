import { test, expect } from "@playwright/test";
import { approveVendorByEmail, loginAs, uniqueEmail } from "./helpers";

test("vendor profile → admin approve → customer book → vendor accepts", async ({
  browser,
}) => {
  const vendEmail = uniqueEmail("vend_flow");
  const custEmail = uniqueEmail("cust_flow");

  // --- vendor side ---
  const vendCtx = await browser.newContext();
  const vendPage = await vendCtx.newPage();
  await loginAs(vendPage, vendEmail, "password123", "Flow Vendor", "vendor");
  await vendPage.goto("/vendor");
  await expect(vendPage.getByRole("heading", { name: /My profile/i })).toBeVisible({
    timeout: 10000,
  });

  await vendPage
    .getByPlaceholder(/Rixos Almaty Ballroom/i)
    .fill("E2E Studio");
  await vendPage.getByPlaceholder("500000").fill("450000");
  await vendPage
    .getByPlaceholder(/Premier venue/i)
    .fill("E2E test vendor description.");
  await vendPage.getByRole("button", { name: /^Save$/i }).click();
  await expect(vendPage.getByText(/Pending admin approval/i)).toBeVisible({
    timeout: 10000,
  });

  // admin approves out of band
  await approveVendorByEmail(vendEmail);
  await vendPage.reload();
  await expect(vendPage.getByText(/^Approved/i)).toBeVisible({ timeout: 10000 });

  // --- customer side ---
  const custCtx = await browser.newContext();
  const custPage = await custCtx.newPage();
  await loginAs(custPage, custEmail, "password123", "Flow Customer", "customer");
  await custPage.goto("/vendors");

  await expect(custPage.getByText("E2E Studio").first()).toBeVisible({
    timeout: 10000,
  });
  await custPage.getByText("E2E Studio").first().click();

  await custPage.locator('input[type="date"]').fill("2026-08-12");
  await custPage.getByPlaceholder("150").fill("100");
  await custPage.getByRole("button", { name: /Request booking/i }).click();
  await expect(custPage.getByText(/Request sent/i)).toBeVisible({
    timeout: 10000,
  });

  // --- vendor accepts ---
  await vendPage.goto("/vendor/bookings");
  await expect(vendPage.getByText(/2026-08-12/)).toBeVisible({ timeout: 10000 });
  await vendPage.getByRole("button", { name: /^Accept$/i }).click();
  await expect(vendPage.getByText("accepted").first()).toBeVisible({
    timeout: 10000,
  });

  // customer sees status
  await custPage.goto("/bookings");
  await expect(custPage.getByText(/accepted/i).first()).toBeVisible({
    timeout: 10000,
  });

  await vendCtx.close();
  await custCtx.close();
});
