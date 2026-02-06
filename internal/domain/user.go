package domain

import "time"

// User represents a registered user (listing owner).
type User struct {
	ID        string    `json:"id"`
	GoogleID  string    `json:"google_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository defines the interface for user persistence.
// Check if context is needed (it is in other repos).
// Using the same pattern as ListingRepository.
/*
type UserRepository interface {
    Save(ctx context.Context, user User) error
    FindByGoogleID(ctx context.Context, googleID string) (User, error)
    FindByID(ctx context.Context, id string) (User, error)
}
*/
