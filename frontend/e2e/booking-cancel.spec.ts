import { test, expect } from "@playwright/test";
import {
  approveVendorByEmail,
  loginAs,
  uniqueEmail,
} from "./helpers";

test("customer creates booking then cancels via UI", async ({ browser }) => {
  const vendEmail = uniqueEmail("vend_cancel");
  const custEmail = uniqueEmail("cust_cancel");

  // setup approved vendor via API
  const vendCtx = await browser.newContext();
  const vendPage = await vendCtx.newPage();
  await loginAs(vendPage, vendEmail, "password123", "CV", "vendor");
  await vendPage.goto("/vendor");
  await vendPage.getByPlaceholder(/Rixos Almaty Ballroom/i).fill("Cancel Studio");
  await vendPage.getByRole("button", { name: /^Save$/i }).click();
  await expect(vendPage.getByText(/Pending/i)).toBeVisible({ timeout: 10000 });
  await approveVendorByEmail(vendEmail);
  await vendCtx.close();

  // customer books
  const custCtx = await browser.newContext();
  const custPage = await custCtx.newPage();
  await loginAs(custPage, custEmail, "password123", "CC", "customer");
  await custPage.goto("/vendors");
  await expect(custPage.getByText("Cancel Studio").first()).toBeVisible({
    timeout: 10000,
  });
  await custPage.getByText("Cancel Studio").first().click();

  await custPage.locator('input[type="date"]').fill("2026-10-01");
  await custPage.getByPlaceholder("150").fill("50");
  await custPage.getByRole("button", { name: /Request booking/i }).click();
  await expect(custPage.getByText(/Request sent/i)).toBeVisible({
    timeout: 10000,
  });

  // bookings page
  await custPage.goto("/bookings");
  await expect(custPage.getByText("2026-10-01")).toBeVisible({ timeout: 10000 });
  await custPage.getByRole("button", { name: /Cancel/i }).click();
  await expect(custPage.getByText("cancelled").first()).toBeVisible({
    timeout: 10000,
  });

  await custCtx.close();
});
