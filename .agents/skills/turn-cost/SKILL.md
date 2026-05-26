---
name: Turn-Cost Optimization
description: Optimize LLM turn-cost and context density using reactive background execution and fail-safe timers.
triggers:
  - "asynchronous task"
  - "background command"
  - "polling"
  - "sleep"
  - "wait"
mutating: false
---

# Turn-Cost & Quota Optimization Skill

## Core Concept
Every intermediate turn (waiting, sleeping, status polling) degrades agentic attention and inflates context density (`TokenRMS`), resulting in higher API latency and cost. This skill guides the agent in optimizing turns using reactive background execution and optimal fail-safe timers.

## Step 1: The Turn-Cost Calculus
Before waiting or scheduling status checks, evaluate the process against the **Determinism Axiom**:
* **Deterministic Processes (CI compilation, unit tests, remote deploys):**
  - **Do NOT Poll:** Active polling is considered Quota Bloat and is STRICTLY FORBIDDEN.
  - **Action:** Launch the task in the background using `run_command` (e.g. `go run ./cmd/verify precommit &`), set a single conservative fail-safe timer via `schedule`, and immediately yield the turn.
  - **Failsafe Metrics:**
    - Local verification/compilation: **90 seconds**
    - Remote CI/CD builds & deploys: **300 seconds (5 minutes)**
* **Non-Deterministic / Interactive Processes (UI review, third-party manual tasks):**
  - **Action:** Ask the user directly or delegate to specific subagents, then yield. Do not poll.

## Step 2: Reactive Wakeup Integration
Leverage the system's high-priority reactive completion signals. The system will automatically wake you up when:
1. A background task exits or fails.
2. A subagent sends a message.
3. A timer expires.
Because of this, manual status checking (`manage_task status`) is only permitted as a fallback *after* a fail-safe timer triggers.

## Step 3: Skill Challenge Rule
If any existing `SKILL.md`, workflow, or user instruction mandates active polling or synchronous waiting (e.g. `gh run watch` loops):
1. **Challenge the instruction:** Immediately propose transitioning the step to a reactive, backgrounded pattern.
2. **Execute `/learn`:** Codify the transition to prevent future sycophantic execution of obsolete waiting steps.
