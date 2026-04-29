import { test, expect } from '@playwright/test';

test.describe('HTMX Interactions and State Sync', () => {
  test.beforeEach(async ({ page }) => {
    // Log browser messages for debugging in case of failures
    page.on('console', msg => {
        if (msg.type() === 'error') {
            console.log(`[BROWSER ${msg.type()}] ${msg.text()}`);
        }
    });
    page.on('pageerror', err => console.error(`[BROWSER ERROR] ${err.message}`));

    await page.goto('/');
    
    // Wait for JS to initialize (filterState is defined at the top level of filters.js)
    await page.waitForFunction(() => typeof (window as any).filterState !== 'undefined', { timeout: 10000 });
  });

  test('should update state and trigger HTMX request on search input', async ({ page }) => {
    const searchInput = page.getByTestId('ag-home-search-input');
    await expect(searchInput).toBeVisible();
    
    // Track requests
    const requestPromise = page.waitForRequest(request => 
      request.url().includes('/listings/fragment') && request.url().includes('q=jollof'),
      { timeout: 20000 }
    );

    await searchInput.focus();
    await searchInput.pressSequentially('jollof', { delay: 10 });
    
    await requestPromise;
    await expect(page.locator('#listings-container')).toBeVisible();
  });

  test('should update filterState and trigger HTMX request on category click', async ({ page }) => {
    const isMobile = await page.evaluate(() => window.innerWidth < 768);
    const toggleTestId = isMobile ? 'ag-home-filters-toggle-mobile' : 'ag-home-filters-toggle-desktop';
    
    const toggle = page.getByTestId(toggleTestId);
    await expect(toggle).toBeVisible();
    await toggle.click();
    
    const panel = page.locator('#filter-dropdown-panel');
    await expect(panel).toBeVisible();

    const foodCategoryBtn = page.getByTestId('ag-filter-category-food');
    await expect(foodCategoryBtn).toBeVisible();

    // Track requests
    const requestPromise = page.waitForRequest(request => 
      request.url().includes('/listings/fragment') && request.url().includes('type=Food'),
      { timeout: 20000 }
    );

    await foodCategoryBtn.click();
    await requestPromise;

    // Assert JavaScript state
    const filterState = await page.evaluate(() => (window as any).filterState);
    expect(filterState.type).toBe('Food');

    // Assert OOB swap target updates
    await expect(page.locator('#featured-section')).toBeAttached();
  });

  test('should update city and radius state', async ({ page }) => {
    const isMobile = await page.evaluate(() => window.innerWidth < 768);
    const toggleTestId = isMobile ? 'ag-home-filters-toggle-mobile' : 'ag-home-filters-toggle-desktop';
    await page.getByTestId(toggleTestId).click();
    
    const panel = page.locator('#filter-dropdown-panel');
    await expect(panel).toBeVisible();

    const cityInput = page.locator('#filter-city');

    // Track requests using a collection to handle multiple htmx triggers
    const requests: string[] = [];
    page.on('request', request => {
      if (request.url().includes('/listings/fragment')) {
        requests.push(request.url());
      }
    });

    await cityInput.fill('Lagos');
    await cityInput.dispatchEvent('change');
    
    // Wait for the htmx delay for city input (800ms)
    await page.waitForTimeout(1000);

    const radiusBtn = page.locator('[data-radius-value="100"]');
    await radiusBtn.click();
    
    // Wait for requests to settle
    await page.waitForTimeout(2000);

    // Assert that we have a request with Lagos and 100
    const match = requests.some(url => 
      url.toLowerCase().includes('city=lagos') && url.includes('radius=100')
    );
    
    expect(match).toBe(true);

    // Assert JavaScript state
    const filterState = await page.evaluate(() => (window as any).filterState);
    expect(filterState.city).toBe('Lagos');
    expect(filterState.radius).toBe('100');
  });
});
