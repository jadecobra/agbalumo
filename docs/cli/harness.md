# agbalumo CLI: Agent Harness & Testing

Internal commands for managing the agent workflow harness and platform testing.

### harness

Agent harness and testing toolkit.

```bash
agbalumo harness [command]
```

### init
Initialize agent workflow.
```bash
agbalumo init [feature-name] [workflow-type]
```

### status
Show current agent workflow status.
```bash
agbalumo status
```

### set-phase
Set current agent phase (RED, GREEN, REFACTOR).
```bash
agbalumo set-phase [phase]
```

### gate
Verify a specific agent workflow gate.
```bash
agbalumo gate [gate-id]
```

### verify
High-level verification of current gate status.
```bash
agbalumo verify [gate-id]
```

### handoff
Generate a HANDOFF.md bridge between personas.
```bash
agbalumo handoff [target-persona]
```

### chaos
Inject failures into the harness for resilience testing.
```bash
agbalumo chaos [flags]
```

### cost
Audit the codebase to measure the agent context cost (RMS of LOC).
```bash
agbalumo cost [dir]
```

### update-coverage
Update coverage threshold for specific packages.
```bash
agbalumo update-coverage <package_path> <threshold>
```

### stress
Generate stress test data (listings/users).
```bash
agbalumo stress [flags]
```

### benchmark
Run performance query benchmarks.
```bash
agbalumo benchmark [flags]
```

### hello
Protocol test command for agent handoff verification.
```bash
agbalumo hello
```
