import { test, expect } from '@playwright/test';

/**
 * Ada's UX Constraint: Discovery under 60s
 * This test simulates the exact real-world journey of a user (Ada)
 * searching for high-quality food to ensure it completes within the budget.
 */
test.describe('UX Constraint: Ada Journey', () => {
  // Enforce a strict test timeout to ensure the journey doesn't hang
  test.setTimeout(60000);

  test('Ada should find and view a listing in under 60 seconds', async ({ page }) => {
    const startTime = Date.now();

    // 1. Landing on the site
    // We wait for networkidle to ensure initial assets are loaded
    await page.goto('/', { waitUntil: 'networkidle' });
    
    // Ensure the app is initialized (filterState is our JS initialization signal)
    await page.waitForFunction(() => typeof (window as any).filterState !== 'undefined', { timeout: 10000 });

    // 2. Typing a search query with realistic delay
    const searchInput = page.getByTestId('ag-home-search-input');
    await expect(searchInput).toBeVisible();
    
    // Human-like typing speed
    await searchInput.focus();
    await searchInput.pressSequentially('Nigerian', { delay: 100 });
    await searchInput.press('Enter'); // Explicitly trigger search

    // 3. Pause to "read" results
    // Wait for the HTMX request to complete
    await page.waitForResponse(res => 
      res.url().includes('/listings/fragment') && res.status() === 200
    );
    
    // Simulate Ada reading the titles of the results
    await page.waitForTimeout(2000);

    // 4. Clicking a listing
    const listingCard = page.getByTestId('ag-listing-card').first();
    await expect(listingCard).toBeVisible({ timeout: 10000 });
    
    // Target the click overlay specifically for reliable interaction
    const overlay = listingCard.locator('div.absolute.inset-0').first();
    await overlay.scrollIntoViewIfNeeded();
    await overlay.click({ force: true });

    // 5. Verify the detail modal appears
    // The detail view is a <dialog>. We wait for it to be open.
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible({ timeout: 15000 });
    
    // Verify the listing title is present in the modal
    const modalTitle = modal.locator('h2.font-serif');
    await expect(modalTitle).toBeVisible();
    await expect(modalTitle).not.toBeEmpty();
    
    const endTime = Date.now();
    const totalDuration = endTime - startTime;
    
    console.log(`Ada Journey completed in: ${totalDuration}ms`);
    
    // Performance Assertion: Entire journey must be well under 60s
    expect(totalDuration).toBeLessThan(60000);
  });
});
