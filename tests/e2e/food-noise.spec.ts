import { test, expect } from '@playwright/test';

test.describe('Noise Removal Verification for Food Listings', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForFunction(() => typeof (window as any).filterState !== 'undefined', { timeout: 10000 });
  });

  test('Food Listings must hide ContactEmail', async ({ page }) => {
    // 1. Filter by Food to get a Food listing
    const isMobile = await page.evaluate(() => window.innerWidth < 768);
    const toggleTestId = isMobile ? 'ag-home-filters-toggle-mobile' : 'ag-home-filters-toggle-desktop';
    await page.getByTestId(toggleTestId).click();
    
    const panel = page.locator('#filter-dropdown-panel');
    await expect(panel).toBeVisible();

    const foodCategoryBtn = page.getByTestId('ag-filter-category-food');
    await expect(foodCategoryBtn).toBeVisible();
    await foodCategoryBtn.click();

    // Wait for the htmx load
    await page.waitForTimeout(2000);

    // 2. Select a Food listing card for Mama Put Dallas
    const foodCard = page.getByTestId('ag-listing-card').filter({ hasText: 'Mama Put Dallas' }).first();
    await expect(foodCard).toBeVisible();

    // 3. Open detail modal
    const overlay = foodCard.locator('div.absolute.inset-0').first();
    await overlay.click({ force: true });

    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible({ timeout: 15000 });

    // 4. Assert ContactEmail is ABSENT
    const emailLink = modal.locator('[data-ada-discovery="email"]');
    await expect(emailLink).not.toBeAttached();
  });

  test('Business Listings must keep ContactEmail', async ({ page }) => {
    // 1. Filter by Business
    const isMobile = await page.evaluate(() => window.innerWidth < 768);
    const toggleTestId = isMobile ? 'ag-home-filters-toggle-mobile' : 'ag-home-filters-toggle-desktop';
    await page.getByTestId(toggleTestId).click();

    // Since Business is a dynamic category, Test ID uses the capitalized name
    const businessCategoryBtn = page.getByTestId('ag-filter-category-Business');
    await expect(businessCategoryBtn).toBeVisible();
    await businessCategoryBtn.click();

    // Wait for the htmx load
    await page.waitForTimeout(2000);

    // 2. Select a Business listing card
    const businessCard = page.getByTestId('ag-listing-card').first();
    await expect(businessCard).toBeVisible();

    // 3. Open detail modal
    const overlay = businessCard.locator('div.absolute.inset-0').first();
    await overlay.click({ force: true });

    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible({ timeout: 15000 });

    // 4. Assert ContactEmail is PRESENT
    const emailLink = modal.locator('[data-ada-discovery="email"]');
    await expect(emailLink).toBeAttached();
  });
});
