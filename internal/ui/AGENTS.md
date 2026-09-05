# UI Package (internal/ui)

HTML template rendering engine and error response handlers for Echo.

## Constraints
- `TemplateRenderer` implements `echo.Renderer` using `html/template` with helper functions from `renderer_funcs.go`.
- `RespondError(c, err)` and `RespondErrorMsg(c, code, msg)` in `error.go` are canonical error responders for handlers.
- All template data passed to `c.Render` must use strictly typed ViewModel structs (enforced by `verify template-contract`).
- HTML templates live in `ui/templates/`. See `ui/templates/AGENTS.md` for Global Editorial Brutalist rules.
