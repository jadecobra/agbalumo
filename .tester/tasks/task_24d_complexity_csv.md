# Task 24d: Reduce Complexity — `(*CSVService).parseRow` (score: 13 → target ≤ 10)

## File
`internal/service/csv.go`

## Context

`parseRow` scores **13** due to the number of sequential `if` checks across different
types of validation logic (title requires, description requires, contact method requires,
origin fallback, and geocoding fallback). Each adds structural complexity.

## Current Code (lines 212–253)

```go
func (s *CSVService) parseRow(record []string, headerMap map[string]int) (*domain.Listing, error) {
    get := func(col string) string { /* ... */ }

    title := get("title")
    if title == "" { return nil, fmt.Errorf("title is required") }

    desc := get("description")
    if desc == "" { return nil, fmt.Errorf("description is required") }

    email, phone, whatsapp, website := get("email"), get("phone"), get("whatsapp"), get("website")
    if email == "" && phone == "" && whatsapp == "" && website == "" {
        return nil, fmt.Errorf("at least one contact method ... is required")
    }

    origin := get("origin")
    if origin == "" { origin = "Nigeria" }

    address, city := get("address"), get("city")
    if city == "" && address != "" && s.Geocoding != nil {
        if foundCity, err := s.Geocoding.GetCity(context.Background(), address); err == nil && foundCity != "" {
            city = foundCity
        }
    }

    return &domain.Listing{ /* ... */ }, nil
}
```

## Required Changes

**Step 1**: Extract the geocoding logic into a helper `resolveCity(s *CSVService, city, address string) string`:

```go
func resolveCity(s *CSVService, city, address string) string {
    if city != "" || address == "" || s.Geocoding == nil {
        return city
    }
    if foundCity, err := s.Geocoding.GetCity(context.Background(), address); err == nil && foundCity != "" {
        return foundCity
    }
    return city
}
```

**Step 2**: Extract the validation logic into `validateParsedRow(title, desc, email, phone, whatsapp, website string) error`:

```go
func validateParsedRow(title, desc, email, phone, whatsapp, website string) error {
    if title == "" {
        return fmt.Errorf("title is required")
    }
    if desc == "" {
        return fmt.Errorf("description is required")
    }
    if email == "" && phone == "" && whatsapp == "" && website == "" {
        return fmt.Errorf("at least one contact method (email, phone, whatsapp, or website) is required")
    }
    return nil
}
```

**Step 3**: Simplify `parseRow`:

```go
func (s *CSVService) parseRow(record []string, headerMap map[string]int) (*domain.Listing, error) {
    get := func(col string) string {
        if idx, ok := headerMap[col]; ok && idx < len(record) {
            return strings.TrimSpace(record[idx])
        }
        return ""
    }

    title, desc := get("title"), get("description")
    email, phone, whatsapp, website := get("email"), get("phone"), get("whatsapp"), get("website")
    
    if err := validateParsedRow(title, desc, email, phone, whatsapp, website); err != nil {
        return nil, err
    }

    origin := get("origin")
    if origin == "" {
        origin = "Nigeria"
    }

    city := resolveCity(s, get("city"), get("address"))

    return &domain.Listing{
        ID: uuid.New().String(), Title: title, Type: parseCategory(get("type")),
        Description: desc, OwnerOrigin: origin, ContactEmail: email, WebsiteURL: website,
        ContactPhone: phone, ContactWhatsApp: whatsapp, Address: get("address"), City: city,
        HoursOfOperation: get("hours"), CreatedAt: time.Now(),
    }, nil
}
```

## What NOT to change
- Keep the `get` closure inside `parseRow`
- Do not modify any tests
- Do not export the helpers

## Verification

```bash
go test ./internal/service/...
go run cmd/verify/main.go critique 2>&1 | grep "parseRow"
```

The function must no longer appear, or score ≤ 10.

## Commit
```
refactor(service): extract CSV parseRow helpers to reduce cognitive complexity
```
