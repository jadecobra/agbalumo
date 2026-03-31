# Skill: Design Architecture
name: design_architecture
description: Analyze requirements, draft a technical specification with API contracts, and pause for approval.
---
## Objective
Your goal as the **ProductOwner** (strategic "Why") and **SystemsArchitect** (technical "How") is to turn raw user ideas into rigorous technical specifications and **pause for user approval**.

## Mandatory Sections for `implementation_plan.md`
Every architecture draft MUST contain the following sections:

1. **Target User Avatar**: Define exactly who in the West African diaspora community this feature serves (e.g., "The First-Gen Student", "The Remittance Sender").
2. **Pain Point Mapping**: Explicitly list the specific user pain points addressed (e.g., "High transaction fees", "Community silos").
3. **Strategic Critique**: The ProductOwner MUST provide a pushback summary: **Why is this not "dumb"?** Explain why a simpler solution won't work and how this provides genuine utility.
4. **Technical Contract**: 
   - Mandatory Go interfaces (`internal/domain`).
   - Mandatory DB schema changes (if any).
   - Mandatory JSON API schemas.
5. **Security STRIDE**: A security boundary analysis identifying potential Spoofing, Tampering, Repudiation, Information Disclosure, Denial of Service, and Elevation of Privilege risks.

## Rules of Engagement
- **Artifact Handover**: Save your final output to `implementation_plan.md`.
- **ChiefCritic Gate**: Your plan will be audited by the ChiefCritic for "Programmer Art" and lack of depth.
- **Approval Gate**: You MUST halt execution and systematically ask the user if they approve the architecture.
- **Rework Loop**: If the user provides feedback, apply it, and request approval again.

## Scripts
- Default execution script: `scripts/run.sh`
- Validation script: `scripts/critic-gate.sh`
