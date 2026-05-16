import { test, expect } from "@playwright/test";
import { loginAs, uniqueEmail } from "./helpers";
import path from "node:path";
import fs from "node:fs";
import os from "node:os";

// 1x1 PNG bytes
const PNG_B64 =
  "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII=";

function writeTempPng(name: string): string {
  const p = path.join(os.tmpdir(), `${name}-${Date.now()}.png`);
  fs.writeFileSync(p, Buffer.from(PNG_B64, "base64"));
  return p;
}

test("vendor uploads multiple photos and deletes one", async ({ page }) => {
  await loginAs(page, uniqueEmail("photo"), "password123", "P", "vendor");
  await page.goto("/vendor");
  await expect(page.getByRole("heading", { name: /My profile/i })).toBeVisible();

  await page.getByPlaceholder(/Rixos Almaty Ballroom/i).fill("Photo Studio");
  await page.getByRole("button", { name: /^Save$/i }).click();
  await expect(page.getByText(/Pending/i)).toBeVisible({ timeout: 10000 });

  // upload photo 1
  const file1 = writeTempPng("p1");
  await page.locator('input[type="file"]').setInputFiles(file1);
  // upload photo 2
  await page.waitForTimeout(500);
  const file2 = writeTempPng("p2");
  await page.locator('input[type="file"]').setInputFiles(file2);
  await page.waitForTimeout(800);

  // both photos in gallery
  const imgs = page.locator('img[alt=""]');
  await expect(imgs).toHaveCount(2, { timeout: 10000 });

  // hover + delete one
  await page.locator('img[alt=""]').first().hover();
  await page.getByRole("button", { name: /Delete/i }).first().click();
  await expect(imgs).toHaveCount(1, { timeout: 5000 });
});
