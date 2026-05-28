import { test, expect } from '@playwright/test';

test.describe('Feedback Modal Closing and Leak Prevention', () => {
  test.beforeEach(async ({ page }) => {
    // Log browser messages for debugging
    page.on('console', msg => {
      console.log(`[BROWSER LOG (${msg.type()})] ${msg.text()}`);
    });
  });

  test('should not reopen feedback modal when clicking on a listing card after closing feedback modal', async ({ page }, testInfo) => {
    const isMobile = testInfo.project.name === 'Mobile';
    
    // 1. Auth as dev user
    const authResponse = await page.goto('/auth/dev');
    expect(authResponse?.status()).toBe(200);
    await page.waitForURL('/');

    // Ensure JS is initialized
    await page.waitForFunction(() => typeof (window as any).filterState !== 'undefined', { timeout: 10000 });

    // Ensure listings are loaded
    await page.locator('#listings-container').waitFor({ state: 'attached', timeout: 15000 });

    // 2. Open feedback modal
    const feedbackBtn = page.getByTestId(isMobile ? 'mobile-feedback-btn' : 'ag-nav-feedback-btn');
    await expect(feedbackBtn).toBeVisible();
    await feedbackBtn.click();

    // Verify feedback modal is open
    const feedbackModal = page.locator('dialog#feedback-modal');
    await expect(feedbackModal).toBeVisible();

    // 3. Close the feedback modal using the close button (X)
    const closeBtn = feedbackModal.locator('.btn-close').first();
    await expect(closeBtn).toBeVisible();
    await closeBtn.click();

    // Verify feedback modal is closed (either hidden or removed)
    await expect(feedbackModal).not.toBeVisible();

    // 4. Click on the first listing card to open details modal
    const listingCard = page.getByTestId('ag-listing-card').first();
    await expect(listingCard).toBeVisible();

    const cardOverlay = listingCard.locator('button.absolute.inset-0').first();
    await cardOverlay.scrollIntoViewIfNeeded();

    // Click to open details
    await cardOverlay.click({ force: true });

    // Verify details modal is open
    const detailsModal = page.locator('dialog[id^="detail-modal-"]');
    await expect(detailsModal).toBeVisible({ timeout: 15000 });

    // EXPECTATION: Feedback modal should NOT be open/visible!
    await expect(feedbackModal).not.toBeVisible();
  });
});
