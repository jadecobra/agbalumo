// Deprecated: The custom location permission modal has been removed from the product.
// Geolocation is now requested directly on load and on NEAR ME clicks.
// This entire test file is obsolete and kept only as a stub.

import { test, expect } from '@playwright/test';

test.describe.skip('Location Permission Prompt (DEPRECATED - custom modal removed)', () => {
  test.beforeEach(async ({ page }) => {
    // Log browser messages for debugging
    page.on('console', msg => {
      console.log(`[BROWSER ${msg.type().toUpperCase()}] ${msg.text()}`);
    });
    // Clear localStorage and mock geolocation permission query to return 'prompt'
    await page.addInitScript(() => {
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
    await page.goto('/');
  });

  test('should not display the location permission prompt on home load (user must click NEAR ME first)', async ({ page }) => {
    const state = await page.evaluate(async () => {
      try {
        const res = await navigator.permissions.query({ name: 'geolocation' });
        return res.state;
      } catch (e) {
        return 'error: ' + e.message;
      }
    });
    console.log('=== GEOLOCATION PERMISSION STATE ===:', state);

    const prompt = page.getByTestId('location-permission-prompt');
    await expect(prompt).not.toBeVisible();

    // Per spec: explicit user action (NEAR ME) is required to surface the app's permission explainer
    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await nearMeBtn.click();

    await expect(prompt).toBeVisible({ timeout: 5000 });
    await expect(prompt).toContainText('Allow location to instantly show African owned spots');
  });

  test('should dismiss the prompt when clicking dismiss (and denial persists for new tab / reload)', async ({ page }) => {
    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await nearMeBtn.click();

    const prompt = page.getByTestId('location-permission-prompt');
    await expect(prompt).toBeVisible({ timeout: 5000 });

    await page.evaluate(() => {
        const btn = document.getElementById('location-dismiss-btn');
        if (btn) btn.click();
    });

    await expect(prompt).not.toBeVisible();

    // Denial is now in localStorage; reload (simulates new tab in some scenarios) must not auto-show
    await page.reload();
    await expect(prompt).not.toBeVisible();

    // But explicit NEAR ME click still surfaces the explainer (user can change mind)
    await nearMeBtn.click();
    await expect(prompt).toBeVisible();
  });

  test('should trigger geolocation when clicking allow', async ({ page }) => {
    const nearMeBtn = page.getByTestId('ag-home-near-me-btn');
    await nearMeBtn.click();

    const prompt = page.getByTestId('location-permission-prompt');
    await expect(prompt).toBeVisible({ timeout: 5000 });

    const allowBtn = page.getByTestId('location-permission-allow');

    // Mock geolocation success
    await page.addInitScript(() => {
      const mockGeolocation = {
        getCurrentPosition: (success: any) => {
          success({
            coords: {
              latitude: 51.5074,
              longitude: 0.1278,
            },
          });
        },
      };
      (navigator as any).geolocation = mockGeolocation;
    });

    await page.evaluate(() => {
        const btn = document.getElementById('location-allow-btn');
        if (btn) btn.click();
    });

    await expect(prompt).not.toBeVisible();
  });
});
