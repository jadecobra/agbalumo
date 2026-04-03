```
cmd/           CLI commands (Cobra)
internal/
  agent/       High-governance harness, context cost, and security auditors
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
