import { test, expect } from '@playwright/test';

test.describe('Reset Filters UX', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForFunction(() => typeof (window as any).filterState !== 'undefined', { timeout: 10000 });
  });

  test('should reset all filters when clicking "VIEW ALL LISTINGS"', async ({ page }) => {
    // 1. Open filters panel and click a category (e.g. Food)
    const toggle = page.getByTestId('ag-home-filters-toggle-desktop');
    await expect(toggle).toBeVisible();
    await toggle.click();

    const panel = page.locator('#filter-dropdown-panel');
    await expect(panel).toBeVisible();

    const foodCategoryBtn = page.getByTestId('ag-filter-category-food');
    await expect(foodCategoryBtn).toBeVisible();
    await foodCategoryBtn.click();

    // Verify state type is Food
    let filterState = await page.evaluate(() => (window as any).filterState);
    expect(filterState.type).toBe('Food');

    // 2. Type a query that yields no results to show the "View All Listings" button
    const searchInput = page.getByTestId('ag-home-search-input');
    await expect(searchInput).toBeVisible();
    await searchInput.focus();
    
    // We'll type something completely random that has 0 results
    await searchInput.fill('xyzabc123nonexistent');
    await Promise.all([
      page.waitForResponse(resp => resp.url().includes('/listings/fragment') && resp.status() === 200),
      searchInput.press('Enter')
    ]);

    // Wait for the View All Listings button to be visible
    const viewAllLink = page.getByTestId('view-all-listings-link');
    await expect(viewAllLink).toBeVisible({ timeout: 10000 });

    // 3. Click "View All Listings" and intercept request
    await page.waitForTimeout(500);
    const [request] = await Promise.all([
      page.waitForRequest(req => req.url().includes('/listings/fragment')),
      viewAllLink.click()
    ]);
    const url = new URL(request.url());
    expect(url.searchParams.get('type')).toBe('All');

    // 4. Assert filterState has been reset
    filterState = await page.evaluate(() => (window as any).filterState);
    expect(filterState.type).toBe('All');
    expect(filterState.city).toBe('');
    expect(filterState.radius).toBe('25');

    // 5. Assert UI inputs are cleared
    const searchInputValue = await searchInput.inputValue();
    expect(searchInputValue).toBe('');

    const cityInput = page.locator('#filter-city');
    const cityInputValue = await cityInput.inputValue();
    expect(cityInputValue).toBe('');
  });

  test('should reset all filters when clicking "VIEW ALL LISTINGS" on an empty category (e.g. Event)', async ({ page }) => {
    // 1. Open filters panel
    const toggle = page.getByTestId('ag-home-filters-toggle-desktop');
    await expect(toggle).toBeVisible();
    await toggle.click();

    const panel = page.locator('#filter-dropdown-panel');
    await expect(panel).toBeVisible();

    // 2. Click the Event category (which should have 0 listings by default)
    const eventCategoryBtn = page.getByTestId('ag-filter-category-Event');
    await expect(eventCategoryBtn).toBeVisible();
    await Promise.all([
      page.waitForResponse(resp => resp.url().includes('/listings/fragment') && resp.status() === 200),
      eventCategoryBtn.click()
    ]);

    // Verify filterState type is Event
    let filterState = await page.evaluate(() => (window as any).filterState);
    expect(filterState.type).toBe('Event');

    // 3. Wait for the View All Listings button to be visible in the empty state
    const viewAllLink = page.getByTestId('view-all-listings-link');
    await expect(viewAllLink).toBeVisible({ timeout: 10000 });

    // 4. Click "View All Listings" and intercept request
    await page.waitForTimeout(500);
    const [request] = await Promise.all([
      page.waitForRequest(req => req.url().includes('/listings/fragment')),
      viewAllLink.click()
    ]);

    const url = new URL(request.url());
    expect(url.searchParams.get('type')).toBe('All');

    // 5. Assert filterState has been reset
    filterState = await page.evaluate(() => (window as any).filterState);
    expect(filterState.type).toBe('All');
    expect(filterState.city).toBe('');
    expect(filterState.radius).toBe('25');

    // 6. Assert UI inputs are cleared
    const searchInput = page.getByTestId('ag-home-search-input');
    const searchInputValue = await searchInput.inputValue();
    expect(searchInputValue).toBe('');

    const cityInput = page.locator('#filter-city');
    const cityInputValue = await cityInput.inputValue();
    expect(cityInputValue).toBe('');

    // 7. Stronger post-swap assertions capturing the exact reported symptom:
    // After reset from a 0-result category (Event/Job), results must be populated
    // AND the category UI must visually reset (All active, Event not active).
    // Wait for the HTMX fragment response + swap to settle.
    await page.waitForResponse(
      (resp) => resp.url().includes('/listings/fragment') && resp.status() === 200,
      { timeout: 15000 }
    );

    // No empty state (the core "still shows 0 listings" failure).
    const noResultsMsg = page.getByText(/No listings found|No saved listings yet/i);
    await expect(noResultsMsg).not.toBeVisible({ timeout: 10000 });

    // Container has real content after the swap.
    const listingsContainer = page.locator('#listings-container');
    await expect(listingsContainer).not.toBeEmpty();

    // Category buttons visually reset (covers "does not reset all filters/categories").
    const allCat = page.getByTestId('ag-filter-category-all');
    await expect(allCat).toHaveClass(/(^|\s)bg-earth-ochre\/10/);

    const eventCat = page.getByTestId('ag-filter-category-Event');
    await expect(eventCat).not.toHaveClass(/(^|\s)bg-earth-ochre\/10/);
  });

  test('should clear geolocation coordinates and reset "Near Me" button when clicking "RESET ALL FILTERS"', async ({ page }) => {
    // 1. Start with already located state in sessionStorage before page navigation
    await page.addInitScript(() => {
      sessionStorage.setItem('agbalumo_lat', '6.5244');
      sessionStorage.setItem('agbalumo_lng', '3.3792');
    });

    await page.goto('/');

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toContainText('Nearby');

    // 2. Type a query that yields no results to show the reset button
    const searchInput = page.getByTestId('ag-home-search-input');
    await expect(searchInput).toBeVisible();
    await searchInput.focus();
    await searchInput.fill('xyzabc123nonexistent');
    await Promise.all([
      page.waitForResponse(resp => resp.url().includes('/listings/fragment') && resp.status() === 200),
      searchInput.press('Enter')
    ]);

    // Wait for the Reset All Filters button to be visible
    const viewAllLink = page.getByTestId('view-all-listings-link');
    await expect(viewAllLink).toBeVisible({ timeout: 10000 });

    // 3. Click "Reset All Filters" and intercept request
    await page.waitForTimeout(500);
    const [request] = await Promise.all([
      page.waitForRequest(req => req.url().includes('/listings/fragment')),
      viewAllLink.click()
    ]);

    const url = new URL(request.url());
    expect(url.searchParams.get('type')).toBe('All');
    expect(url.searchParams.get('lat')).toBeNull();
    expect(url.searchParams.get('lng')).toBeNull();

    // 4. Assert sessionStorage coordinates are cleared
    const sessionStorageCoords = await page.evaluate(() => {
      return {
        lat: sessionStorage.getItem('agbalumo_lat'),
        lng: sessionStorage.getItem('agbalumo_lng')
      };
    });
    expect(sessionStorageCoords.lat).toBeNull();
    expect(sessionStorageCoords.lng).toBeNull();

    // 5. Assert Near Me button is reset back to 'Near Me'
    await expect(nearMeBtn).toContainText('Near Me');
  });
});
