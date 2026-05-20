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
});

