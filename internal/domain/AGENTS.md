# Internal/Domain - Core Domain Types

## Overview
Core business entities and validation logic for West African diaspora platform.

## Domain Structure

### Entity Patterns
```go
type Listing struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Category    string    `json:"category"`
    Type        string    `json:"type"` // "listing", "job", "event"
    Status      string    `json:"status"` // "active", "pending", "rejected"
    OwnerID     string    `json:"owner_id"`
    // ... other fields
}

// Constructor pattern
func NewListing(title string, description string, category string, ownerID string) (*Listing, error) {
    // Validate inputs
    // Initialize fields
    // Return pointer
}
```

## Validation Patterns

### Business Validation
```go
// Validate listing before processing
func (l *Listing) Validate() error {
    if len(l.Title) == 0 {
        return errors.New("title is required")
    }
    if len(l.Description) < 10 {
        return errors.New("description must be at least 10 characters")
    }
    if !isValidCategory(l.Category) {
        return errors.New("invalid category")
    }
    return nil
}
```

### Type-Specific Validation
```go
// Validate job-specific fields
func validateJob(listing *Listing) error {
    if listing.Type != "job" {
        return nil
    }
    if len(listing.Company) == 0 {
        return errors.New("company is required for jobs")
    }
    if listing.Salary == 0 {
        return errors.New("salary is required for jobs")
    }
    return nil
}
```

## Error Types

### Domain Errors
```go
// Custom domain error types
var (
    ErrInvalidTitle       = errors.New("invalid title")
    ErrInvalidDescription = errors.New("invalid description")
    ErrInvalidCategory    = errors.New("invalid category")
    ErrInvalidType        = errors.New("invalid listing type")
)
```

## Conversion Patterns

### Struct Conversion
```go
// Convert between domain and API models
func ToListing(input *ListingInput) (*Listing, error) {
    // Validate input
    // Convert fields
    // Return domain object
}

func FromListing(listing *Listing) *ListingOutput {
    // Convert to API model
    // Format dates
    // Return output
}
```

## Repository Interface

### Domain Store Interface
```go
type ListingStore interface {
    Save(ctx context.Context, listing *Listing) error
    FindByID(ctx context.Context, id string) (*Listing, error)
    FindAll(ctx context.Context) ([]*Listing, error)
    Delete(ctx context.Context, id string) error
    // ... other methods
}
```

## Testing Patterns

### Table-Driven Validation Tests
```go
func TestListingValidation(t *testing.T) {
    tests := []struct {
        name     string
        listing  *Listing
        expectErr bool
    }{
        {name: "valid listing", listing: validListing, expectErr: false},
        {name: "empty title", listing: emptyTitleListing, expectErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.listing.Validate()
            if (err != nil) != tt.expectErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.expectErr)
            }
        })
    }
}
```

## Common Patterns

### Constructor Validation
```go
// Always validate in constructors
func NewListing(title string, description string, category string) (*Listing, error) {
    if len(title) == 0 {
        return nil, ErrInvalidTitle
    }
    // ... other validation
    return &Listing{Title: title, Description: description, Category: category}, nil
}
```

### Type Safety
```go
// Use constants for valid values
const (
    TypeListing = "listing"
    TypeJob     = "job"
    TypeEvent   = "event"
)

// Validate type
func isValidType(t string) bool {
    switch t {
    case TypeListing, TypeJob, TypeEvent:
        return true
    }
    return false
}
```

## Anti-Patterns

- **No business logic in handlers** - Keep in domain
- **No direct database access** - Use repository interfaces
- **No unvalidated input** - Always validate before processing
- **No magic strings** - Use constants for valid values
- **No error suppression** - Always return or handle errors

## Dependencies
- **context.Context** - Request context
- **errors** - Error handling
- **time** - Date/time operations
- **validation** - Input validation

## Coverage Requirements
- All validation functions must have tests
- Error cases must be tested
- Type-specific validation must be tested
- Constructor validation must be tested

## Testing Commands
```bash
go test -v -run TestDomainName ./internal/domain/
go test -v -run TestDomainName/SubtestName ./internal/domain/
```

## Notes
- Keep domain logic pure (no external dependencies)
- Use table-driven tests for validation
- Test both success and error paths
- Validate all input before processing
- Use constructor patterns for object creation
