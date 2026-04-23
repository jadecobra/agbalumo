# Handler Constraints

- Use `RespondError(c, err)` — never raw `c.JSON()`
- All form bindings use `form` struct tags
- No raw HTML in handlers — use `ui/templates/components/`
