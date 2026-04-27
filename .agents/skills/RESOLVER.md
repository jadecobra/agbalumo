# Skill Resolver

Read this file at session start. Match intent against triggers. Read the skill file BEFORE acting.

## Workflow Commands
| Trigger | Skill |
|---------|-------|
| `/build-feature` | `.agents/workflows/build-feature.md` |
| `/learn` | `.agents/workflows/learn.md` |
| `/coding-standards` | `.agents/workflows/coding-standards.md` |
| `/stress-test` | `.agents/workflows/stress-test.md` |
| `/deploy-secrets` | `.agents/workflows/deploy-secrets.md` |
| `/skillify` | `.agents/workflows/skillify.md` |

## Procedural Skills
| Trigger | Skill |
|---------|-------|
| Writing tests, fixing bugs, implementing features, TDD | `.agents/skills/go-tdd/SKILL.md` |
| UI change, browser verification, layout check, viewport audit | `.agents/skills/browser-verify/SKILL.md` |
| Pushing changes, CI failure, production parity | `.agents/skills/ci-parity/SKILL.md` |


## Disambiguation
1. Slash command → Workflow Commands table.
2. Modifying `*_test.go` or user says "test" → `go-tdd`.
3. Modifying templates/CSS/JS or user says "UI"/"layout" → `browser-verify`.
4. Both apply → read BOTH skills.
5. Uncertain → ask user.
