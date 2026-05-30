import { test, expect } from '@playwright/test';

test.describe('Saved/Favorites Feature', () => {
  // Use serial mode because tests share the same dev user and state
  test.describe.configure({ mode: 'serial' });

  test.beforeEach(async ({ page }) => {
    // Log browser messages for debugging
    page.on('console', msg => {
        console.log(`[BROWSER LOG (${msg.type()})] ${msg.text()}`);
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
    await page.locator('#listings-container').waitFor({ state: 'attached', timeout: 15000 });
    const count = await listings.count();
    test.skip(count === 0, 'No listings in dev DB');

    // Verify heart button appears on listing cards
    console.log('=== VIEWPORT SIZE ===:', page.viewportSize());
    
    const card = page.locator('[data-testid="ag-listing-card"]').first();
    console.log('Card bounding box:', await card.boundingBox());
    console.log('Card computed display:', await card.evaluate(el => window.getComputedStyle(el).display));
    
    const btns = await page.locator('[data-testid="ag-save-btn"]').all();
    console.log('=== SAVE BUTTONS COUNT ===:', btns.length);
    for (let i = 0; i < btns.length; i++) {
        const btnStyle = await btns[i].evaluate(el => {
            const style = window.getComputedStyle(el);
            return { display: style.display, visibility: style.visibility, opacity: style.opacity, width: style.width, height: style.height };
        });
        console.log(`Save button ${i}: visible = ${await btns[i].isVisible()}, boundingBox = ${JSON.stringify(await btns[i].boundingBox())}, style = ${JSON.stringify(btnStyle)}, classes = ${await btns[i].getAttribute('class')}`);
    }

    const heart = page.locator('[data-testid="ag-save-btn"]:visible').first();
    await expect(heart).toBeVisible({ timeout: 15000 });
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
    const heart = firstListing.locator('[data-testid="ag-save-btn"]:visible').first();
    await expect(heart).toBeVisible();

    // Check if already saved and unsave first if so, to have a clean start
    const isSaved = await heart.evaluate(el => el.getAttribute('aria-label') === 'Unsave listing');
    if (isSaved) {
        const responsePromise = page.waitForResponse(res => 
            res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
        );
        await heart.click();
        await responsePromise;
        const unsavedHeart = firstListing.locator('[data-testid="ag-save-btn"]:visible').first();
        await expect(unsavedHeart).toHaveClass(/text-white/);
        await expect(unsavedHeart).not.toHaveClass(/text-earth-accent/);
    }
    
    // Click heart to save
    const heartToSave = firstListing.locator('[data-testid="ag-save-btn"]:visible').first();
    const savePromise = page.waitForResponse(res => 
        res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
    );
    await heartToSave.click();
    await savePromise;
    
    // Wait for HTMX swap and assert class changes (text-earth-accent for saved)
    await expect(firstListing.locator('[data-testid="ag-save-btn"]:visible').first()).toHaveClass(/text-earth-accent/);

    // Click again to unsave
    const unsavePromise = page.waitForResponse(res => 
        res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
    );
    await firstListing.locator('[data-testid="ag-save-btn"]:visible').first().click();
    await unsavePromise;
    
    // Assert toggled back (text-white for unsaved)
    const toggledHeart = firstListing.locator('[data-testid="ag-save-btn"]:visible').first();
    await expect(toggledHeart).toHaveClass(/text-white/);
    await expect(toggledHeart).not.toHaveClass(/text-earth-accent/);
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
    const overlay = firstListing.locator('button[hx-get^="/listings/"]').first();
    const hxGet = await overlay.getAttribute('hx-get');
    const listingId = hxGet?.split('/').pop();
    
    test.skip(!listingId, 'Could not determine listing ID');

    test.skip(isMobile, 'Saving from card not supported on mobile in this test');

    const heart = firstListing.locator('[data-testid="ag-save-btn"]:visible').first();
    const isSaved = await heart.evaluate(el => el.getAttribute('aria-label') === 'Unsave listing');
    if (!isSaved) {
        const savePromise = page.waitForResponse(res => 
            res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
        );
        await heart.click();
        await savePromise;
        await expect(firstListing.locator('[data-testid="ag-save-btn"]:visible').first()).toHaveClass(/text-earth-accent/);
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

  test('listing cards show heart button in saved listing view and clicking it removes it', async ({ page }, testInfo) => {
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
    const overlay = firstListing.locator('button[hx-get^="/listings/"]').first();
    const hxGet = await overlay.getAttribute('hx-get');
    const listingId = hxGet?.split('/').pop();
    
    test.skip(!listingId, 'Could not determine listing ID');
    test.skip(isMobile, 'Saving from card not supported on mobile in this test');

    const heart = firstListing.locator('[data-testid="ag-save-btn"]:visible').first();
    const isSaved = await heart.evaluate(el => el.getAttribute('aria-label') === 'Unsave listing');
    if (!isSaved) {
        const savePromise = page.waitForResponse(res => 
            res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
        );
        await heart.click();
        await savePromise;
        await expect(firstListing.locator('[data-testid="ag-save-btn"]:visible').first()).toHaveClass(/text-earth-accent/);
    }

    // Click nav heart
    const navBtn = page.getByTestId(isMobile ? 'ag-nav-saved-btn-mobile' : 'ag-nav-saved-btn');
    await expect(navBtn).toBeVisible();
    
    const savedListPromise = page.waitForResponse(res => 
      res.url().includes('/saved') && res.status() === 200
    );
    await navBtn.click();
    await savedListPromise;

    // Verify card is visible in saved view
    const savedCard = page.locator(`#listing-${listingId}`);
    await expect(savedCard).toBeVisible();

    // Verify heart button is visible on card in saved listing view and has the homepage background color when active
    const savedCardHeart = savedCard.locator('[data-testid="ag-save-btn"]:visible').first();
    await expect(savedCardHeart).toBeVisible();
    await expect(savedCardHeart).toHaveClass(/text-earth-accent/);

    // Click heart to remove from saved list
    const removePromise = page.waitForResponse(res => 
        res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
    );
    await savedCardHeart.click();
    await removePromise;

    // Verify it toggles to unsaved state (text-white)
    await expect(savedCardHeart).toHaveClass(/text-white/);
    await expect(savedCardHeart).not.toHaveClass(/text-earth-accent/);

    // Navigate back to saved view and verify the listing is no longer listed
    const refreshPromise = page.waitForResponse(res => 
      res.url().includes('/saved') && res.status() === 200
    );
    await page.goto('/saved');
    await refreshPromise;

    // The card should no longer be visible in the saved listings container
    await expect(savedCard).not.toBeVisible();
  });

  test('profile modal shows saved listings and clicking heart button removes listing', async ({ page }, testInfo) => {
    const isMobile = testInfo.project.name === 'Mobile';

    // Auth
    const response = await page.goto('/auth/dev');
    expect(response?.status()).toBe(200);

    // Ensure listings exist
    const listings = page.getByTestId('ag-listing-card');
    await page.locator('#listings-container').waitFor({ state: 'attached', timeout: 10000 });
    const count = await listings.count();
    test.skip(count === 0, 'No listings in dev DB');

    // Save a listing from the homepage first to ensure we have a saved listing
    const firstListing = listings.first();
    const overlay = firstListing.locator('button[hx-get^="/listings/"]').first();
    const hxGet = await overlay.getAttribute('hx-get');
    const listingId = hxGet?.split('/').pop();
    
    test.skip(!listingId, 'Could not determine listing ID');
    test.skip(isMobile, 'Saving from card not supported on mobile in this test');

    const heart = firstListing.locator('[data-testid="ag-save-btn"]:visible').first();
    const isSaved = await heart.evaluate(el => el.getAttribute('aria-label') === 'Unsave listing');
    if (!isSaved) {
        const savePromise = page.waitForResponse(res => 
            res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
        );
        await heart.click();
        await savePromise;
        await expect(firstListing.locator('[data-testid="ag-save-btn"]:visible').first()).toHaveClass(/text-earth-accent/);
    }

    // 1. Click on profile button to show profile modal
    const profileTrigger = page.locator('[data-testid="ag-nav-profile-btn"], [data-testid="mobile-account-btn"]')
      .filter({ visible: true })
      .first();
    await expect(profileTrigger).toBeVisible({ timeout: 10000 });
    await profileTrigger.click();

    // Verify profile modal appears
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible({ timeout: 10000 });
    await expect(modal.getByText('My Profile')).toBeVisible();

    // 2. Locate the listing card in Section 2: "Saved Listings" inside profile modal
    // It should have id `modal-saved-listing-{listingId}`
    const modalSavedCard = modal.locator(`#modal-saved-listing-${listingId}`);
    await expect(modalSavedCard).toBeVisible();

    // 3. Verify heart button is visible on that listing card in profile modal and matches the homepage background color when active
    const modalSavedCardHeart = modalSavedCard.locator('[data-testid="ag-save-btn"]:visible').first();
    await expect(modalSavedCardHeart).toBeVisible();
    await expect(modalSavedCardHeart).toHaveClass(/text-earth-accent/);

    // 4. Click heart button to remove the listing from the saved list
    const removePromise = page.waitForResponse(res => 
        res.url().includes('/save') && res.request().method() === 'POST' && res.status() === 200
    );
    await modalSavedCardHeart.click();
    await removePromise;

    // 5. Verify the heart button updates to unsaved state (text-white)
    await expect(modalSavedCardHeart).toHaveClass(/text-white/);
    await expect(modalSavedCardHeart).not.toHaveClass(/text-earth-accent/);

    // 6. Close the profile modal and reopen to verify the listing card is gone
    const closeBtn = modal.locator('button[aria-label="Close"]').first();
    await closeBtn.click();
    await expect(modal).not.toBeVisible();

    // Reopen profile modal
    await profileTrigger.click();
    await expect(modal).toBeVisible();

    // The card should no longer exist under "Saved Listings" section in profile modal
    await expect(modalSavedCard).not.toBeVisible();
  });
});

