---
description: Workspace cleanup — remove entropy, archive stale docs, surface tech debt markers
---

# /janitor Workflow

**Executor**: SystemsArchitect (run this in your own window — no new persona, no separate handoff)
**Trigger**: After every 5 REFACTOR phase completions, or when `critique_report.md` lists > 10 unresolved items.

---

## Step 1 — Find Stale Docs

```bash
find . -name "*.md" -newer .agents/state.json -not -path "./.git/*"
```
Review each result. If a doc was not updated since the last feature commit, either update or archive it to `progress_archive.md`.

## Step 2 — Surface Tech Debt Markers

```bash
grep -rn "// TODO\|// FIXME\|// HACK" --include="*.go" . | grep -v "_test.go" | sort
```
Log all results in `.tester/tasks/tech_debt.md`. Each entry must have: file, line, description, owner persona.

## Step 3 — Detect Stale HANDOFF

```bash
test -f .tester/tasks/HANDOFF.md && git log --oneline -1 .tester/tasks/HANDOFF.md
```
If `HANDOFF.md` exists and its last git commit predates the last feature commit, delete it:
```bash
rm .tester/tasks/HANDOFF.md
```

## Step 4 — Archive Resolved Progress Items

Run:
```bash
./scripts/agent-exec.sh status --text
```
Move all `[x]` completed categories older than 2 features from `progress.md` into `progress_archive.md`.

## Step 5 — Lint and Confirm Clean State

```bash
task lint
go build ./...
```
Both must pass before Janitor work is considered done.

## Step 6 — Commit

```bash
git add -A
git commit -m "chore: janitor cleanup — stale docs archived, tech debt logged"
```

> [!NOTE]
> Janitor never deletes business logic or test files. When in doubt, archive — don't delete.
