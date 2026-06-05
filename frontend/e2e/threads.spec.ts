import { test, expect } from "@playwright/test";
import {
  seedApprovedVendor,
  apiSignup,
  apiCreateBooking,
  apiSetBookingStatus,
  plantToken,
  uniqueEmail,
} from "./helpers";

const BACKEND = "http://localhost:8080";

// Realtime DM: a thread opens automatically when the vendor accepts a
// booking; customer and vendor exchange messages (realtime service).
test("customer and vendor exchange messages in a booking thread", async ({
  browser,
}) => {
  const vendor = await seedApprovedVendor("MsgVendor");

  const custEmail = uniqueEmail("msgcust");
  const custToken = await apiSignup(custEmail, "Msg Customer", "customer");

  const booking = await apiCreateBooking(custToken, {
    vendorId: vendor.vendorId,
    serviceId: vendor.serviceId,
    eventDate: "2026-11-20",
    guestCount: 120,
    amount: 200000,
  });
  // accepting the booking triggers EnsureThread in core → realtime
  await apiSetBookingStatus(vendor.token, booking.id, "accepted");

  // resolve the thread id from the customer's thread list
  const listRes = await fetch(`${BACKEND}/api/threads`, {
    headers: { authorization: `Bearer ${custToken}` },
  });
  const list = (await listRes.json()) as {
    items: Array<{ thread: { id: string } }> | null;
  };
  const threadId = list.items?.[0]?.thread.id;
  expect(threadId, "thread should open on accept").toBeTruthy();

  // --- customer sends a message ---
  const custCtx = await browser.newContext();
  const custPage = await custCtx.newPage();
  await plantToken(custPage, custToken);
  await custPage.goto(`/threads/${threadId}`);
  await expect(custPage.getByText(/No messages yet/i)).toBeVisible({
    timeout: 10000,
  });
  const custInput = custPage.getByPlaceholder(/Type a message/i);
  await custInput.fill("Hi, is the date available?");
  await custInput.press("Enter");
  await expect(
    custPage.getByText("Hi, is the date available?")
  ).toBeVisible({ timeout: 10000 });

  // --- vendor sees it and replies ---
  const vendCtx = await browser.newContext();
  const vendPage = await vendCtx.newPage();
  await plantToken(vendPage, vendor.token);
  await vendPage.goto(`/threads/${threadId}`);
  await expect(
    vendPage.getByText("Hi, is the date available?")
  ).toBeVisible({ timeout: 10000 });
  const vendInput = vendPage.getByPlaceholder(/Type a message/i);
  await vendInput.fill("Yes, it is free — happy to host you.");
  await vendInput.press("Enter");
  await expect(
    vendPage.getByText("Yes, it is free — happy to host you.")
  ).toBeVisible({ timeout: 10000 });

  // --- customer sees the reply after reloading history ---
  await custPage.reload();
  await expect(
    custPage.getByText("Yes, it is free — happy to host you.")
  ).toBeVisible({ timeout: 10000 });

  await custCtx.close();
  await vendCtx.close();
});
