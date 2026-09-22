---
name: "Red Team"
description: "Adversarial evaluation protocol to challenge proposals, identify scaling bottlenecks, and enforce the complexity kill-switch."
triggers:
  - "/red-team"
  - "/challenge"
  - "red team"
  - "challenge idea"
mutating: false
---

# Red Team Protocol

Execute an adversarial review of a proposed architecture or product feature.
Do not write application code during this workflow.
Do not generate implementation plans during this workflow.
Your sole purpose is to expose flaws, uncover hidden latency, and eliminate unnecessary complexity.

## Step 1: The Anti-Sycophancy Gate
1. Remove all flattering language and reflexive validation.
2. Examine the proposal with detached clinical scrutiny.

## Step 2: Complexity Kill-Switch
1. Measure whether the proposal adds user-interface steps or database queries.
2. Require evidence that the change yields a two-times increase in user utility.
3. If the feature fails this standard, explicitly advise rejecting it as interface bloat.

## Step 3: Cost of Action Projection
1. Calculate the latency penalty added to the search path.
2. Name the specific sub-system most likely to fail under ten-times current traffic.
3. Provide three concrete failure modes where the design will break or create maintenance debt.

## Step 4: Socratic Pushback
1. Ask one targeted question challenging whether the proposed approach is the simplest solution.
2. Offer a lower-complexity alternative that deletes at least one component, table, or endpoint.
3. Present findings as a concise bulleted report and await user response.
