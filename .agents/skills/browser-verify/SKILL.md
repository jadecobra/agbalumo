---
name: Browser Verification
description: Verify UI changes using browser subagent with proper environment detection
---

# Browser Verification Skill

## Pre-flight (MANDATORY)

1. Read `.env` for `BASE_URL` — extract protocol and port
2. If no BASE_URL, default to `https://localhost:8443`
3. Verify server is running: `lsof -i :<port>`

## Verification Checklist

For every UI element verified, you MUST check:

- [ ] Element exists (`querySelector`)
- [ ] Element is visible (`offsetHeight > 0`)
- [ ] Element has content (`innerText.length > 0`)
- [ ] Interactive behavior works (click → state change)
- [ ] Viewport clearance at 1440x900

## Post-flight

- Document each check with pass/fail in task.md
- Screenshot proof required for layout changes
