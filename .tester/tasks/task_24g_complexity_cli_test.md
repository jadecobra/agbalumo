# Task 24g: Reduce Complexity — `TestCLIJSONOutput` (score: 14 → target ≤ 10)

## File
`cmd/cli_json_test.go`

## Context

`TestCLIJSONOutput` scores **14** because it contains three subtests, each with an
`if err := rootCmd.Execute(); err != nil` block and custom output parsing/assertions.

## Current Code (lines 14–75)

```go
func TestCLIJSONOutput(t *testing.T) {
    // ... setup ...

    // 1. Test listing list --json (empty)
    t.Run("listing list --json empty", func(t *testing.T) {
        buf := new(bytes.Buffer)
        rootCmd.SetOut(buf)
        rootCmd.SetArgs([]string{"listing", "list"})

        if err := rootCmd.Execute(); err != nil {
            t.Fatalf("Execute failed: %v", err)
        }

        output := strings.TrimSpace(buf.String())
        if !strings.Contains(output, "[]") {
            t.Errorf("Expected output to contain '[]', got %q", output)
        }
    })

    // 2. Test category list --json
    t.Run("category list --json", func(t *testing.T) {
        buf := new(bytes.Buffer)
        // ... same setup & execute block ...
        jsonPart := extractJSONFromOutput(t, buf.String())
        var categories []domain.CategoryData
        if err := json.Unmarshal([]byte(jsonPart), &categories); err != nil {
            t.Fatalf("Unmarshal failed: %v", err)
        }
    })

    // 3. Test listing create --json
    t.Run("listing create --json", func(t *testing.T) {
        buf := new(bytes.Buffer)
        // ... same setup & execute block ...
        jsonPart := extractJSONFromOutput(t, buf.String())
        var listing domain.Listing
        if err := json.Unmarshal([]byte(jsonPart), &listing); err != nil {
            t.Fatalf("Unmarshal failed: %v", err)
        }
        if listing.Title != "JSON Test Listing" {
            t.Errorf("Expected title 'JSON Test Listing', got %q", listing.Title)
        }
    })
}
```

## Required Changes

Extract a helper `executeCommand(t *testing.T, args ...string) string` that
handles the buffer setup, argument injection, and error checking.

```go
func executeCommand(t *testing.T, args ...string) string {
    t.Helper()
    buf := new(bytes.Buffer)
    rootCmd.SetOut(buf)
    rootCmd.SetArgs(args)

    if err := rootCmd.Execute(); err != nil {
        t.Fatalf("Execute failed: %v", err)
    }

    return buf.String()
}
```

Then refactor the three subtests to use it:

```go
func TestCLIJSONOutput(t *testing.T) {
    // ... setup (keep as is) ...

    t.Run("listing list --json empty", func(t *testing.T) {
        output := executeCommand(t, "listing", "list")
        output = strings.TrimSpace(output)
        if !strings.Contains(output, "[]") {
            t.Errorf("Expected output to contain '[]', got %q", output)
        }
    })

    t.Run("category list --json", func(t *testing.T) {
        output := executeCommand(t, "category", "list")
        jsonPart := extractJSONFromOutput(t, output)
        var categories []domain.CategoryData
        if err := json.Unmarshal([]byte(jsonPart), &categories); err != nil {
            t.Fatalf("Unmarshal failed: %v", err)
        }
    })

    t.Run("listing create --json", func(t *testing.T) {
        output := executeCommand(t, "listing", "create", "--title", "JSON Test Listing")
        jsonPart := extractJSONFromOutput(t, output)
        var listing domain.Listing
        if err := json.Unmarshal([]byte(jsonPart), &listing); err != nil {
            t.Fatalf("Unmarshal failed: %v", err)
        }
        if listing.Title != "JSON Test Listing" {
            t.Errorf("Expected title 'JSON Test Listing', got %q", listing.Title)
        }
    })
}
```

## What NOT to change
- Do not modify `extractJSONFromOutput`

## Verification

```bash
go test ./cmd/...
go run cmd/verify/main.go critique 2>&1 | grep "TestCLIJSONOutput"
```

The function must no longer appear, or score ≤ 10.

## Commit
```
test(cmd): extract command execution helper to reduce test complexity
```
