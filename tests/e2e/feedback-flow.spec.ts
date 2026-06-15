import { test, expect } from '@playwright/test';

test.describe('User Feedback Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Log browser messages
    page.on('console', msg => {
      console.log(`[BROWSER LOG (${msg.type()})] ${msg.text()}`);
    });
  });

  test('should allow anonymous users to leave feedback via feedback button', async ({ page }, testInfo) => {
    const isMobile = testInfo.project.name === 'Mobile';
    
    // 1. Go to homepage anonymously
    await page.goto('/');

    // 2. Find and click the feedback button
    const feedbackBtn = page.getByTestId(isMobile ? 'mobile-feedback-btn' : 'ag-nav-feedback-btn');
    await expect(feedbackBtn).toBeVisible();
    await feedbackBtn.click();

    // 3. Fill and submit the feedback modal
    const feedbackModal = page.locator('dialog#feedback-modal');
    await expect(feedbackModal).toBeVisible();

    // Select feedback type (Issue is default, let's select Other)
    await feedbackModal.locator('input#feedback-type-other').check();
    
    // Fill content
    const testMessage = `Anonymous Feedback: ${Date.now()}`;
    await feedbackModal.locator('textarea[name="content"]').fill(testMessage);

    // Submit feedback
    await feedbackModal.locator('button[type="submit"]').click();

    // 4. Verify success screen
    await expect(feedbackModal).toContainText('Thank You!');
    
    // Close modal
    await feedbackModal.locator('button:has-text("Close")').click();
    await expect(feedbackModal).not.toBeVisible();
  });

  test('should allow signed-in users to leave feedback via feedback button', async ({ page }, testInfo) => {
    const isMobile = testInfo.project.name === 'Mobile';

    // 1. Sign in
    await page.goto('/auth/dev?email=user@agbalumo.com');
    await page.waitForURL('/');

    // 2. Find and click the feedback button
    const feedbackBtn = page.getByTestId(isMobile ? 'mobile-feedback-btn' : 'ag-nav-feedback-btn');
    await expect(feedbackBtn).toBeVisible();
    await feedbackBtn.click();

    // 3. Fill and submit the feedback modal
    const feedbackModal = page.locator('dialog#feedback-modal');
    await expect(feedbackModal).toBeVisible();

    await feedbackModal.locator('input#feedback-type-feature').check();
    const testMessage = `User Feedback: ${Date.now()}`;
    await feedbackModal.locator('textarea[name="content"]').fill(testMessage);

    // Submit
    await feedbackModal.locator('button[type="submit"]').click();

    // 4. Verify success
    await expect(feedbackModal).toContainText('Thank You!');
  });

  test('admin should see feedback left by anonymous/signed in users from dashboard', async ({ page }, testInfo) => {
    const isMobile = testInfo.project.name === 'Mobile';
    const timeMark = Date.now();
    const anonymousMessage = `Anon-${timeMark}`;
    const signedInMessage = `User-${timeMark}`;

    // 1. Submit anonymous feedback
    await page.goto('/');
    const feedbackBtn = page.getByTestId(isMobile ? 'mobile-feedback-btn' : 'ag-nav-feedback-btn');
    await expect(feedbackBtn).toBeVisible();
    await feedbackBtn.click();

    let feedbackModal = page.locator('dialog#feedback-modal');
    await expect(feedbackModal).toBeVisible();
    await feedbackModal.locator('input#feedback-type-issue').check();
    await feedbackModal.locator('textarea[name="content"]').fill(anonymousMessage);
    await feedbackModal.locator('button[type="submit"]').click();
    await expect(feedbackModal).toContainText('Thank You!');
    await feedbackModal.locator('button:has-text("Close")').click();
    await expect(feedbackModal).not.toBeVisible();

    // 2. Sign in and submit user feedback
    await page.goto('/auth/dev?email=test-user@agbalumo.com');
    await page.waitForURL('/');
    
    await expect(feedbackBtn).toBeVisible();
    await feedbackBtn.click();

    feedbackModal = page.locator('dialog#feedback-modal');
    await expect(feedbackModal).toBeVisible();
    await feedbackModal.locator('input#feedback-type-feature').check();
    await feedbackModal.locator('textarea[name="content"]').fill(signedInMessage);
    await feedbackModal.locator('button[type="submit"]').click();
    await expect(feedbackModal).toContainText('Thank You!');
    await feedbackModal.locator('button:has-text("Close")').click();
    await expect(feedbackModal).not.toBeVisible();

    // 3. Log out and log in as admin
    await page.goto('/auth/logout');
    await page.waitForURL('/');

    await page.goto('/auth/dev?email=admin-user@agbalumo.com');
    await page.waitForURL('/');

    // Go to admin login
    await page.goto('/admin/login');
    
    // If not redirected to /admin immediately, fill the code and submit
    if (page.url().endsWith('/admin/login')) {
      await page.locator('input[name="code"]').fill('agbalumo2024');
      await page.locator('button[type="submit"]').click();
      await page.waitForURL('/admin');
    }

    // 4. Verify dashboard lists both feedback messages
    const body = page.locator('body');
    await expect(body).toContainText(anonymousMessage);
    await expect(body).toContainText(signedInMessage);
  });
});
