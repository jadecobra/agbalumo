from playwright.sync_api import sync_playwright

def test():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page(ignore_https_errors=True)
        page.goto("https://127.0.0.1:8443")
        page.wait_for_timeout(1000)
        
        # Click post
        page.locator("button:has-text('POST')").first.click()
        page.wait_for_selector("#create-listing-modal")
        
        # Select Job
        page.locator("#create-listing-modal .dropdown-display").first.click()
        page.wait_for_timeout(500)
        page.locator("text=Job").click()
        page.wait_for_timeout(500)
        
        # Check if job fields are visible
        is_company_visible = page.locator("input[name='company']").is_visible()
        is_photo_visible = page.locator("#image-upload-section").is_visible()
        
        print(f"Company Field Visible: {is_company_visible}")
        print(f"Photo Field Visible: {is_photo_visible}")
        browser.close()

test()
