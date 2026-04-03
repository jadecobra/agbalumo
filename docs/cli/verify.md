# agbalumo CLI: Verification & Maintenance

The `verify` command provides subcommands for ensuring codebase health, documentation sync, and coverage standards.

## Commands

### verify

Agbalumo Maintenance and Verification Utility.

```bash
agbalumo verify [command]
```

#### Subcommands

##### audit

Run comprehensive security and health audit.

```bash
agbalumo verify audit
```

##### api-spec

Detect drift between Code, OpenAPI, and Markdown docs.

```bash
agbalumo verify api-spec
```

##### template-drift

Detect undefined template functions in HTML templates.

```bash
agbalumo verify template-drift
```

##### context-cost

Calculate codebase token density and context window usage (advisory).

```bash
agbalumo verify context-cost
```

##### coverage

Enforce coverage threshold anti-degradation.

```bash
agbalumo verify coverage
```
