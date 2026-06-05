import { test, expect } from "@playwright/test";
import {
  seedApprovedVendor,
  apiSignup,
  apiCreateBooking,
  apiSetBookingStatus,
  apiAdminToken,
  plantToken,
  uniqueEmail,
} from "./helpers";

const BACKEND = "http://localhost:8080";

// Review lifecycle: completed booking → customer leaves review (UI) →
// review surfaces on the public vendor page → admin removes it.
test("customer reviews a completed booking, admin deletes it", async ({
  page,
}) => {
  const vendor = await seedApprovedVendor("RevVendor", { price: 150000 });

  const custEmail = uniqueEmail("revcust");
  const custToken = await apiSignup(custEmail, "Review Customer", "customer");

  const booking = await apiCreateBooking(custToken, {
    vendorId: vendor.vendorId,
    serviceId: vendor.serviceId,
    eventDate: "2026-10-01",
    guestCount: 50,
    amount: 150000,
  });
  // vendor accepts then marks completed (only completed/paid is reviewable)
  await apiSetBookingStatus(vendor.token, booking.id, "accepted");
  await apiSetBookingStatus(vendor.token, booking.id, "completed");

  // --- customer submits review via UI ---
  await plantToken(page, custToken);
  await page.goto("/bookings");

  await page.getByRole("button", { name: /^Reviews$/i }).click();
  await page.getByLabel("Rate 5").click();
  await page
    .getByPlaceholder(/Share your experience/i)
    .fill("Outstanding service, highly recommend!");
  await page.getByRole("button", { name: /Submit review/i }).click();
  // Note: the bookings page closes the form on submit (onSubmitted clears it),
  // so the inline "Thanks" never lingers — assert the persisted outcome instead.

  // --- review is visible on the public vendor page ---
  await page.goto(`/vendors/${vendor.vendorId}`);
  await expect(
    page.getByText(/Outstanding service, highly recommend/i)
  ).toBeVisible({ timeout: 10000 });

  // --- admin deletes the review ---
  const admin = await apiAdminToken();
  const reviewsRes = await fetch(
    `${BACKEND}/api/vendors/${vendor.vendorId}/reviews`
  );
  const reviews = (await reviewsRes.json()) as {
    items: Array<{ id: string }> | null;
  };
  const reviewId = reviews.items?.[0]?.id;
  expect(reviewId).toBeTruthy();

  const del = await fetch(`${BACKEND}/api/admin/reviews/${reviewId}`, {
    method: "DELETE",
    headers: { authorization: `Bearer ${admin}` },
  });
  expect(del.ok).toBeTruthy();

  await page.reload();
  await expect(
    page.getByText(/Outstanding service, highly recommend/i)
  ).toHaveCount(0, { timeout: 10000 });
});
