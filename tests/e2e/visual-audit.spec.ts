import { test, expect, Page } from '@playwright/test';

const PUBLIC_ROUTES = ['/', '/about'];
const AUTH_ROUTES = ['/profile', '/saved'];
const ADMIN_ROUTES = ['/admin', '/admin/listings', '/admin/users'];

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

async function getFirstListingId(page: Page): Promise<string | null> {
  await page.goto('/');
  await page.waitForLoadState('networkidle');
  // Ensure we get a listing card, not a pagination fragment
  const cardOverlay = page.locator('[data-testid="ag-listing-card"] button[hx-get^="/listings/"]').first();
  const hxGet = await cardOverlay.getAttribute('hx-get');
  if (!hxGet) return null;
  const parts = hxGet.split('/');
  const id = parts.pop();
  return id === 'fragment' ? null : id || null;
}

test.describe('Visual Audit', () => {
  test.afterEach(async ({ page }, testInfo) => {
    if (testInfo.status !== testInfo.expectedStatus) {
      const screenshotPath = `test-results/failure-${testInfo.title.replace(/\s+/g, '-')}.png`;
      await page.screenshot({ path: screenshotPath, fullPage: true });
      console.log(`Screenshot saved to ${screenshotPath}`);
    }
  });

  test('no console errors on public routes', async ({ page }) => {
    for (const route of PUBLIC_ROUTES) {
      const errors: string[] = [];
      page.on('console', msg => {
        const text = msg.text();
        if (msg.type() === 'error' && 
            !text.includes('favicon.ico') && 
            !text.includes('t0.gstatic.com') && 
            !text.startsWith('Failed to load resource:')) {
          errors.push(`${route}: ${text}`);
        }
      });
      await page.goto(route);
      await page.waitForLoadState('networkidle');
      expect(errors, `Errors on ${route}: ${errors.join(', ')}`).toHaveLength(0);
    }
  });

  test('no console errors on auth routes', async ({ page }) => {
    await loginAsDev(page);
    for (const route of AUTH_ROUTES) {
      const errors: string[] = [];
      page.on('console', msg => {
        const text = msg.text();
        if (msg.type() === 'error' && 
            !text.includes('favicon.ico') && 
            !text.includes('t0.gstatic.com') && 
            !text.startsWith('Failed to load resource:')) {
          errors.push(`${route}: ${text}`);
        }
      });
      await page.goto(route);
      await page.waitForLoadState('networkidle');
      expect(errors, `Errors on ${route}: ${errors.join(', ')}`).toHaveLength(0);
    }
  });

  test('no console errors on admin routes', async ({ page }) => {
    await loginAsAdmin(page);
    for (const route of ADMIN_ROUTES) {
      const errors: string[] = [];
      page.on('console', msg => {
        const text = msg.text();
        if (msg.type() === 'error' && 
            !text.includes('favicon.ico') && 
            !text.includes('t0.gstatic.com') && 
            !text.startsWith('Failed to load resource:')) {
          errors.push(`${route}: ${text}`);
        }
      });
      await page.goto(route);
      await page.waitForLoadState('networkidle');
      expect(errors, `Errors on ${route}: ${errors.join(', ')}`).toHaveLength(0);
    }
  });

  test('listing detail page opens modal and shows title', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    const firstListing = page.locator('[data-testid="ag-listing-card"]').first();
    const rawTitle = await firstListing.locator('h3').first().innerText();
    const title = rawTitle.replace('🍊', '').trim();
    const overlay = firstListing.locator('button[hx-get^="/listings/"]').first();
    await overlay.click();
    
    // Wait for modal to appear (dialog element)
    const modalTitle = page.locator('dialog[open] h2').first();
    await expect(modalTitle).toBeVisible();
    expect(await modalTitle.innerText()).toContain(title.trim());
  });

  test('listing detail modal preserves background scroll position after closing', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    
    // Dismiss location prompt to avoid blocking or shifting layouts
    const dismissBtn = page.getByTestId('location-permission-dismiss');
    if (await dismissBtn.isVisible()) {
      await dismissBtn.click();
      await page.waitForTimeout(500);
    }

    // Scroll down to ensure we are not at the top
    await page.evaluate(() => window.scrollTo(0, 400));
    await page.waitForTimeout(200);

    // Open first listing detail modal
    const firstListing = page.locator('[data-testid="ag-listing-card"]').first();
    const overlay = firstListing.locator('button[hx-get^="/listings/"]').first();
    
    // Scroll element into view and wait to settle before measuring baseline scroll
    await overlay.scrollIntoViewIfNeeded();
    await page.waitForTimeout(200);
    const scrollBefore = await page.evaluate(() => window.scrollY);
    expect(scrollBefore).toBeGreaterThan(100);

    await overlay.click();

    // Verify modal is open
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible();

    // Close the modal
    const closeBtn = modal.locator('[data-modal-action="close"]').first();
    await closeBtn.click();
    await expect(modal).not.toBeVisible();
    await page.waitForTimeout(200);

    // Assert that the scroll position was preserved
    const scrollAfter = await page.evaluate(() => window.scrollY);
    expect(scrollAfter).toBeCloseTo(scrollBefore, 10);
  });


  test('single h1 per major page', async ({ page }) => {
    // Check Home
    await page.goto('/');
    await expect(page.locator('h1')).toHaveCount(1);
    
    // Check About
    await page.goto('/about');
    await expect(page.locator('h1')).toHaveCount(1);
    
    // Check Admin Dashboard
    await loginAsAdmin(page);
    await page.goto('/admin');
    await expect(page.locator('h1')).toHaveCount(1);
  });

  test('create listing modal opens and has no console errors', async ({ page }, testInfo) => {
    await loginAsDev(page);
    const errors: string[] = [];
    page.on('console', msg => {
      const text = msg.text();
      if (msg.type() === 'error' && 
          !text.includes('favicon.ico') && 
          !text.includes('t0.gstatic.com') && 
          !text.startsWith('Failed to load resource:')) {
        errors.push(text);
      }
    });
    
    await page.goto('/');
    const postBtn = testInfo.project.name === 'Mobile' 
      ? page.locator('[data-testid="ag-nav-post-btn-mobile"]')
      : page.locator('[data-testid="ag-nav-post-btn-desktop"]');
    await postBtn.click();
    
    // Wait for modal
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible();
    await expect(modal.locator('h2').first()).toContainText('Add');
    
    // Verify some key fields exist
    await expect(page.locator('input[name="title"]')).toBeVisible();
    await expect(page.locator('textarea[name="description"]')).toBeVisible();
    
    expect(errors, `Errors on Post modal: ${errors.join(', ')}`).toHaveLength(0);
  });

  test('cards visible above fold at desktop', async ({ page }, testInfo) => {
    if (testInfo.project.name !== 'Desktop') test.skip();

    await page.goto('/');
    await page.waitForLoadState('load');
    const firstCard = page.locator('[data-testid="ag-listing-card"]').first();
    await expect(firstCard).toBeVisible();
    const box = await firstCard.boundingBox();
    if (box) expect(box.y).toBeLessThan(900);
  });

  test('touch targets >= 44px on all major surfaces', async ({ page }, testInfo) => {
    if (testInfo.project.name !== 'Mobile') test.skip();

    const routesToAudit = [...PUBLIC_ROUTES, '/admin/login'];
    for (const route of routesToAudit) {
      await page.goto(route);
      await page.waitForLoadState('load');
      const interactiveElements = page.locator('button, a[href], [hx-get], [hx-post]');
      const count = await interactiveElements.count();
      const failures: string[] = [];

      for (let i = 0; i < count; i++) {
        const element = interactiveElements.nth(i);
        if (!await element.isVisible()) continue;
        if (await element.evaluate(el => el.classList.contains('sr-only'))) continue;

        const box = await element.boundingBox();
        if (box && (box.width < 44 || box.height < 44)) {
          const text = (await element.innerText()).trim() || await element.getAttribute('aria-label') || 'unnamed';
          failures.push(`${route}: ${text} [${box.width}x${box.height}]`);
        }
      }
      expect(failures, `Small touch targets on ${route}: \n${failures.join('\n')}`).toHaveLength(0);
    }
  });

  test('all HTMX buttons have data-testid on all routes', async ({ page }) => {
    const routes = [...PUBLIC_ROUTES, '/admin/login'];
    for (const route of routes) {
      await page.goto(route);
      await page.waitForLoadState('load');
      const htmxElements = page.locator('button[hx-get], button[hx-post], button[hx-delete], a[hx-get]');
      const count = await htmxElements.count();
      for (let i = 0; i < count; i++) {
        await expect(htmxElements.nth(i)).toHaveAttribute('data-testid');
      }
    }
  });
});




