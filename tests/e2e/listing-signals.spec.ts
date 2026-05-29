import { test, expect } from '@playwright/test';

test.describe('Create Listing Sensory Signals Toggling', () => {
  test.beforeEach(async ({ page }) => {
    // Dismiss location permission prompt automatically
    await page.addInitScript(() => {
      window.sessionStorage.setItem('agbalumo_geo_dismissed', 'true');
    });

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

    // Verify initial state: dropdown shows Business (the default)
    const dropdownBtn = modal.locator('#listing-type-btn');
    await expect(dropdownBtn).toBeVisible();
    await expect(dropdownBtn).toHaveText(/Business/);

    // Verify that the sensory signals container (the gray bar) is HIDDEN initially
    const adaSignalsContainer = modal.locator('[data-agent-template="listing_form_ada_signals"]');
    await expect(adaSignalsContainer).toBeHidden();

    // Click the category dropdown to open it
    await dropdownBtn.click();

    // Select Food category
    const foodOption = modal.locator('[data-dropdown-value="Food"]');
    await expect(foodOption).toBeVisible();
    await foodOption.click();

    // Dropdown text should update to Food
    await expect(dropdownBtn).toHaveText(/Food/);

    // Verify that the sensory signals container (the gray bar) is now VISIBLE
    await expect(adaSignalsContainer).toBeVisible();

    // Verify that both fields (Heat Level and Regional Specialty) are visible and enabled
    const heatInput = modal.locator('input[name="heat_level"]');
    const regionalInput = modal.locator('input[name="regional_specialty"]');
    await expect(heatInput).toBeVisible();
    await expect(heatInput).not.toBeDisabled();
    await expect(regionalInput).toBeVisible();
    await expect(regionalInput).not.toBeDisabled();

    // Change back to Business
    await dropdownBtn.click();
    const businessOption = modal.locator('[data-dropdown-value="Business"]');
    await expect(businessOption).toBeVisible();
    await businessOption.click();

    // Verify container is hidden again
    await expect(adaSignalsContainer).toBeHidden();
  });
});
