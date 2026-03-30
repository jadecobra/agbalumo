# agbalumo - agbalumo-Squad Guidelines

## Project Overview

agbalumo is a Go web application for the West African diaspora community. It is built and maintained by the **agbalumo-Squad**, a modular multi-agent system designed for high-governance, 10x engineering.

---

## 1. Modular Persona System
The squad's configuration is managed via **[.agents/config.yaml](file:///Users/johnnyblase/gym/agbalumo/.agents/config.yaml)**. This registry maps specialized personas to their respective instruction sets in `.agents/personas/`.

### Core Governance Roles
- **[ProductOwner](file:///Users/johnnyblase/gym/agbalumo/.agents/personas/product_owner.yaml)**: Owns the **"Why"**. Final authority on user value, cultural context, and reversible (Two-Way Door) decisions.
- **[SDET / SecurityEngineer](file:///Users/johnnyblase/gym/agbalumo/.agents/personas/security_engineer.yaml)**: Owns the **"Proof"**. Final authority on quality and security. They write the failing tests (RED) but NEVER the fix. They are the "Witness" to the implementation's success.
- **[BackendEngineer](file:///Users/johnnyblase/gym/agbalumo/.agents/personas/backend_engineer.yaml)**: Owns the **"Fix"**. Responsible for minimal, high-performance Go logic to pass the witness's tests.

---

## 2. Build/Lint/Test Commands
See **[.agents/rules/commands.md](file:///Users/johnnyblase/gym/agbalumo/.agents/rules/commands.md)** for testing, linting, and build instructions.
Mandatory verification via **`scripts/verify-persona.go`** must pass before any squad-level changes are finalized.

---

## 3. Specialized Workflows
See **[.agents/rules/workflows.md](file:///Users/johnnyblase/gym/agbalumo/.agents/rules/workflows.md)** for the available workflow registry and deep-dive guidelines.
All feature work MUST follow the **`/build-feature`** sequence.

---

## 4. Architecture
See **[.agents/rules/architecture.md](file:///Users/johnnyblase/gym/agbalumo/.agents/rules/architecture.md)** for the project's directory structure and component layout.
- **cmd/**: Application entry points.
- **internal/domain/**: Core business logic and interfaces.
- **internal/handler/**: HTTP/Echo handlers and boundary validation.
- **internal/repository/**: SQLite/Persistence logic.
