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

    const suyaPill = page.getByTestId('ag-home-filter-pill-suya');
    await expect(suyaPill).toBeVisible();

    // Intercept HTMX category filter request
    const [request] = await Promise.all([
      page.waitForRequest(req => req.url().includes('/listings/fragment')),
      suyaPill.click()
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
});
