package domain

import (
	"testing"
	"time"
)

func BenchmarkValidate(b *testing.B) {
	l := Listing{
		ID:           "bench-1",
		OwnerOrigin:  "Nigeria",
		Type:         Business,
		Title:        "Benchmark Business",
		City:         "Abuja",
		ContactEmail: "bench@example.com",
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l.Validate()
	}
}

func FuzzListing_Validate(f *testing.F) {
	// Seed with a valid listing
	f.Add("Nigeria", string(Business), "Standard Biz", "A desc", "123 Main St", "test@test.com")
	f.Add("", string(Request), "", "", "", "")
	f.Add("Ghana", string(Food), "Foodie", "Yum", "Street 1", "food@test.com")

	f.Fuzz(func(t *testing.T, origin, cat, title, desc, address, email string) {
		l := Listing{
			OwnerOrigin:  origin,
			Type:         Category(cat),
			Title:        title,
			Description:  desc,
			Address:      address,
			ContactEmail: email,
			CreatedAt:    time.Now(),
		}
		// We just care if it panics
		_ = l.Validate()
	})
}
