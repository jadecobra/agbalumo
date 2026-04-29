---
name: "UI Cohesion Guard"
description: "Enforce visual consistency, legibility, and theme parity across templates"
triggers:
  - "template_change"
  - "ui_cohesion"
  - "design review"
  - "visual audit"
mutating: false
---

# UI Cohesion Guard
## When to Use
Run this skill on ANY template or CSS modification. It combines deterministic tooling
with a manual checklist to prevent the visual fragmentation documented in
ADR `2026-04-28-surface-theme-unification.md`.
## Step 1: Deterministic Gate (MANDATORY)
Run `go run ./cmd/verify design` — catches:
- Font sizes below 10px (`text-[8px]`, `text-[9px]`)
- Low-contrast opacity (`text-text-sub/60`)
- Hardcoded dark backgrounds in modals bypassing theme sync
- Rounding violations (existing)
- Hardcoded hex codes (existing)
**If violations exist, fix them before proceeding.**
## Step 2: Card ↔ Modal Parity Check
For any change to `listing_card.html` or `modal_detail.html`:
- [ ] Card surface uses `bg-white dark:bg-surface-dark`
- [ ] Modal scrollable content uses `bg-white dark:bg-surface-dark`
- [ ] Modal footer uses `bg-white dark:bg-surface-dark`
- [ ] Text colors use `text-text-main dark:text-earth-cream` (not hardcoded `text-earth-cream` alone)
- [ ] Borders use `border-stone-200 dark:border-stone-800` (not `border-white/10` alone)
## Step 3: Badge Density Audit
For any change to card or modal header areas:
- [ ] Card header shows ≤3 metadata items (Type + Rating + Title)
- [ ] All additional metadata (TopDish, RegionalSpecialty, Origin, HeatLevel) is in the card body or modal only
- [ ] No badge uses font size below `text-[10px]`
## Step 4: Typography Hierarchy Check
- [ ] All `h1`/`h2` use `font-serif` (Playfair Display)
- [ ] All functional text uses `font-sans`/`font-display` (Inter)
- [ ] No `uppercase tracking-[0.2em] font-bold` is applied to more than 2 elements per visible section
- [ ] **The Attention Budget**: There MUST be a maximum of ONE primary high-contrast CTA (e.g., `bg-earth-ochre`) per viewport. Any competing primary buttons must be autonomously demoted to secondary ghost buttons (`bg-transparent border border-earth-ochre`).

## Step 5: Browser Verification (if layout changed)
Follow `.agents/skills/browser-verify/SKILL.md` — verify at all mandatory viewports.

## Step 6: Automated Fix Loop
If a visual regression or design violation is found during Steps 1 or 5, you MUST NOT stop at reporting the error. Instead, you must automatically apply the minimal CSS/Tailwind fix, commit it atomically using the `style(design): <fix description>` conventional format, and capture an "After" screenshot to prove the resolution.