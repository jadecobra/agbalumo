# Cmd - CLI Commands

## OVERVIEW
Cobra-based CLI for server management, listing CRUD, seeding, and security auditing.

## WHERE TO LOOK

| Task | File |
|------|------|
| Server start | `serve.go` - `serveCmd`, `serveHTTPS()` |
| Listing CRUD | `listing.go` - `listingCmd` with create/list/get/delete |
| Database seeding | `seed.go` - `seedCmd` |
| Admin operations | `admin.go` - `adminCmd` |
| Security audit | `security-audit/` - standalone security checks |

## CONVENTIONS

```go
// Cobra command structure
var rootCmd = &cobra.Command{Use: "agbalumo"}
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start HTTP server",
    RunE:  func(cmd *cobra.Command, args []string) error { ... },
}

// Flag binding in init()
func init() {
    serveCmd.Flags().StringVar(&flagAddr, "addr", ":8080", "server address")
    rootCmd.AddCommand(serveCmd)
}

// Persistent flags for subcommands
listingCmd.Flags().StringVar(&flagTitle, "title", "", "listing title")
```

## ANTI-PATTERNS

- Do NOT use `flag.Parse()` directly - Cobra handles it
- Do NOT hardcode DB paths - use config or env vars
- Do NOT skip `cmd.Help()` on flag errors
