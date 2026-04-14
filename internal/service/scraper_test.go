package service

import (
	"strings"
	"testing"
)

func TestWebsiteScraper_Heuristics(t *testing.T) {
	scraper := &WebsiteScraper{}

	tests := []struct {
		name     string
		html     string
		expected AdaSignals
	}{
		{
			name: "Spicy and Payment Keywords",
			html: `<html><body>
				<h1>Our Signature Jollof</h1>
				<p>It is extra spicy and hot!</p>
				<p>We accept Zelle and Venmo.</p>
				<a href="/downloads/menu.pdf">View Menu</a>
			</body></html>`,
			expected: AdaSignals{
				HeatLevel:      2, // "spicy", "hot"
				PaymentMethods: "Zelle, Venmo",
				TopDish:        "Our Signature Jollof",
				MenuURL:        "http://example.com/downloads/menu.pdf",
			},
		},
		{
			name: "Multiple spicy keywords",
			html: `<p>Scotch bonnet, habanero, and chili are used in our scotch bonnet sauce.</p>`,
			expected: AdaSignals{
				HeatLevel: 4, // "scotch bonnet" (twice), "habanero", "chili"
			},
		},
		{
			name: "CashApp and Menu link",
			html: `<div>Pay with cash app</div><a href="https://other.com/menu.jpg">The Menu</a>`,
			expected: AdaSignals{
				PaymentMethods: "CashApp",
				MenuURL:        "https://other.com/menu.jpg",
			},
		},
		{
			name: "No Signals",
			html: `<html><body>Welcome to our restaurant</body></html>`,
			expected: AdaSignals{
				HeatLevel:      0,
				PaymentMethods: "",
				TopDish:        "",
				MenuURL:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := scraper.parseHTML(strings.NewReader(tt.html), "http://example.com")
			if signals.HeatLevel != tt.expected.HeatLevel {
				t.Errorf("%s: HeatLevel = %d, want %d", tt.name, signals.HeatLevel, tt.expected.HeatLevel)
			}
			if signals.PaymentMethods != tt.expected.PaymentMethods {
				t.Errorf("%s: PaymentMethods = %q, want %q", tt.name, signals.PaymentMethods, tt.expected.PaymentMethods)
			}
			if signals.TopDish != tt.expected.TopDish {
				t.Errorf("%s: TopDish = %q, want %q", tt.name, signals.TopDish, tt.expected.TopDish)
			}
			if signals.MenuURL != tt.expected.MenuURL {
				t.Errorf("%s: MenuURL = %q, want %q", tt.name, signals.MenuURL, tt.expected.MenuURL)
			}
		})
	}
}
