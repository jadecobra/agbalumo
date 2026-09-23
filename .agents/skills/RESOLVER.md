# Skill Resolver

Read this file at session start. Match intent against triggers. Read the skill file BEFORE acting.

## Workflow Commands
| Trigger | Skill |
|---------|-------|
| `/build-feature` | `.agents/skills/build-feature/SKILL.md` |
| `/learn` | `.agents/skills/learn/SKILL.md` |
| `/coding-standards` | `.agents/coding-standards.md` |
| `/stress-test` | `.agents/skills/stress-test/SKILL.md` |
| `/deploy-secrets` | `.agents/skills/deploy-secrets/SKILL.md` |
| `/skill-audit` | `.agents/skills/verify-authoring/SKILL.md` |
| `/refactor` | `.agents/skills/go-tdd/SKILL.md` |
| `/doc-prune` | `.agents/skills/doc-prune/SKILL.md` |
| `/debug` | `.agents/skills/go-tdd/SKILL.md` |
| `/hotfix` | `.agents/skills/go-tdd/SKILL.md` |
| `/red-team`, `/challenge` | `.agents/skills/red-team/SKILL.md` |

## Procedural Skills
| Trigger | Skill |
|---------|-------|
| Writing tests, fixing bugs, implementing features, TDD, /hotfix, /refactor, /debug | `.agents/skills/go-tdd/SKILL.md` |
| UI change, browser verification, layout check, viewport audit | `.agents/skills/browser-verify/SKILL.md` |
| Push changes, CI failure, production parity | `.agents/skills/ci-parity/SKILL.md` |
| /plan, /architect, let's plan, plan for flash, break this down, split into prompts, decompose, flash prompt, design for | `.agents/skills/flash-plan/SKILL.md` |
| /design-critique, critique design, review ui, harsh review, redesigns, modals, overlays, design variants, theme harmonization | `.agents/skills/design-critique/SKILL.md` |
| review flash output, check implementation, verify flash changes | `.agents/skills/flash-review/SKILL.md` |
| add verify subcommand, new verify tool, automate this check, /skill-audit | `.agents/skills/verify-authoring/SKILL.md` |
| audit codebase, health check, score the codebase, review infrastructure, how healthy is the codebase | `.agents/skills/codebase-audit/SKILL.md` |
| migrate handler, typed viewmodel, fix deprecated map, viewmodel migration | `.agents/skills/viewmodel-migration/SKILL.md` |
| asynchronous task, background command, polling, sleep, wait | `.agents/skills/turn-cost/SKILL.md` |
| /build-feature, build feature, implement feature, new feature | `.agents/skills/build-feature/SKILL.md` |
| /learn, learn, codify lesson, record mistake | `.agents/skills/learn/SKILL.md` |
| /red-team, /challenge, red team, challenge idea | `.agents/skills/red-team/SKILL.md` |
| /doc-prune, doc prune, prune documentation, prune docs | `.agents/skills/doc-prune/SKILL.md` |
| /deploy-secrets, deploy secrets, rotate keys, push secrets | `.agents/skills/deploy-secrets/SKILL.md` |
| /stress-test, stress test, benchmark system, load test, performance optimization | `.agents/skills/stress-test/SKILL.md` |




## Disambiguation
1. Slash command → Workflow Commands table.
2. Modifying `*_test.go` or user says "test" → `go-tdd`.
3. Modifying templates/CSS/JS or user says "UI"/"layout" → `browser-verify`.
4. Both apply → read BOTH skills.
5. User says "plan the tests" or "plan the feature" → ask: "Do you want to plan with an expensive model (flash-plan) or execute directly (go-tdd)?"
6. Uncertain → ask user.
