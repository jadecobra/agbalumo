#!/bin/bash
# scripts/check-ignored.sh - Check for ignored files staged for commit
STAGED_NEW_OR_MOD=$(git diff --cached --name-only --diff-filter=ACMR || true)
if [ -n "$STAGED_NEW_OR_MOD" ]; then
    IGNORED_STAGED=$(echo "$STAGED_NEW_OR_MOD" | git check-ignore --no-index --stdin || true)
    if [ -n "$IGNORED_STAGED" ]; then
        echo "Error: The following ignored files are staged for commit:"
        echo "$IGNORED_STAGED" | sed 's/^/    /'
        exit 1
    fi
fi
