import { test, expect } from '@playwright/test';

test.describe('Sandbox Visual Regression', () => {
  test('capture and diff the entire sandbox layout', async ({ page }) => {
    await page.goto('/sandbox');
    // Wait for the sandbox header to ensure the page has loaded
    await expect(page.locator('h1')).toBeVisible();
    
    // Capture and diff the entire sandbox layout
    await expect(page).toHaveScreenshot('sandbox-baseline.png', {
      fullPage: true,
      timeout: 30000,
    });
  });
});

test.describe('Listing Modal Visual Regression', () => {
  test('capture and diff the listing detail modal', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByTestId('ag-listing-card').first()).toBeVisible();
    await page.getByTestId('ag-listing-card').first().click();
    await expect(page.locator('dialog[id^="detail-modal-"]')).toBeVisible();
    await page.waitForTimeout(500);
    await expect(page).toHaveScreenshot('listing-modal-baseline.png', {
      fullPage: false,
      maxDiffPixelRatio: 0.05,
      timeout: 30000,
    });
  });
});
