# Context Cost & Efficiency (Tokens)

To maintain a 10x engineering pace, we monitor the codebase context's efficiency for agentic assistants. We use token-based metrics (Tiktoken/BPE) to understand the cognitive load on the models.

### Key Metrics (Advisory Health Indicators)
- **TotalTokens**: Total count of segments (excluding `vendor/`, `.git/`, and `ui/static/`).
- **TokenRMS**: Root-mean-square density of tokens across files. **Warning Threshold: > 110**.
- **ContextWindowPct**: Percentage of a standard 200k context window consumed. **Warning Threshold: > 50%**.

### 🚫 The "No Gaming" Rule
Context cost metrics are **Advisory Health Indicators**, not hard gates.
- **NEVER** reduce context cost by removing useful documentation, comments, safety checks, or unit tests.
- **NEVER** sacrifice code readability or architectural clarity to lower token counts.
- **ONLY** refactor for context when a file truly becomes "un-agentable" (e.g., > 1000 tokens) and splitting would improve modularity and focus.

### Monitoring
Run the harness cost tool for visibility:
```bash
./scripts/agent-exec.sh cost
```

### Efficiency Patterns
1. **Modularization**: High-density files (> 500 tokens) should be reviewed by the **SystemsArchitect** for potential splitting.
2. **Asset Exclusion**: Large static assets must remain ignored in `internal/agent/cost.go` to prevent skewing the cognitive load metrics.
