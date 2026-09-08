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
    const toggle = page.getByTestId('ag-home-filters-toggle-desktop');
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
    await page.getByTestId('ag-home-filters-toggle-desktop').click();
    
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

    const radiusBtn = page.locator('[data-radius-value="10"]');
    await radiusBtn.click();
    
    // Wait for requests to settle
    await page.waitForTimeout(2000);

    // Assert that we have a request with Lagos and 10
    const match = requests.some(url => 
      url.toLowerCase().includes('city=lagos') && url.includes('radius=10')
    );
    
    expect(match).toBe(true);

    // Assert JavaScript state
    const filterState = await page.evaluate(() => (window as any).filterState);
    expect(filterState.city).toBe('Lagos');
    expect(filterState.radius).toBe('10');
  });

  test('should provide a visible desktop login button', async ({ page }) => {
    // Desktop viewport
    await page.setViewportSize({ width: 1440, height: 900 });
    const signInBtn = page.getByTestId('ag-nav-signin-btn');
    await expect(signInBtn).toBeVisible();
    
    // Check that it doesn't use the dark secondary text class which hides it on the dark header
    const classList = await signInBtn.getAttribute('class');
    expect(classList).not.toContain('text-secondary');
  });

  test('should auto-open post/login modal when action=post query param is present', async ({ page }) => {
    await page.goto('/?action=post');
    const loginModal = page.locator('#login-prompt-modal');
    await expect(loginModal).toBeVisible();
  });

  test('should auto-open post/login modal when hash is #post', async ({ page }) => {
    await page.goto('/#post');
    const loginModal = page.locator('#login-prompt-modal');
    await expect(loginModal).toBeVisible();
  });

  test('should track inbound UTM campaign parameters and log campaign_visit', async ({ page }) => {
    let capturedPayload: any = null;
    await page.route('**/api/metrics', async route => {
      const postData = route.request().postDataJSON();
      if (postData && postData.event === 'campaign_visit') {
        capturedPayload = postData;
      }
      await route.continue();
    });

    await page.goto('/?utm_source=facebook&utm_medium=social&utm_campaign=coverage_gaps');
    await expect.poll(() => capturedPayload).not.toBeNull();
    expect(capturedPayload.metadata.utm_source).toBe('facebook');
    expect(capturedPayload.metadata.utm_medium).toBe('social');
    expect(capturedPayload.metadata.utm_campaign).toBe('coverage_gaps');
  });

});
