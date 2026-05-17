import { test, expect } from '@playwright/test';

test.describe('Location Permission Prompt', () => {
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

  test('should display the location permission prompt on home load', async ({ page }) => {
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
    await expect(prompt).toBeVisible({ timeout: 10000 });
    await expect(prompt).toContainText('Allow location to instantly show African owned spots');
  });

  test('should dismiss the prompt when clicking dismiss', async ({ page }) => {
    const prompt = page.getByTestId('location-permission-prompt');
    await expect(prompt).toBeVisible({ timeout: 15000 });
    
    // Wait for JS delay (1500ms) + animation (500ms) + buffer
    await page.waitForTimeout(3000);

    await page.evaluate(() => {
        const btn = document.getElementById('location-dismiss-btn');
        if (btn) btn.click();
    });
    
    await expect(prompt).not.toBeVisible();
    
    // Verify persistence
    await page.reload();
    await expect(prompt).not.toBeVisible();
  });

  test('should trigger geolocation when clicking allow', async ({ page }) => {
    const prompt = page.getByTestId('location-permission-prompt');
    await expect(prompt).toBeVisible({ timeout: 15000 });

    const allowBtn = page.getByTestId('location-permission-allow');
    
    // Mock geolocation
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

    // Wait for JS delay (1500ms) + animation (500ms) + buffer
    await page.waitForTimeout(3000);
    
    await page.evaluate(() => {
        const btn = document.getElementById('location-allow-btn');
        if (btn) btn.click();
    });
    
    // The prompt should hide after allowing
    await expect(prompt).not.toBeVisible();
  });
});
