# Goal: Obliterate Localized Agent Constraints (Anti-Pattern Cleanup)

## Background
The repository contains several localized `AGENTS.md` files nested deep within the business layers (e.g., `internal/domain`, `internal/service`). 

As identified in the recent architectural critique, these files violate the "No Paperwork" rule and pollute the context window (wasting Tokens and reducing Agentic Attention). The Hexagonal Architecture is already governed by the overarching rules in the project root's `AGENTS.md`, and strict Go-native linting completely replaces the need for markdown anti-pattern reminders. 

Furthermore, context details (such as the `auth_session` string or `PRAGMA` SQLite tuning settings) are already safely self-documented inside `./internal/middleware/session.go` and `./internal/repository/sqlite/sqlite.go`. 

## Implementation Steps for Gemini 3 Flash

1. **Delete Nested `AGENTS.md` files:**
   Run the following command to permanently remove these redundant files from the repository:
   ```bash
   git rm cmd/AGENTS.md \
          internal/domain/AGENTS.md \
          internal/service/AGENTS.md \
          internal/middleware/AGENTS.md \
          internal/repository/sqlite/AGENTS.md
   ```

2. **Verify Deletion:**
   Ensure no stray `AGENTS.md` files exists other than the `AGENTS.md` file located at the project root. You can run the following to verify:
   ```bash
   find . -type f -name "AGENTS.md" | grep -v "^./AGENTS.md$"
   # This should return exactly no output.
   ```

3. **Commit:**
   Execute the atomic commit to finalize this cleanup:
   ```bash
   git commit -m "docs(agents): obliterate redundant localized agent constraints"
   ```

4. **Sanity Check:**
   Run the native CI pipeline one final time to ensure no architectural dependencies were somehow linked to these text files (this is extremely unlikely, but part of standard protocol):
   ```bash
   go run cmd/verify/main.go ci
   ```
