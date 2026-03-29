# Skill: User Journey Audit
name: user_journey_audit
description: Simulate real-world user paths via browser to detect functional regressions and UX friction.
---
## Objective
Perform an autonomous, end-to-end audit of critical user journeys (Visitor, Auth User, Admin) to ensure the application remains functional, performant, and brand-compliant after changes.

## Audit Protocol
1. **Restart Server**: Always run `task restart-server` before auditing to ensure the latest binary and assets are active.
2. **Path Simulation**: Execute journeys defined in `.agents/user_journeys.yaml`.
3. **Friction Analysis**:
    - **Layout Shift**: Flag any element shifting >5px during navigation or interaction.
    - **Interaction Lag**: Detect delays >100ms between `click` and `HTMX swap` or visual state change.
    - **Conversion Efficiency**: Count steps to reach "Value" (e.g., Listing Detail). Max allowed: 3.
4. **Brand Guardrails**:
    - Verify `bg-earth-dark` or `#f4f4f4` (toon) persistence.
    - Verify 0px border-radius unless explicitly allowed.
5. **Failure Reporting**:
    - Capture screenshots **ONLY on failure**.
    - Categorize issues as: `Functional Regression`, `UX Friction`, or `Brand Drift`.

## Decision Authority
- **REJECT**: If any journey step fails or "Steps to Value" exceeds threshold.
- **WARN**: If friction metrics (Lag/Shift) are near thresholds but journey succeeds.

## Scripts
- Orchestrator: `scripts/browser_audit.sh`
