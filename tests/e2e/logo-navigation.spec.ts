import { test, expect } from '@playwright/test';

test.describe('Logo Mobile Navigation', () => {
  test('should navigate to home page when clicking the mobile logo', async ({ page }) => {
    // 1. Start on the about page
    await page.goto('/about');

    // 2. Ensure mobile viewport
    await page.setViewportSize({ width: 375, height: 812 });

    // 3. Locate the mobile logo link
    // The mobile logo is visible on mobile viewports
    const mobileLogo = page.locator('header a[href="/"]').filter({ hasText: 'find what you want' }).first();

    // 4. Assert it is visible and clickable
    await expect(mobileLogo).toBeVisible();
    await mobileLogo.click();

    // 5. Verify we are now on the home page (or base URL)
    await expect(page).toHaveURL('/');
  });
});
