import { test, expect } from '@playwright/test';

test.describe('Owner Claim Flow', () => {
  // Use serial mode because tests share the same dev user and database state
  test.describe.configure({ mode: 'serial' });

  test.beforeEach(async ({ page }) => {
    page.on('console', msg => {
      console.log(`[BROWSER LOG (${msg.type()})] ${msg.text()}`);
    });
  });

  /**
   * Project-isolation: each Playwright project (Mobile, Desktop) uses a
   * different card index and claimer email so they never collide on the
   * shared SQLite database — even when fullyParallel is true at the
   * config level.
   */
  function projectCardIndex(projectName: string): number {
    return projectName === 'Mobile' ? 0 : 1;
  }

  function projectClaimerEmail(projectName: string): string {
    return `claimer-${projectName.toLowerCase()}@agbalumo.com`;
  }

  test('anonymous user clicking claim redirects to login', async ({ page }) => {
    await page.goto('/');
    
    // Toggle filters open
    await page.getByTestId('ag-home-filters-toggle-desktop').click();
    
    // Wait for panel
    const panel = page.locator('#filter-dropdown-panel');
    await expect(panel).toBeVisible();
    
    // Click on Business category to show claimable listings
    const bizCategoryBtn = page.getByTestId('ag-filter-category-Business');
    await expect(bizCategoryBtn).toBeVisible();
    await bizCategoryBtn.click();
    
    // Wait for HTMX swap to complete by asserting that the listings container now contains a Business card
    const businessCard = page.locator('#listings-container').getByText('Business').first();
    await expect(businessCard).toBeVisible({ timeout: 15000 });
    
    const card = page.locator('#listings-container [data-testid="ag-listing-card"]').first();
    await expect(card).toBeVisible();
    
    // Open the detail modal
    const overlay = card.getByTestId('ag-listing-card-overlay').first();
    await overlay.click();
    
    // Wait for the modal/detail content to load
    const claimBtn = page.getByTestId('ag-claim-listing-btn');
    await expect(claimBtn).toBeVisible({ timeout: 10000 });
    
    // Track redirect to /auth/google/login
    let redirectedToLogin = false;
    page.on('request', req => {
      if (req.url().includes('/auth/google/login')) {
        redirectedToLogin = true;
      }
    });

    // Click claim button
    await claimBtn.click();
    
    // Wait a brief moment to allow redirects to trigger
    await page.waitForTimeout(3000);
    
    expect(redirectedToLogin).toBe(true);
  });

  test('authenticated user can successfully submit claim request', async ({ page }, testInfo) => {
    const email = projectClaimerEmail(testInfo.project.name);
    const cardIdx = projectCardIndex(testInfo.project.name);

    // Log in as dev user with project-specific email
    const response = await page.goto(`/auth/dev?email=${email}`);
    expect(response?.status()).toBe(200);
    await page.waitForURL('/');
    
    // Toggle filters open
    await page.getByTestId('ag-home-filters-toggle-desktop').click();
    
    // Wait for panel
    const panel = page.locator('#filter-dropdown-panel');
    await expect(panel).toBeVisible();
    
    // Click on Business category to show claimable listings
    const bizCategoryBtn = page.getByTestId('ag-filter-category-Business');
    await expect(bizCategoryBtn).toBeVisible();
    await bizCategoryBtn.click();
    
    // Wait for HTMX swap
    const businessCard = page.locator('#listings-container').getByText('Business').first();
    await expect(businessCard).toBeVisible({ timeout: 15000 });

    const container = page.locator('#listings-container');
    await container.waitFor({ state: 'attached', timeout: 15000 });
    
    // Select a project-specific card to avoid cross-project collisions
    const card = page.locator('#listings-container [data-testid="ag-listing-card"]').nth(cardIdx);
    await expect(card).toBeVisible();
    
    // Get listing ID
    const overlay = card.getByTestId('ag-listing-card-overlay').first();
    const hxGet = await overlay.getAttribute('hx-get');
    const listingId = hxGet?.split('/').pop();
    expect(listingId).toBeDefined();
    
    // Open the detail modal
    await overlay.click();
    
    const claimBtn = page.getByTestId('ag-claim-listing-btn');
    await expect(claimBtn).toBeVisible({ timeout: 10000 });
    
    // Click the claim button
    await claimBtn.click();
    
    // Verify HTMX response swaps content inside the modal dialog to "Claim Pending Review"
    const pendingText = page.locator('span:has-text("Claim Pending Review")');
    await expect(pendingText).toBeVisible({ timeout: 10000 });
  });
  
  test('admin can see claim request and unassign ownership from dashboard', async ({ page }) => {
    // 1. Authenticate as admin via dev login + access code promoter
    await page.goto('/auth/dev?email=admin@agbalumo.com');
    await page.waitForURL('/');
    
    // Navigate to admin login. If user was already promoted (from a prior run),
    // /admin/login will 307 redirect straight to /admin dashboard.
    await page.goto('/admin/login');
    
    // Check if we need to enter the access code or are already on the dashboard
    if (page.url().includes('/admin/login') || !page.url().endsWith('/admin')) {
      const codeInput = page.locator('input[name="code"]');
      await expect(codeInput).toBeVisible({ timeout: 10000 });
      await codeInput.fill('agbalumo2024');
      await page.locator('button[type="submit"]').click();
      await page.waitForURL('**/admin', { timeout: 10000 });
    }
    
    expect(page.url()).toContain('/admin');
    
    // 2. Open moderation modal from "Pending Claims" banner button
    const pendingClaimsBannerBtn = page.getByTestId('metric-stat-btn-Pending Claims');
    await expect(pendingClaimsBannerBtn).toBeVisible({ timeout: 10000 });
    await pendingClaimsBannerBtn.click();
    
    // Wait for moderation modal to load
    const moderationModal = page.locator('#moderationModal');
    await expect(moderationModal).toBeVisible({ timeout: 10000 });
    
    // Click approve button inside moderation queue modal
    const approveBtn = page.locator('[data-testid^="approve-claim-"]').first();
    await expect(approveBtn).toBeVisible({ timeout: 10000 });
    await approveBtn.click();
    
    // Close moderation modal
    const closeBtn = page.locator('#admin-modal-container button[aria-label="Close"]').first();
    await closeBtn.click();
    
    // 3. Go to admin listings management to verify unassignment
    await page.goto('/admin/listings');
    
    // Wait for the listings table to load
    const table = page.locator('table');
    await expect(table).toBeVisible();
    
    // Locate the first unassign button (rendered now that the claim is approved!)
    const unassignBtn = page.locator('[data-testid^="ag-listing-unassign-btn-"]').first();
    await expect(unassignBtn).toBeVisible({ timeout: 15000 });
    
    // Setup handling for the confirm dialog
    page.once('dialog', async dialog => {
      expect(dialog.message()).toContain('Are you sure you want to unassign');
      await dialog.accept();
    });
    
    // Click Unassign button
    await unassignBtn.click();
    
    // Verify HTMX swaps the row and it now says "UNOWNED"
    const unownedLabel = page.locator('span:has-text("UNOWNED")').first();
    await expect(unownedLabel).toBeVisible({ timeout: 10000 });
  });
});
