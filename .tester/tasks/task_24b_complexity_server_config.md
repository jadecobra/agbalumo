# Task 24b: Reduce Complexity — `ResolveServerConfig` (score: 14 → target ≤ 10)

## File
`cmd/serve.go`

## Context

`ResolveServerConfig` scores **14** due to three levels of nested branching:

1. `if port == ""` → `if appURL != ""` → `if strings.Contains(...)` (3 nested nodes)
2. `if port == ""` → `if hasCerts && env != "production"` (2 more nodes)
3. `if env == "production"` (1 node)
4. `if hasCerts` (1 node)

Total: the port-resolution path alone accounts for 5+ nodes in a 40-line function.

## Current Code (lines 27–68)

```go
func ResolveServerConfig(env, port string, fileExists func(string) bool) ServerConfig {
    certFile := "certs/cert.pem"
    keyFile := "certs/key.pem"
    hasCerts := fileExists(certFile) && fileExists(keyFile)

    if port == "" {
        if appURL := os.Getenv("APP_URL"); appURL != "" {
            if strings.Contains(appURL, ":") {
                parts := strings.Split(appURL, ":")
                port = parts[len(parts)-1]
                port = strings.TrimSuffix(port, "/")
            }
        }
    }
    if port == "" {
        if hasCerts && env != "production" {
            port = "8443"
        } else {
            port = "8080"
        }
    }

    if env == "production" {
        return ServerConfig{Addr: ":" + port, TLS: false}
    }
    if hasCerts {
        return ServerConfig{Addr: ":" + port, TLS: true, CertFile: certFile, KeyFile: keyFile}
    }
    return ServerConfig{Addr: ":" + port, TLS: false}
}
```

## Required Changes

**Step 1**: Extract port resolution into its own function:

```go
// resolvePort returns the effective port to listen on.
// It checks: explicit arg → APP_URL env → cert-aware default.
func resolvePort(port, env string, hasCerts bool) string {
    if port != "" {
        return port
    }
    if appURL := os.Getenv("APP_URL"); appURL != "" {
        if strings.Contains(appURL, ":") {
            parts := strings.Split(appURL, ":")
            port = strings.TrimSuffix(parts[len(parts)-1], "/")
            if port != "" {
                return port
            }
        }
    }
    if hasCerts && env != "production" {
        return "8443"
    }
    return "8080"
}
```

**Step 2**: Slim `ResolveServerConfig` to:

```go
func ResolveServerConfig(env, port string, fileExists func(string) bool) ServerConfig {
    const certFile = "certs/cert.pem"
    const keyFile = "certs/key.pem"
    hasCerts := fileExists(certFile) && fileExists(keyFile)
    port = resolvePort(port, env, hasCerts)

    if env == "production" {
        return ServerConfig{Addr: ":" + port, TLS: false}
    }
    if hasCerts {
        return ServerConfig{Addr: ":" + port, TLS: true, CertFile: certFile, KeyFile: keyFile}
    }
    return ServerConfig{Addr: ":" + port, TLS: false}
}
```

## What NOT to change
- Do not touch `serveCmd`, `init()`, or any test files
- `resolvePort` is unexported — keep it in `cmd/serve.go` alongside `ResolveServerConfig`

## Verification

```bash
go test ./cmd/...
go run cmd/verify/main.go critique 2>&1 | grep "ResolveServerConfig"
```

`ResolveServerConfig` must no longer appear in the cognitive complexity output,
OR its score must be ≤ 10.

## Commit

```
refactor(cmd): extract resolvePort to reduce ResolveServerConfig complexity
```
