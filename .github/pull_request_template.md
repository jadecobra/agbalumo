# Pull Request Checklist

## Summary
- **Type**: (Choose one: Feat, Fix, Refactor, Perf, Docs, Chore)
- **Problem**: What issue does this PR solve?
- **Solution**: How was it solved?
- **Impact**: What areas might be affected?

## Verification (SDET Approval)
- [ ] **Automated Tests**: All unit/integration tests pass (`go test ./...`).
- [ ] **Test Coverage**: No regression in coverage (>80% required).
- [ ] **Linting**: Code is formatted and linted (`go fmt`, `golangci-lint`).
- [ ] **Security**: No secrets committed. Input validation confirmed.

## User Experience (UI/UX)
- [ ] **Delight**: Visual feedback confirmed (hover states, animations).
- [ ] **Responsiveness**: Tested on Mobile/Desktop breakpoints.
- [ ] **Performance**: Fast load times verified.
- [ ] **Screenshots/Video**: Attached for visual changes.

## Compliance
- [ ] **Documentation**: Updated `task.md`, `README.md`, or `CONTRIBUTING.md` if needed.
- [ ] **Standards**: Follows `CONTRIBUTING.md` guidelines (10x Protocol).

---
**Note:** Please tag relevant reviewers. Do not merge until verification is complete.
