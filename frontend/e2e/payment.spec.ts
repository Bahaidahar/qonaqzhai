import { test, expect } from "@playwright/test";
import {
  seedApprovedVendor,
  apiSignup,
  apiCreateBooking,
  apiSetBookingStatus,
  apiAddCard,
  apiMyBookings,
  plantToken,
  uniqueEmail,
} from "./helpers";

// Full payment saga: booking accepted by vendor → customer pays (mock) →
// core marks booking paid (core → payment → core).
test("customer pays an accepted booking and it becomes paid", async ({
  page,
}) => {
  const vendor = await seedApprovedVendor("PayVendor", { price: 200000 });

  const custEmail = uniqueEmail("paycust");
  const custToken = await apiSignup(custEmail, "Pay Customer", "customer");
  await apiAddCard(custToken);

  const booking = await apiCreateBooking(custToken, {
    vendorId: vendor.vendorId,
    serviceId: vendor.serviceId,
    eventDate: "2026-09-15",
    guestCount: 80,
    amount: 200000,
  });
  await apiSetBookingStatus(vendor.token, booking.id, "accepted");

  // pay through the UI
  await plantToken(page, custToken);
  await page.goto("/bookings");

  const payBtn = page.getByRole("button", { name: /Pay/i });
  await expect(payBtn).toBeVisible({ timeout: 10000 });
  await payBtn.click();

  // the saga (core → payment → core) settles the booking to "paid".
  // Poll the backend so a slow gRPC round-trip under load doesn't flake.
  await expect
    .poll(
      async () => {
        const after = await apiMyBookings(custToken);
        return after.items?.find((b) => b.id === booking.id)?.status;
      },
      { timeout: 15000 }
    )
    .toBe("paid");

  // and the Pay action is gone from the UI once it's no longer "accepted"
  await expect(page.getByRole("button", { name: /Pay/i })).toHaveCount(0, {
    timeout: 10000,
  });
});
