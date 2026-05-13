package maintenance

import (
	"encoding/json"
	"fmt"

	"github.com/playwright-community/playwright-go"
)

// HitboxViolation represents a UI interaction violation.
type HitboxViolation struct {
	Tag       string  `json:"tag"`
	Text      string  `json:"text"`
	Reason    string  `json:"reason"`
	BlockedBy string  `json:"blocked_by,omitempty"`
	Width     float64 `json:"width,omitempty"`
	Height    float64 `json:"height,omitempty"`
}

// CheckHitboxes performs a deterministic interaction-layer verifier to prevent mobile filter interception bugs.
func CheckHitboxes(targetURL string) ([]HitboxViolation, error) {
	// Ensure playwright is installed
	err := playwright.Install()
	if err != nil {
		return nil, fmt.Errorf("could not install playwright: %w", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start playwright: %w", err)
	}
	defer func() { _ = pw.Stop() }()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}
	defer func() { _ = browser.Close() }()

	// Mobile viewport audit (44x44 rule is specifically for touch targets)
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Viewport:          &playwright.Size{Width: 375, Height: 812},
		IgnoreHttpsErrors: playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("could not create browser context: %w", err)
	}
	defer func() { _ = context.Close() }()

	page, err := context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %w", err)
	}

	if _, err = page.Goto(targetURL); err != nil {
		return nil, fmt.Errorf("could not navigate to %s: %w", targetURL, err)
	}

	// Inject evaluation script
	const script = `
		(() => {
			const interactive = document.querySelectorAll('button, a, [role="button"]');
			const violations = [];
			for (const el of interactive) {
				const rect = el.getBoundingClientRect();
				const isSmall = rect.width < 44 || rect.height < 44;
				
				// Verify if blocked by transparent overlay
				const centerX = rect.left + rect.width / 2;
				const centerY = rect.top + rect.height / 2;
				const topEl = document.elementFromPoint(centerX, centerY);
				const isBlocked = topEl && !el.contains(topEl) && !topEl.contains(el);
				
				if (isSmall || isBlocked) {
					violations.push({
						tag: el.tagName,
						text: el.innerText.trim().substring(0, 30),
						reason: isSmall ? "Touch target too small (min 44x44)" : "Interaction blocked by overlay",
						width: rect.width,
						height: rect.height,
						blocked_by: isBlocked ? (topEl.tagName + (topEl.className ? "." + topEl.className.split(" ").join(".") : "")) : ""
					});
				}
			}
			return violations;
		})()
	`

	raw, err := page.Evaluate(script)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate hitbox script: %w", err)
	}

	// Result from Evaluate is an interface{}, we need to convert it to our struct
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal raw violations: %w", err)
	}

	var violations []HitboxViolation
	if err := json.Unmarshal(data, &violations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal violations: %w", err)
	}

	return violations, nil
}
