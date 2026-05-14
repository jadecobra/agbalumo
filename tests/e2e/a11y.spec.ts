import { test, expect, Page } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

async function loginAsDev(page: Page) {
  await page.goto('/auth/dev');
  await page.waitForURL('/');
}

async function loginAsAdmin(page: Page) {
  await loginAsDev(page);
  await page.goto('/admin/login');
  await page.waitForLoadState('networkidle');
  
  // Check if we were automatically redirected because we're already admin
  if (page.url().includes('/admin') && !page.url().includes('/login')) {
    return;
  }

  const codeInput = page.locator('input[name="code"]');
  try {
    await codeInput.waitFor({ state: 'visible', timeout: 5000 });
    await codeInput.fill('agbalumo2024');
    await page.click('button[type="submit"]');
    await page.waitForURL(url => url.pathname === '/admin', { timeout: 10000 });
  } catch (e) {
    // If it's not visible, maybe we're already redirected or something else is wrong
    if (!page.url().includes('/admin')) {
      throw new Error(`Admin login failed: Not redirected to /admin and code input not found. Current URL: ${page.url()}`);
    }
  }
}

test.describe('Accessibility Audits', () => {
  test('sandbox should have no detectable a11y violations', async ({ page }) => {
    await page.goto('/sandbox', { waitUntil: 'networkidle', timeout: 30000 });
    const results = await new AxeBuilder({ page })
      .include('main')
      .analyze();
    expect(results.violations).toEqual([]);
  });

  test('main feed should have no detectable a11y violations', async ({ page }) => {
    await page.goto('/', { waitUntil: 'networkidle' });
    await expect(page.locator('[data-testid="ag-listing-card"]').first()).toBeVisible();
    
    const results = await new AxeBuilder({ page })
      .include('body')
      .analyze();
    
    if (results.violations.length > 0) {
      console.log('Main Feed A11y Violations:', JSON.stringify(results.violations, null, 2));
    }
    expect(results.violations).toEqual([]);
  });

  test('listing detail modal should have no detectable a11y violations', async ({ page }) => {
    await page.goto('/', { waitUntil: 'networkidle' });
    await page.locator('[data-testid="ag-listing-card"]').first().click();
    
    const modal = page.locator('dialog[id^="detail-modal-"]');
    await expect(modal).toBeVisible();
    
    const results = await new AxeBuilder({ page })
      .include('dialog')
      .analyze();
    
    if (results.violations.length > 0) {
      console.log('Listing Detail Modal A11y Violations:', JSON.stringify(results.violations, null, 2));
    }
    expect(results.violations).toEqual([]);
  });

  test('post listing modal should have no detectable a11y violations', async ({ page }, testInfo) => {
    await loginAsDev(page);
    await page.goto('/', { waitUntil: 'networkidle' });
    
    const postBtn = testInfo.project.name === 'Mobile' 
      ? page.locator('[data-testid="ag-nav-post-btn-mobile"]')
      : page.locator('[data-testid="ag-nav-post-btn-desktop"]');
    await postBtn.click();
    
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible();
    
    const results = await new AxeBuilder({ page })
      .include('dialog')
      .analyze();
    
    if (results.violations.length > 0) {
      console.log('Post Listing Modal A11y Violations:', JSON.stringify(results.violations, null, 2));
    }
    expect(results.violations).toEqual([]);
  });

  test('admin dashboard should have no detectable a11y violations', async ({ page }) => {
    await loginAsAdmin(page);
    await page.goto('/admin', { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toBeVisible();
    
    const results = await new AxeBuilder({ page })
      .include('main')
      .analyze();
    
    if (results.violations.length > 0) {
      console.log('Admin Dashboard A11y Violations:', JSON.stringify(results.violations, null, 2));
    }
    expect(results.violations).toEqual([]);
  });
});
