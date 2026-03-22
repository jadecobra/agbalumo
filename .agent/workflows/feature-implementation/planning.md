# Planning & Pre-Flight Rules

## 0. Initialize Workflow State

Before any code is written, initialize the state machine to track gate progress.

// turbo
1. Run `./scripts/agent-exec.sh workflow init <feature-name>`
// turbo
2. Run `./scripts/agent-exec.sh workflow set-phase RED`
