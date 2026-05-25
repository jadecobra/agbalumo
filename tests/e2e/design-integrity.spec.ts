import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('Design Integrity & Cohesion', () => {
  test('modals should have exactly one close button', async ({ page }) => {
    await page.goto('/');
    
    // Open Listing Detail Modal
    const firstCard = page.locator('[data-testid="ag-listing-card"]').first();
    await firstCard.locator('div[hx-get^="/listings/"]').click();
    
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible();
    
    // Check for redundant close buttons
    const closeButtons = modal.locator('button[data-modal-action="close"]');
    const count = await closeButtons.count();
    expect(count, `Expected exactly 1 close button in listing modal, found ${count}`).toBe(1);
    
    // Close modal
    await closeButtons.first().click();
    await expect(modal).not.toBeVisible();
    
    // Open Create Listing Modal (as dev)
    await page.goto('/auth/dev');
    await page.goto('/');
    const postBtn = page.locator('[data-testid="ag-nav-post-btn-desktop"]:visible, [data-testid="ag-nav-post-btn-mobile"]:visible').first();
    await postBtn.click();
    
    await expect(modal).toBeVisible();
    const createCloseButtons = modal.locator('button[data-modal-action="close"]');
    const createCount = await createCloseButtons.count();
    expect(createCount, `Expected exactly 1 close button in create modal, found ${createCount}`).toBe(1);
  });

  test('modals should pass accessibility contrast checks', async ({ page }) => {
    await page.goto('/');
    
    // Open Listing Detail Modal
    await page.locator('[data-testid="ag-listing-card"]').first().locator('div[hx-get^="/listings/"]:visible').click();
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible();
    
    const results = await new AxeBuilder({ page })
      .include('dialog[open]')
      .analyze();
      
    const contrastViolations = results.violations.filter(v => v.id === 'color-contrast');
    expect(contrastViolations, `Contrast violations in detail modal: ${JSON.stringify(contrastViolations, null, 2)}`).toHaveLength(0);
  });

  test('profile modal should have exactly one close button and pass contrast', async ({ page }) => {
    await page.goto('/auth/dev');
    await page.goto('/');
    
    // Open Profile Modal
    const profileBtn = page.locator('[data-testid="ag-nav-profile-btn"]:visible, [data-testid="mobile-account-btn"]:visible, [data-testid="mobile-profile-btn"]:visible').first();
    await profileBtn.click();
    
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible();
    
    const closeButtons = modal.locator('button[data-modal-action="close"]');
    const count = await closeButtons.count();
    expect(count, `Expected exactly 1 close button in profile modal, found ${count}`).toBe(1);
    
    const results = await new AxeBuilder({ page })
      .include('dialog[open]')
      .analyze();
      
    const contrastViolations = results.violations.filter(v => v.id === 'color-contrast');
    expect(contrastViolations, `Contrast violations in profile modal: ${JSON.stringify(contrastViolations, null, 2)}`).toHaveLength(0);
  });

  test('icon containers should not have borders', async ({ page }) => {
    // 1. Location Permission Prompt Icon Container
    await page.goto('/');
    const promptIcon = page.locator('[data-testid="location-permission-prompt"] .bg-earth-accent\\/20');
    const promptClassList = await promptIcon.evaluate(el => el.className);
    expect(promptClassList).not.toContain('border');

    // 2. Verified / Poster Origin Flag Badge in Detail Modal
    const firstCard = page.locator('[data-testid="ag-listing-card"]').first();
    await firstCard.locator('div[hx-get^="/listings/"]').click();
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible();
    
    const flagBadge = modal.locator('div[title*="Origin:"]');
    if (await flagBadge.count() > 0) {
      const flagClassList = await flagBadge.first().evaluate(el => el.className);
      expect(flagClassList).not.toContain('border');
    }

    // 3. Error Page Bounce Orange Icon
    await page.goto('/listings/99999');
    const errorIcon = page.locator('.animate-bounce');
    const errorClassList = await errorIcon.evaluate(el => el.className);
    expect(errorClassList).not.toContain('border');
  });

  test('heart icon should not have borders or backgrounds', async ({ page }) => {
    await page.goto('/auth/dev');
    await page.goto('/');
    const heartBtn = page.locator('[data-testid="ag-save-btn"]').first();
    await expect(heartBtn).toBeVisible({ timeout: 15000 });
    const classList = await heartBtn.evaluate(el => el.className);
    
    // Background and border removal validation
    expect(classList).not.toContain('btn-action-overlay');
    expect(classList).not.toContain('backdrop-blur-sm');
    expect(classList).not.toContain('drop-shadow-md');
    expect(classList).not.toContain('bg-');
    expect(classList).not.toContain('border');
  });

  test('close button should not have borders or backgrounds', async ({ page }) => {
    await page.goto('/');
    
    // Open Listing Detail Modal
    const firstCard = page.locator('[data-testid="ag-listing-card"]').first();
    await firstCard.locator('div[hx-get^="/listings/"]').click();
    const modal = page.locator('dialog[open]');
    await expect(modal).toBeVisible();
    
    // Locate the close button inside the modal
    const closeBtn = modal.locator('button[data-modal-action="close"]');
    await expect(closeBtn).toBeVisible();
    const classList = await closeBtn.evaluate(el => el.className);
    
    // Background and border removal validation
    expect(classList).not.toContain('btn-action-overlay');
    expect(classList).not.toContain('backdrop-blur-sm');
    expect(classList).not.toContain('drop-shadow-md');
    expect(classList).not.toContain('bg-');
  });

  test('cards and badges should have strict visual tokens and casing', async ({ page }) => {
    await page.goto('/');
    
    // 1. Featured Card Badge: Uppercase, high-contrast dark translucent glassmorphism
    const featuredCardBadge = page.locator('.card-juicy .absolute.top-4.left-4 span').first();
    await expect(featuredCardBadge).toBeVisible();
    const featuredBadgeText = await featuredCardBadge.innerText();
    expect(featuredBadgeText).toBe(featuredBadgeText.toUpperCase());
    
    const featuredBadgeClassList = await featuredCardBadge.evaluate(el => el.className);
    expect(featuredBadgeClassList).toContain('bg-earth-dark/60');
    expect(featuredBadgeClassList).not.toContain('bg-earth-cream/20');
    expect(featuredBadgeClassList).toContain('rounded-none');
    
    // 2. Listing Card Badge: Uppercase, sophisticated sand/clay styling
    const listingCard = page.locator('[data-testid="ag-listing-card"]').first();
    // Locate the category tag container
    const listingBadge = listingCard.locator('span.bg-earth-sand\\/60, span.bg-earth-accent').first();
    await expect(listingBadge).toBeVisible();
    const listingBadgeText = await listingBadge.innerText();
    expect(listingBadgeText).toBe(listingBadgeText.toUpperCase());
    
    const listingBadgeClassList = await listingBadge.evaluate(el => el.className);
    expect(listingBadgeClassList).toContain('bg-earth-sand/60');
    expect(listingBadgeClassList).toContain('text-earth-clay');
    expect(listingBadgeClassList).toContain('border-earth-clay/20');
    expect(listingBadgeClassList).not.toContain('bg-earth-accent');
    expect(listingBadgeClassList).not.toContain('capitalize');
    expect(listingBadgeClassList).toContain('uppercase');
    expect(listingBadgeClassList).toContain('rounded-none');
    
    // 3. Fallback Listing Image Section: Category-dynamic dynamic gradients, no grey bg-stone-100 background
    const listingWithFallback = page.locator('[data-testid="ag-listing-card"]', { has: page.locator('.blur-2xl') }).first();
    if (await listingWithFallback.count() > 0) {
      const fallbackDiv = listingWithFallback.locator('.blur-2xl').first().locator('..');
      const fallbackClassList = await fallbackDiv.evaluate(el => el.className);
      expect(fallbackClassList).toContain('bg-gradient-to-br');
      expect(fallbackClassList).not.toContain('bg-stone-100');
      expect(fallbackClassList).toContain('rounded-none');
    }
  });
});

