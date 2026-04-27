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
		{
			name: "Anchor text variation - Order Online",
			html: `<div><a href="/checkout">Order Online</a></div>`,
			expected: AdaSignals{
				MenuURL: "http://example.com/checkout",
			},
		},
		{
			name: "Anchor text variation - Our Menu",
			html: `<div><a href="/carte">Our Menu</a></div>`,
			expected: AdaSignals{
				MenuURL: "http://example.com/carte",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := scraper.parseHTML(strings.NewReader(tt.html), "http://example.com")
			assertAdaSignalsMatch(t, tt.name, signals, tt.expected)
		})
	}
}

func assertAdaSignalsMatch(t *testing.T, name string, got, want AdaSignals) {
	t.Helper()
	if got.HeatLevel != want.HeatLevel {
		t.Errorf("%s: HeatLevel = %d, want %d", name, got.HeatLevel, want.HeatLevel)
	}
	if got.PaymentMethods != want.PaymentMethods {
		t.Errorf("%s: PaymentMethods = %q, want %q", name, got.PaymentMethods, want.PaymentMethods)
	}
	if got.TopDish != want.TopDish {
		t.Errorf("%s: TopDish = %q, want %q", name, got.TopDish, want.TopDish)
	}
	if got.MenuURL != want.MenuURL {
		t.Errorf("%s: MenuURL = %q, want %q", name, got.MenuURL, want.MenuURL)
	}
}
