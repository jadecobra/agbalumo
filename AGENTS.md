# agbalumo - Agent Guidelines

## Project Overview

agbalumo is a Go web application for the West African diaspora community, featuring a business directory, job board, event listings, and community requests. Built with Echo framework, SQLite, HTMX, and Tailwind CSS.

---

## Build/Lint/Test Commands

### Run All Tests
```bash
go test -json ./...
```

### Run Single Test
```bash
go test -json -v -run TestFunctionName ./internal/package/
go test -json -v -run TestFunctionName/SubtestName ./internal/package/
```

### Run Tests with Race Detection
```bash
go test -json -race ./...
```

### Pre-Commit Quality Gate
```bash
./scripts/pre-commit.sh
```
Runs: `gofmt`, `go mod tidy`, `go vet`, race tests, coverage (threshold in `.agent/coverage-threshold`), secret scanning.

### Build & Restart Server
```bash
./scripts/verify_restart.sh
```
Builds CSS, compiles binary, restarts server on :8443.

### Build CSS
```bash
npm run build:css
```

---

## Specialized Workflows
See `.agent/workflows/` for deep-dive guidelines:
- `/coding-standards`: Code Style Guidelines, Naming, Error Handling, Testing Structure
- `/feature-implementation`: Feature Development Workflow
- `/audit`: Performance, Auth, Security gates
- `/restart-server`: Commands to rebuild CSS and binary

---

## Architecture

```
cmd/           CLI commands (Cobra)
internal/
  config/      Configuration
  domain/      Core types, interfaces, business rules
  handler/     HTTP handlers (Echo)
  middleware/  Auth, sessions, rate limiting
  mock/        Test mocks
  repository/  Data access interfaces
  service/     Business logic layer
  ui/          Template renderer
ui/
  templates/   HTML templates (Go templates)
  static/      CSS, JS, images
```
