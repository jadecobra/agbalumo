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
    await searchInput.press('Enter');

    // Wait for the View All Listings button to be visible
    const viewAllLink = page.getByTestId('view-all-listings-link');
    await expect(viewAllLink).toBeVisible({ timeout: 10000 });

    // 3. Click "View All Listings" and intercept request
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

    // 2. Click the Event category (which should have 0 listings by default)
    const eventCategoryBtn = page.getByTestId('ag-filter-category-Event');
    await expect(eventCategoryBtn).toBeVisible();
    await eventCategoryBtn.click();

    // Verify filterState type is Event
    let filterState = await page.evaluate(() => (window as any).filterState);
    expect(filterState.type).toBe('Event');

    // 3. Wait for the View All Listings button to be visible in the empty state
    const viewAllLink = page.getByTestId('view-all-listings-link');
    await expect(viewAllLink).toBeVisible({ timeout: 10000 });

    // 4. Click "View All Listings" and intercept request
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
  });
});
