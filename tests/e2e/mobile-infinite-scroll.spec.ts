import { test, expect } from '@playwright/test';

test.describe('Mobile Infinite Scroll & Desktop Pagination', () => {
  test('Mobile viewport: should hide desktop pagination, display loading sentinel, and load next page when sentinel is revealed', async ({ page }) => {
    // 1. Emulate Mobile Viewport (iPhone 11 style)
    await page.setViewportSize({ width: 375, height: 812 });

    // Track the HTMX request for page 2 before navigation, since it triggers immediately on load
    const requestPromise = page.waitForRequest(request => 
      request.url().includes('/listings/fragment') && request.url().includes('page=2'),
      { timeout: 20000 }
    );

    // 2. Open Home page
    await page.goto('/?type=All&limit=3');
    await page.waitForLoadState('networkidle');

    // 3. Wait for the fragment request to complete
    await requestPromise;

    // 4. Desktop pagination controls should be hidden
    const pagination = page.locator('#pagination');
    await expect(pagination).toBeHidden();

    // 5. Verify initial cards exist
    const listingsContainer = page.locator('#listings-container');
    await expect(listingsContainer).toBeVisible();
    
    // We expect listing cards from both page 1 and page 2 (3 + 3 = 6 cards)
    const cards = listingsContainer.locator('[data-testid="ag-listing-card"]');
    const count = await cards.count();
    expect(count).toBe(6);

    // 6. Sentinel should exist and be active for page 3
    const sentinel = page.locator('#infinite-scroll-sentinel');
    await expect(sentinel).toBeVisible();
  });

  test('Desktop viewport: should display traditional pagination, hide infinite scroll sentinel, and paginate normally', async ({ page }) => {
    // 1. Set Desktop Viewport
    await page.setViewportSize({ width: 1440, height: 900 });

    // 2. Open Home page
    await page.goto('/?type=All&limit=3');
    await page.waitForLoadState('networkidle');

    // 3. Desktop pagination controls should be visible
    const pagination = page.locator('#pagination');
    await expect(pagination).toBeVisible();

    // 4. Infinite scroll sentinel should be hidden
    const sentinel = page.locator('#infinite-scroll-sentinel');
    await expect(sentinel).toBeHidden();

    // 5. Click the page 2 button in pagination and verify content reload
    const page2Button = pagination.locator('a', { hasText: '2' }).first();
    await expect(page2Button).toBeVisible();

    const requestPromise = page.waitForRequest(request => 
      request.url().includes('/listings/fragment') && request.url().includes('page=2'),
      { timeout: 20000 }
    );

    await page2Button.click();
    await requestPromise;

    // Verify that the listings container remains visible and reloaded
    await expect(page.locator('#listings-container')).toBeVisible();
  });
});
