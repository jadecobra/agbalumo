import { test, expect } from '@playwright/test';

test.describe('Create Listing Sensory Signals Toggling', () => {
  test.beforeEach(async ({ page }) => {
    // Login as developer to allow opening the post modal
    await page.goto('/auth/dev');
    await page.waitForURL('/');

    // Wait for JS initialization
    await page.waitForFunction(() => typeof (window as any).filterState !== 'undefined', { timeout: 10000 });
  });

  test('should hide sensory signals bar when type is Business, show when Food is selected', async ({ page }) => {
    // Open the post modal
    const postBtn = page.locator('[data-testid="ag-nav-post-btn-desktop"], [data-testid="ag-nav-post-btn-mobile"]').filter({ visible: true }).first();
    await expect(postBtn).toBeVisible();
    await postBtn.click();

    const modal = page.locator('dialog[id="create-listing-modal"]');
    await expect(modal).toBeVisible();

    // Verify initial state: dropdown shows Food/Restaurant (the default)
    const dropdownBtn = modal.locator('#listing-type-btn');
    await expect(dropdownBtn).toBeVisible();
    await expect(dropdownBtn).toHaveText(/Food\/Restaurant/i);

    // Verify that the sensory signals container (the gray bar) is VISIBLE initially
    const adaSignalsContainer = modal.locator('[data-agent-template="listing_form_ada_signals"]');
    await expect(adaSignalsContainer).toBeVisible();

    // Verify that both fields (Heat Level and Regional Specialty) are visible and enabled
    const heatInput = modal.locator('input[name="heat_level"]');
    const regionalInput = modal.locator('input[name="regional_specialty"]');
    await expect(heatInput).toBeVisible();
    await expect(heatInput).not.toBeDisabled();
    await expect(regionalInput).toBeVisible();
    await expect(regionalInput).not.toBeDisabled();

    // Click the category dropdown to open it
    await dropdownBtn.click();

    // Select Business category
    const businessOption = modal.locator('[data-dropdown-value="Business"]');
    await expect(businessOption).toBeVisible();
    await businessOption.click();

    // Dropdown text should update to Business
    await expect(dropdownBtn).toHaveText(/Business/i);

    // Verify container is hidden
    await expect(adaSignalsContainer).toBeHidden();

    // Change back to Food/Restaurant
    await dropdownBtn.click();
    const foodOption = modal.locator('[data-dropdown-value="Food"]');
    await expect(foodOption).toBeVisible();
    await foodOption.click();

    // Verify container is visible again
    await expect(adaSignalsContainer).toBeVisible();
  });

  test('should change description placeholder depending on category selection', async ({ page }) => {
    // Open the post modal
    const postBtn = page.locator('[data-testid="ag-nav-post-btn-desktop"], [data-testid="ag-nav-post-btn-mobile"]').filter({ visible: true }).first();
    await expect(postBtn).toBeVisible();
    await postBtn.click();

    const modal = page.locator('dialog[id="create-listing-modal"]');
    await expect(modal).toBeVisible();

    const descriptionInput = modal.locator('textarea[name="description"]');
    await expect(descriptionInput).toBeVisible();

    // Default selection: Food -> "Tell us about your restaurant..."
    await expect(descriptionInput).toHaveAttribute('placeholder', 'Tell us about your restaurant...');

    const dropdownBtn = modal.locator('#listing-type-btn');

    // Change to Business -> "Tell us about your business..."
    await dropdownBtn.click();
    const businessOption = modal.locator('[data-dropdown-value="Business"]');
    await businessOption.click();
    await expect(descriptionInput).toHaveAttribute('placeholder', 'Tell us about your business...');

    // Change to Event -> "Tell us about your event..."
    await dropdownBtn.click();
    const eventOption = modal.locator('[data-dropdown-value="Event"]');
    await eventOption.click();
    await expect(descriptionInput).toHaveAttribute('placeholder', 'Tell us about your event...');
  });
});
