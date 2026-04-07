# Task 24e: Reduce Complexity — `parseEventDates` (score: 11 → target ≤ 10)

## File
`internal/module/listing/listing_form.go`

## Context

`parseEventDates` scores **11** because it contains a highly regular pattern repeated
twice within a parent `if` statement:
1. `if l.Type == domain.Event`
2. `if req.EventStart != ""`
3. `if err != nil`
4. `if req.EventEnd != ""`
5. `if err != nil`

## Current Code (lines 93–111)

```go
func parseEventDates(req *ListingFormRequest, l *domain.Listing) error {
    if l.Type == domain.Event {
        if req.EventStart != "" {
            parsedTime, err := time.Parse("2006-01-02T15:04", req.EventStart)
            if err != nil {
                return echo.NewHTTPError(http.StatusBadRequest, "Invalid Start Date Format")
            }
            l.EventStart = parsedTime
        }
        if req.EventEnd != "" {
            parsedTime, err := time.Parse("2006-01-02T15:04", req.EventEnd)
            if err != nil {
                return echo.NewHTTPError(http.StatusBadRequest, "Invalid End Date Format")
            }
            l.EventEnd = parsedTime
        }
    }
    return nil
}
```

## Required Changes

**Step 1**: There is already a repetitive pattern with `parseJobStartDate` and `parseDeadline`.
First, extract a simple date parsing helper:

```go
func parseFormDate(val, format, errMsg string) (time.Time, error) {
    if val == "" {
        return time.Time{}, nil
    }
    parsed, err := time.Parse(format, val)
    if err != nil {
        return time.Time{}, echo.NewHTTPError(http.StatusBadRequest, errMsg)
    }
    return parsed, nil
}
```

**Step 2**: Refactor `parseEventDates` to use it:

```go
func parseEventDates(req *ListingFormRequest, l *domain.Listing) error {
    if l.Type != domain.Event {
        return nil
    }
    
    start, err := parseFormDate(req.EventStart, "2006-01-02T15:04", "Invalid Start Date Format")
    if err != nil {
        return err
    }
    if !start.IsZero() {
        l.EventStart = start
    }

    end, err := parseFormDate(req.EventEnd, "2006-01-02T15:04", "Invalid End Date Format")
    if err != nil {
        return err
    }
    if !end.IsZero() {
        l.EventEnd = end
    }

    return nil
}
```

> **Bonus Cleanup (Optional but recommended)**:
> Since you extracted `parseFormDate`, you can also clean up `parseDeadline` and
> `parseJobStartDate` in the same file to use it, dropping their complexity too. Example:
> ```go
> func parseDeadline(req *ListingFormRequest, l *domain.Listing) error {
>     if l.Type == domain.Request {
>         parsed, err := parseFormDate(req.DeadlineDate, "2006-01-02", "Invalid Date Format")
>         if err == nil && !parsed.IsZero() {
>             l.Deadline = parsed
>         }
>         return err
>     }
>     return nil
> }
> ```

## What NOT to change
- Do not modify any other files

## Verification

```bash
go test ./internal/module/listing/...
go run cmd/verify/main.go critique 2>&1 | grep "parseEventDates"
```

The function must no longer appear, or score ≤ 10.

## Commit
```
refactor(listing): extract form date parsing to reduce cognitive complexity
```
