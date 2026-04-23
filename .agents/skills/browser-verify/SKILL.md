---
name: Browser Verification
description: Verify UI changes using browser subagent with proper environment detection
---
# Browser Verification Skill
## Pre-flight (MANDATORY — always run before browser tasks)
1. Read `.agents/invariants.json` — get `protocol` and `port` fields
2. Construct the base URL: `{protocol}://localhost:{port}`
3. Verify server is running: `lsof -i :{port}` — if no output, start server first
4. Read `.env` for `BASE_URL` as override — if set, use that instead
## Verification Checklist
For EVERY UI element verified, you MUST check ALL of:
- [ ] **Exists**: `document.querySelector(selector)` returns non-null
- [ ] **Visible**: `element.offsetHeight > 0 && element.offsetWidth > 0`
- [ ] **Has Content**: `element.innerText.trim().length > 0`
- [ ] **Interactive**: Click/hover produces expected state change
- [ ] **Viewport Safe**: Element fully visible at 1440x900 without scrolling main window
## Common Failure Patterns
| Symptom | Root Cause | Fix |
|---------|-----------|-----|
| Connection refused | Wrong port/protocol | Check `.agents/invariants.json` |
| Element exists but invisible | CSS `display: none` or `opacity: 0` | Check computed styles |
| Click has no effect | Duplicate event listeners | Check `initApp()` in `app.js` |
| Dropdown clipped | Viewport clearance | Add `open-upwards` logic |
| Stale content after fix | Browser cache | Increment `?v=N` in `head_meta.html` |
## Post-flight
1. Document each check result in `task.md` with pass/fail
2. For layout changes: embed screenshot in walkthrough
3. For interactive changes: describe the state transition verified
