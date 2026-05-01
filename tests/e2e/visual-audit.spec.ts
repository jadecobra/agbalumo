import { test, expect } from '@playwright/test';

test.describe('Visual Audit', () => {
  test('no console errors or CSP violations', async ({ page }) => {
    const errors: string[] = [];
    page.on('console', msg => {
      if (msg.type() === 'error') {
        const text = msg.text();
        // Filter out known noise (favicon 404 etc)
        if (!text.includes('favicon.ico') && !text.includes('404')) {
          errors.push(text);
        }
      }
    });

    // Navigate to root
    await page.goto('/');
    
    // Wait for network to be idle to catch late-loading script errors
    await page.waitForLoadState('networkidle');

    expect(errors, `Found console errors: ${errors.join(', ')}`).toHaveLength(0);
  });

  test('cards visible above fold at desktop', async ({ page }, testInfo) => {
    // Only at 1440x900 (Desktop project)
    if (testInfo.project.name !== 'Desktop') {
      test.skip();
    }

    await page.goto('/');
    await page.waitForLoadState('load');

    const firstCard = page.locator('[data-testid="ag-listing-card"]').first();
    await expect(firstCard).toBeVisible();

    const box = await firstCard.boundingBox();
    expect(box).not.toBeNull();
    if (box) {
      // Assert y < 900 (Desktop height)
      expect(box.y).toBeLessThan(900);
    }
  });

  test('nav-to-content gap is reasonable', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('load');

    const header = page.locator('header').first();
    const exploreSection = page.locator('#explore-section').first();

    await expect(header).toBeVisible();
    await expect(exploreSection).toBeVisible();

    const headerBox = await header.boundingBox();
    const exploreBox = await exploreSection.boundingBox();

    if (headerBox && exploreBox) {
      const gap = exploreBox.y - (headerBox.y + headerBox.height);
      // Assert < 150px gap between header and first section
      expect(gap).toBeLessThan(150);
    }
  });

  test('all HTMX buttons have data-testid', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('load');

    const htmxElements = page.locator('button[hx-get], button[hx-post], button[hx-delete], a[hx-get]');
    const count = await htmxElements.count();
    
    const failures: string[] = [];
    for (let i = 0; i < count; i++) {
      const element = htmxElements.nth(i);
      const testId = await element.getAttribute('data-testid');
      if (!testId) {
        const html = await element.evaluate(el => el.outerHTML);
        failures.push(html);
      }
    }

    expect(failures, `HTMX elements missing data-testid: \n${failures.join('\n')}`).toHaveLength(0);
  });

  test('touch targets >= 44px on mobile', async ({ page }, testInfo) => {
    // Only in Mobile project
    if (testInfo.project.name !== 'Mobile') {
      test.skip();
    }

    await page.goto('/');
    await page.waitForLoadState('load');

    // Query all button, a[href], [hx-get], [hx-post] visible elements
    const interactiveElements = page.locator('button, a[href], [hx-get], [hx-post]');
    const count = await interactiveElements.count();

    const failures: string[] = [];
    for (let i = 0; i < count; i++) {
      const element = interactiveElements.nth(i);
      
      // Skip elements that are display:none or visibility:hidden
      if (!await element.isVisible()) continue;

      // Skip sr-only elements
      const isSrOnly = await element.evaluate(el => el.classList.contains('sr-only'));
      if (isSrOnly) continue;

      const box = await element.boundingBox();
      if (box) {
        // 40px allows 4px tolerance
        if (box.width < 40 || box.height < 40) {
          const text = (await element.innerText()).trim() || await element.getAttribute('aria-label') || 'unnamed element';
          const selector = await element.evaluate(el => {
            const path = [];
            let curr = el;
            while (curr && curr.nodeType === Node.ELEMENT_NODE) {
              let name = curr.nodeName.toLowerCase();
              if (curr.id) name += `#${curr.id}`;
              else if (curr.className) name += `.${curr.className.split(' ').join('.')}`;
              path.unshift(name);
              curr = (curr as Element).parentNode as Element;
            }
            return path.join(' > ');
          });
          failures.push(`${text} [${box.width}x${box.height}] at ${selector}`);
        }
      }
    }

    expect(failures, `Small touch targets on mobile: \n${failures.join('\n')}`).toHaveLength(0);
  });
});
