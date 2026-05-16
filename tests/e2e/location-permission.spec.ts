import { test, expect } from '@playwright/test';

test.describe('Location Permission Prompt', () => {
  test.beforeEach(async ({ page }) => {
    // Clear localStorage to ensure prompt shows
    await page.addInitScript(() => {
      window.localStorage.clear();
    });
    await page.goto('/');
  });

  test('should display the location permission prompt on home load', async ({ page }) => {
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
