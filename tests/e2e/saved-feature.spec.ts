import { test, expect } from '@playwright/test';

test.describe('Saved/Favorites Feature', () => {
  // Use serial mode because tests share the same dev user and state
  test.describe.configure({ mode: 'serial' });

  test.beforeEach(async ({ page }) => {
    // Log browser messages for debugging
    page.on('console', msg => {
        if (msg.type() === 'error') {
            console.log(`[BROWSER ERROR] ${msg.text()}`);
        }
    });
  });

  test('heart icon visible on listing cards when authenticated', async ({ page }, testInfo) => {
    // Only run on Desktop project (hearts are hidden on mobile for authenticated users in some layouts)
    test.skip(testInfo.project.name === 'Mobile', 'Hearts on cards are hidden or positioned differently on mobile');

    // Log in as dev user
    const response = await page.goto('/auth/dev');
    expect(response?.status()).toBe(200);
    await page.waitForURL('/');
    
    // Ensure listings exist
    const listings = page.getByTestId('ag-listing-card');
    // Wait for the container to ensure page load/HTMX initialization
    await page.locator('#listings-container').waitFor({ state: 'attached', timeout: 10000 });
    const count = await listings.count();
    test.skip(count === 0, 'No listings in dev DB');

    // Verify heart button appears on listing cards
    const heart = page.getByTestId('ag-save-btn').first();
    await expect(heart).toBeVisible({ timeout: 10000 });
  });

  test('heart icon NOT visible when anonymous', async ({ page }) => {
    await page.goto('/');
    // Wait for JS to initialize
    await page.waitForFunction(() => typeof (window as any).filterState !== 'undefined', { timeout: 10000 });
    
    const heart = page.getByTestId('ag-save-btn');
    await expect(heart).toHaveCount(0);
  });

  test('heart toggles on click', async ({ page }, testInfo) => {
    test.skip(testInfo.project.name === 'Mobile', 'Hearts on cards are hidden or positioned differently on mobile');

    // Auth
    const response = await page.goto('/auth/dev');
    expect(response?.status()).toBe(200);
    
    // Ensure listings exist
    const listings = page.getByTestId('ag-listing-card');
    await page.locator('#listings-container').waitFor({ state: 'attached', timeout: 10000 });
    const count = await listings.count();
    test.skip(count === 0, 'No listings in dev DB');

    const firstListing = listings.first();
    const heart = firstListing.getByTestId('ag-save-btn');
    await expect(heart).toBeVisible();

    // Check if already saved and unsave first if so, to have a clean start
    const isSaved = await heart.evaluate(el => el.classList.contains('text-red-500'));
    if (isSaved) {
        const responsePromise = page.waitForResponse(res => 
            res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
        );
        await heart.click();
        await responsePromise;
        await expect(firstListing.getByTestId('ag-save-btn')).toHaveClass(/text-stone-400/);
    }
    
    // Click heart to save
    const heartToSave = firstListing.getByTestId('ag-save-btn');
    const savePromise = page.waitForResponse(res => 
        res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
    );
    await heartToSave.click();
    await savePromise;
    
    // Wait for HTMX swap and assert class changes (text-red-500 for saved)
    await expect(firstListing.getByTestId('ag-save-btn')).toHaveClass(/text-red-500/);

    // Click again to unsave
    const unsavePromise = page.waitForResponse(res => 
        res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
    );
    await firstListing.getByTestId('ag-save-btn').click();
    await unsavePromise;
    
    // Assert toggled back (text-stone-400 for unsaved)
    await expect(firstListing.getByTestId('ag-save-btn')).toHaveClass(/text-stone-400/);
  });

  test('saved nav button filters to saved listings', async ({ page }, testInfo) => {
    const isMobile = testInfo.project.name === 'Mobile';
    
    // Auth
    const response = await page.goto('/auth/dev');
    expect(response?.status()).toBe(200);

    // Ensure listings exist
    const listings = page.getByTestId('ag-listing-card');
    await page.locator('#listings-container').waitFor({ state: 'attached', timeout: 10000 });
    const count = await listings.count();
    test.skip(count === 0, 'No listings in dev DB');

    // Save a listing
    const firstListing = listings.first();
    
    // Parse listing ID from hx-get in overlay
    const overlay = firstListing.locator('div[hx-get^="/listings/"]').first();
    const hxGet = await overlay.getAttribute('hx-get');
    const listingId = hxGet?.split('/').pop();
    
    test.skip(!listingId, 'Could not determine listing ID');

    test.skip(isMobile, 'Saving from card not supported on mobile in this test');

    const heart = firstListing.getByTestId('ag-save-btn');
    const isSaved = await heart.evaluate(el => el.classList.contains('text-red-500'));
    if (!isSaved) {
        const savePromise = page.waitForResponse(res => 
            res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
        );
        await heart.click();
        await savePromise;
        await expect(firstListing.getByTestId('ag-save-btn')).toHaveClass(/text-red-500/);
    }

    // Click nav heart
    const navBtn = page.getByTestId(isMobile ? 'ag-nav-saved-btn-mobile' : 'ag-nav-saved-btn');
    await expect(navBtn).toBeVisible();
    
    const responsePromise = page.waitForResponse(res => 
      res.url().includes('/saved') && res.status() === 200
    );
    await navBtn.click();
    await responsePromise;

    // Assert listings container updates via HTMX and shows the saved listing
    // The ID in the DOM is listing-{uuid}
    await expect(page.locator(`#listing-${listingId}`)).toBeVisible();
  });
});
