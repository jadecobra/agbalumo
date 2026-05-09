import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('Accessibility Audits', () => {
  test('sandbox should have no detectable a11y violations', async ({ page }) => {
    // Increase timeout for slow local server
    await page.goto('/sandbox', { waitUntil: 'networkidle', timeout: 30000 });
    
    console.log('Final URL:', page.url());
    
    // Wait for content to be visible
    const h1 = page.locator('h1');
    try {
      await expect(h1).toBeVisible({ timeout: 10000 });
    } catch (e) {
      console.log('H1 visibility failure. Current HTML:', await page.content());
      throw e;
    }

    const results = await new AxeBuilder({ page })
      .include('main')
      .analyze();
    
    if (results.violations.length > 0) {
      console.log('A11y Violations:', JSON.stringify(results.violations, null, 2));
    }

    expect(results.violations).toEqual([]);
  });
});
