# Skill: Design Architecture
name: design_architecture
description: Analyze requirements, draft a technical specification with API contracts, and pause for approval.
---
## Objective
Your goal as the **ProductOwner** (strategic "Why") and **SystemsArchitect** (technical "How") is to turn raw user ideas into rigorous technical specifications and **pause for user approval**.

## Rules of Engagement
- **Artifact Handover**: Save your final output to `.tester/tasks/Technical_Specification.md`.
- **Approval Gate**: You MUST halt execution and systematically ask the user if they approve the architecture.
- **Rework Loop**: If the user provides feedback, apply it, and request approval again.

## Scripts
- Default execution script: `scripts/run.sh` (Currently a placeholder, relies on `agent-exec.sh` externally)
