import { test, expect } from '@playwright/test';

test.describe('Near Me Geolocation UX', () => {
  test.beforeEach(async ({ context, page }) => {
    // Log browser messages and requests for debugging
    page.on('console', msg => {
      console.log(`[BROWSER ${msg.type().toUpperCase()}] ${msg.text()}`);
    });

    page.on('request', req => {
      if (req.url().includes('/listings/fragment')) {
        console.log(`[HTMX REQUEST]: ${req.url()}`);
      }
    });

    // Clear storage and mock permission query before each test
    await page.addInitScript(() => {
      window.sessionStorage.clear();
      window.localStorage.clear();
      
      if (navigator.permissions) {
        const originalQuery = navigator.permissions.query;
        navigator.permissions.query = async (descriptor) => {
          if (descriptor && descriptor.name === 'geolocation') {
            return {
              state: 'prompt',
              onchange: null,
            } as any;
          }
          return originalQuery.call(navigator.permissions, descriptor);
        };
      }
    });

    // Grant geolocation permissions and set coordinates natively at context level
    await context.grantPermissions(['geolocation']);
    await context.setGeolocation({ latitude: 6.5244, longitude: 3.3792 });
  });

  test('Test 1: Near Me button exists and displays "Near Me"', async ({ page }) => {
    await page.goto('/');
    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toBeVisible();
    await expect(nearMeBtn).toContainText('Near Me');
  });

  test('Test 2: Clicking the button triggers mock geolocation, loads geolocated fragment, stores coords in sessionStorage, and changes the button to active state', async ({ page }) => {
    await page.goto('/');

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toBeVisible();

    // Click near me button
    await nearMeBtn.click();

    // Check button updates to active state text
    await expect(nearMeBtn).toContainText('Nearby');

    // Confirm sessionStorage has the correct coords
    const sessionStorageCoords = await page.evaluate(() => {
      return {
        lat: sessionStorage.getItem('agbalumo_lat'),
        lng: sessionStorage.getItem('agbalumo_lng')
      };
    });

    expect(sessionStorageCoords.lat).toBe('6.5244');
    expect(sessionStorageCoords.lng).toBe('3.3792');

    // Confirm that the listings container opacity was reset and is not stuck at 0.3
    const container = page.locator('#listings-container');
    await expect(container).toHaveJSProperty('style.opacity', '');

    // Click near me button again to toggle it off (de-select)
    await nearMeBtn.click();

    // Check button resets back to "Near Me"
    await expect(nearMeBtn).toContainText('Near Me');

    // Confirm sessionStorage has cleared the coordinates
    const clearedCoords = await page.evaluate(() => {
      return {
        lat: sessionStorage.getItem('agbalumo_lat'),
        lng: sessionStorage.getItem('agbalumo_lng')
      };
    });
    expect(clearedCoords.lat).toBeNull();
    expect(clearedCoords.lng).toBeNull();

    // Confirm that the container opacity style remains clean
    await expect(container).toHaveJSProperty('style.opacity', '');
  });

  test('Test 3: Subsequent category filter actions automatically include lat and lng in their HTMX queries', async ({ page }) => {
    // Inject coordinates directly into sessionStorage before page navigation
    await page.addInitScript(() => {
      sessionStorage.setItem('agbalumo_lat', '6.5244');
      sessionStorage.setItem('agbalumo_lng', '3.3792');
    });

    await page.goto('/');

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toContainText('Nearby');

    // Toggle filters panel to access category
    const filtersToggle = page.getByTestId('ag-home-filters-toggle-desktop');
    await expect(filtersToggle).toBeVisible();
    await filtersToggle.click();

    const foodBtn = page.getByTestId('ag-filter-category-food');
    await expect(foodBtn).toBeVisible();

    // Intercept HTMX category filter request
    const [request] = await Promise.all([
      page.waitForRequest(req => req.url().includes('/listings/fragment')),
      foodBtn.click()
    ]);

    const url = new URL(request.url());
    expect(url.searchParams.get('lat')).toBe('6.5244');
    expect(url.searchParams.get('lng')).toBe('3.3792');
  });

  test('Test 4: Typing/searching a custom city clears sessionStorage coordinates and resets Near Me button to default', async ({ page }) => {
    // Start with already located state in sessionStorage before page navigation
    await page.addInitScript(() => {
      sessionStorage.setItem('agbalumo_lat', '6.5244');
      sessionStorage.setItem('agbalumo_lng', '3.3792');
    });

    await page.goto('/');

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toContainText('Nearby');

    // Toggle filters panel to access city filter field
    const filtersToggle = page.getByTestId('ag-home-filters-toggle-desktop');
    await expect(filtersToggle).toBeVisible();
    await filtersToggle.click();

    // Wait 500ms for autofocus event on filters.js to fire and settle
    await page.waitForTimeout(500);

    const cityInput = page.locator('#filter-city');
    await expect(cityInput).toBeVisible();

    // Focus and click the city input explicitly to reclaim active focus from search input
    await cityInput.focus();
    await cityInput.click();

    // Type a city name character-by-character
    await cityInput.pressSequentially('Lagos', { delay: 100 });

    // Give a delay to allow sessionStorage and DOM updates to finalize
    await page.waitForTimeout(2000);

    // Verify sessionStorage coordinates are cleared
    const sessionStorageCoords = await page.evaluate(() => {
      return {
        lat: sessionStorage.getItem('agbalumo_lat'),
        lng: sessionStorage.getItem('agbalumo_lng')
      };
    });

    expect(sessionStorageCoords.lat).toBeNull();
    expect(sessionStorageCoords.lng).toBeNull();

    // Verify Near Me button reset back to 'Near Me'
    await expect(nearMeBtn).toContainText('Near Me');
  });

  test('Test 5: Loading the page with lat/lng query parameters renders active state and clicking it toggles off', async ({ page }) => {
    // Navigate with lat/lng in URL
    await page.goto('/?lat=6.5244&lng=3.3792');

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toBeVisible();
    await expect(nearMeBtn).toContainText('Nearby');

    // Click to toggle off (deselect)
    await nearMeBtn.click();

    // Check button resets back to "Near Me"
    await expect(nearMeBtn).toContainText('Near Me');

    // Verify sessionStorage coordinates are cleared/null
    const clearedCoords = await page.evaluate(() => {
      return {
        lat: sessionStorage.getItem('agbalumo_lat'),
        lng: sessionStorage.getItem('agbalumo_lng')
      };
    });
    expect(clearedCoords.lat).toBeNull();
    expect(clearedCoords.lng).toBeNull();

    // Verify browser URL has no lat or lng
    const currentUrl = new URL(page.url());
    expect(currentUrl.searchParams.get('lat')).toBeNull();
    expect(currentUrl.searchParams.get('lng')).toBeNull();
  });

  test('Test 6: Deselecting NEARBY does not show radius banner', async ({ page }) => {
    await page.goto('/');

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toBeVisible();

    // Click to activate and wait for HTMX response
    await Promise.all([
      page.waitForResponse(resp => resp.url().includes('/listings/fragment') && resp.status() === 200),
      nearMeBtn.click()
    ]);
    await expect(nearMeBtn).toContainText('Nearby');

    // Click to deactivate and wait for HTMX response
    await Promise.all([
      page.waitForResponse(resp => resp.url().includes('/listings/fragment') && resp.status() === 200),
      nearMeBtn.click()
    ]);
    await expect(nearMeBtn).toContainText('Near Me');

    // The location-status banner should NOT contain "miles" text
    const locationStatus = page.locator('#location-status');
    await expect(locationStatus).not.toContainText('miles');
  });

  test('Test 7: Near Me respects user-selected radius from filter', async ({ page }) => {
    await page.goto('/');

    // Open filter panel and select 5-mile radius
    const filtersToggle = page.getByTestId('ag-home-filters-toggle-desktop');
    await expect(filtersToggle).toBeVisible();
    await filtersToggle.click();

    const radius5Btn = page.locator('[data-radius-value="5"]');
    await expect(radius5Btn).toBeVisible();
    await radius5Btn.click();

    // Now click Near Me
    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');

    // Intercept the HTMX request triggered by Near Me click
    const [request] = await Promise.all([
      page.waitForRequest(req => req.url().includes('/listings/fragment') && req.url().includes('lat=')),
      nearMeBtn.click()
    ]);

    // Verify the request uses the user-selected radius (5), not hardcoded 10
    const url = new URL(request.url());
    expect(url.searchParams.get('radius')).toBe('5');
  });

  test('Test 8: denied geolocation is sticky/retryable (button re-enabled so user can click again for prompt)', async ({ page }) => {
    // Mock geolocation to fail with PERMISSION_DENIED (covers decline-then-reclick symptom)
    await page.addInitScript(() => {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition = (_success: any, error: any, _options: any) => {
          error({
            code: 1,
            message: "Permission denied"
          });
        };
      }
    });

    await page.goto('/');

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toBeVisible();

    await nearMeBtn.click();

    // Denied state is shown and button is re-enabled for immediate re-click (new gesture)
    await expect(nearMeBtn).toContainText('Denied');
    await expect(nearMeBtn).toBeEnabled();

    // No auto-revert for denial (user can tap again to request the native prompt).
    // (Transient errors still flash + revert; this test only exercises the denied path.)
  });

  test('Test 9: clicking NEAR ME after dismissing HTML permission prompt shows prompt again', async ({ page }) => {
    await page.goto('/');

    const prompt = page.getByTestId('location-permission-prompt');
    await expect(prompt).toBeVisible({ timeout: 10000 });

    const dismissBtn = page.locator('#location-dismiss-btn');
    await expect(dismissBtn).toBeVisible();
    await dismissBtn.click();

    await expect(prompt).not.toBeVisible();

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toBeVisible();
    await nearMeBtn.click();

    // The location permission prompt should be shown again!
    await expect(prompt).toBeVisible();
  });

  test('Test 10: clicking NEAR ME after native permission denial shows HTML prompt again', async ({ page }) => {
    // Mock geolocation to fail with PERMISSION_DENIED
    await page.addInitScript(() => {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition = (_success: any, error: any, _options: any) => {
          error({
            code: 1,
            message: "Permission denied"
          });
        };
      }
    });

    await page.goto('/');

    // Dismiss custom prompt first so we can click near-me
    const dismissBtn = page.locator('#location-dismiss-btn');
    await expect(dismissBtn).toBeVisible();
    await dismissBtn.click();

    // Clear the dismissed state before first click so it directly queries native geolocation
    await page.evaluate(() => {
      sessionStorage.removeItem('agbalumo_geo_dismissed');
    });

    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await expect(nearMeBtn).toBeVisible();
    await nearMeBtn.click();


    // Near me button should enter denied state
    await expect(nearMeBtn).toContainText('Denied');

    // Click near-me button again after denial
    await nearMeBtn.click();

    // It should show the HTML permission prompt again
    const prompt = page.getByTestId('location-permission-prompt');
    await expect(prompt).toBeVisible();
  });
});


