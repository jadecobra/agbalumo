import { test, expect } from '@playwright/test';

test.describe('Sandbox Visual Regression', () => {
  test('capture and diff the entire sandbox layout', async ({ page }) => {
    await page.goto('/sandbox');
    // Wait for the sandbox header to ensure the page has loaded
    await expect(page.locator('h1')).toBeVisible();
    
    // Capture and diff the entire sandbox layout
    await expect(page).toHaveScreenshot('sandbox-baseline.png', {
      fullPage: true,
    });
  });
});
