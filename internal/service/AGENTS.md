# Internal/Service - Business Logic Layer

## Overview
Business logic layer with CSV processing, image handling, background services, and listing management.

## Service Structure

### Basic Pattern
```go
type XxxService struct {
    Repo         domain.ListingStore
    ImageService service.ImageService
    // other dependencies
}

func NewXxxService(repo domain.ListingStore, is service.ImageService) *XxxService {
    return &XxxService{Repo: repo, ImageService: is}
}
```

## CSV Processing

### Import Patterns
```go
// Parse and import CSV data
func (s *CSVService) ParseAndImport(ctx context.Context, r io.Reader) ([]domain.Listing, error) {
    // Parse CSV rows
    // Validate data
    // Convert to domain objects
    // Save to repository
}
```

### Row Validation
```go
// Validate individual CSV row
func (s *CSVService) parseRow(row []string) (domain.Listing, error) {
    // Validate required fields
    // Parse dates and numbers
    // Handle category mapping
}
```

## Image Processing

### Upload Patterns
```go
// Upload and process image
func (s *LocalImageService) UploadImage(file *multipart.FileHeader) (string, error) {
    // Validate file type
    // Resize/compress image
    // Save to storage
    // Return URL
}
```

### Image Operations
```go
// Compress image to reduce size
func (s *LocalImageService) CompressImage(img image.Image, format string) ([]byte, error) {
    // Resize if needed
    // Compress based on format
    // Return bytes
}
```

## Background Services

### Ticker Patterns
```go
// Background service with ticker
func (s *BackgroundService) StartTicker(ctx context.Context) {
    ticker := time.NewTicker(s.CleanupInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.expireListings()
        }
    }
}
```

### Cleanup Operations
```go
// Cleanup expired listings
func (s *BackgroundService) expireListings() error {
    // Find expired listings
    // Delete or mark as expired
    // Return result
}
```

## Listing Management

### Claim Processing
```go
// Process listing claim
func (s *ListingService) ClaimListing(ctx context.Context, listingID string, claimerID string) error {
    // Validate ownership
    // Update listing status
    // Send notifications
}
```

## Error Handling

### Service Errors
```go
// Return domain-specific errors
if err := s.Repo.Save(ctx, listing); err != nil {
    return fmt.Errorf("failed to save listing: %w", err)
}
```

### Validation Errors
```go
// Validate business rules
if err := validateListing(listing); err != nil {
    return fmt.Errorf("invalid listing data: %w", err)
}
```

## Dependencies
- **domain.ListingStore** - Data access
- **service.ImageService** - Image processing
- **context.Context** - Request context
- **io.Reader** - CSV input
- **time.Ticker** - Background processing

## Testing Patterns

### Table-Driven Tests
```go
func TestCSVService(t *testing.T) {
    tests := []struct {
        name     string
        csvData  string
        expectErr bool
    }{}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Mock Usage
```go
// Mock repository for testing
mockRepo := new(mock.Repository)
mockRepo.On("Save", ctx, mock.Anything).Return(nil)
```

## Common Patterns

### Context Usage
```go
// Always use context for operations
ctx := c.Request().Context()
```

### Error Wrapping
```go
// Wrap errors with context
return fmt.Errorf("operation failed: %w", err)
```

## Anti-Patterns

- **No direct database calls** - Use repository
- **No blocking operations** - Use context for cancellation
- **No unhandled errors** - Always wrap and return
- **No hardcoded paths** - Use configuration
- **No synchronous background work** - Use ticker or channels

## Coverage Requirements
- All service methods must have tests
- Error paths must be tested
- CSV parsing edge cases must be tested
- Image processing must be tested
- Background service lifecycle must be tested

## Testing Commands
```bash
go test -json -v -run TestServiceName ./internal/service/
go test -json -v -run TestServiceName/SubtestName ./internal/service/
```

## Notes
- Always use context for cancellation
- Test both success and error paths
- Use table-driven tests for CSV parsing
- Mock external dependencies (image service, repository)
- Test background service start/stop behavior
