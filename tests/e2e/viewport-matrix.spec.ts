import { test, expect } from '@playwright/test';

/**
 * Viewport Matrix Assertion Suite
 * Verifies layout integrity across critical routes and mandatory viewports.
 * Calculation: Hero content must not overlap with the sticky header (rect.top >= header.height).
 */

const routes = [
  { name: 'Home', path: '/' },
  { name: 'Search', path: '/?q=African' },
  { name: 'Listing Details', path: '/listings/001500a7-70cd-452e-9cfa-46014c625371' },
];

for (const route of routes) {
  test.describe(`Route: ${route.name}`, () => {
    test(`layout integrity check`, async ({ page }) => {
      // Navigate to the route
      await page.goto(route.path);
      
      // Wait for the page to be stable
      await page.waitForLoadState('load');

      // Identify the sticky header
      const header = page.locator('header').first();
      const headerCount = await header.count();
      
      if (headerCount > 0) {
        // Identify the first meaningful content block after the header
        const hero = page.locator('main > div, main > section, section').first();
        
        await expect(header).toBeVisible();
        await expect(hero).toBeVisible();

        const headerBox = await header.boundingBox();
        const heroBox = await hero.boundingBox();

        if (headerBox && heroBox) {
          /**
           * Assertion: Hero content does not overlap with sticky header.
           * Logic: rect.top >= header.height
           * We allow a 1px margin for rounding/rendering sub-pixel variances.
           */
          expect(heroBox.y).toBeGreaterThanOrEqual(headerBox.height - 1);
        }
      } else {
        // Still verify that some content exists for fragment routes
        await expect(page.locator('body')).not.toBeEmpty();
      }
    });
  });
}
