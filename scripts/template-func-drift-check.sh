#!/bin/bash
set -e

# Colors for output
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
BLUE=$(printf '\033[0;34m')
YELLOW=$(printf '\033[1;33m')
NC=$(printf '\033[0m')

printf "${BLUE}Running Template Function Drift Check...${NC}\n"

# 1. Extract function names from ui/renderer.go
# We look for keys in the template.FuncMap
RENDERER_FILE="internal/ui/renderer.go"
if [ ! -f "$RENDERER_FILE" ]; then
    printf "${RED}❌ Renderer file mot found: %s${NC}\n" "$RENDERER_FILE"
    exit 1
fi

DEFINED_FUNCS=$(grep -E '^		"[a-zA-Z0-9]+":' "$RENDERER_FILE" | sed -E 's/.*"([a-zA-Z0-9]+)".*/\1/' | sort -u)

if [ -z "$DEFINED_FUNCS" ]; then
    printf "${RED}❌ Could not extract function names from %s${NC}\n" "$RENDERER_FILE"
    exit 1
fi

# 2. Extract function calls from all .html templates
# Pattern: {{ [a-zA-Z0-9_]+ ... }} or {{ range [a-zA-Z0-9_]+ ... }}
# We'll use a more comprehensive regex to find function calls in templates.
# Standard calls look like: {{ func ... }} or pipeline: {{ .Var | func }}
USED_FUNCS_FILE=$(mktemp)

# Grep for {{ func or {{ range func or | func
find ui/templates -name "*.html" -exec cat {} + | \
    grep -oE '\{\{[[:space:]]*(range[[:space:]]+)?([a-zA-Z0-9]+)[[:space:]]' | \
    sed -E 's/\{\{[[:space:]]*(range[[:space:]]+)?([a-zA-Z0-9]+).*/\2/' >> "$USED_FUNCS_FILE"

find ui/templates -name "*.html" -exec cat {} + | \
    grep -oE '\|[[:space:]]*([a-zA-Z0-9]+)' | \
    sed -E 's/\|[[:space:]]*([a-zA-Z0-9]+).*/\1/' >> "$USED_FUNCS_FILE"

# Filter out built-in/reserved keywords and common template variables
FILTERED_USED_FUNCS=$(sort -u "$USED_FUNCS_FILE" | grep -vE '^(if|else|end|range|with|define|block|template|nil|len|and|or|not|index|slice|printf|print|println|html|urlquery|js|call)$' | grep -vE '^\.')

rm "$USED_FUNCS_FILE"

DRIFT_DETECTED=0

# 3. Check if all USED functions are DEFINED
for func in $FILTERED_USED_FUNCS; do
    if ! echo "$DEFINED_FUNCS" | grep -qxw "$func"; then
        # Check if it might be a field access starting with dot (though filtered, let's be safe)
        if [[ $func =~ ^[A-Z] ]]; then
            # Likely a field direct access if Uppercase, but templates use .Field
            # If it's used as {{ Field }}, it's probably a pipeline or we missed the dot.
            # But Go templates usually require the dot for fields.
            # However, some funcs might be provided but not in our global map (though they should be).
            continue
        fi
        
        printf "${RED}❌ Undefined template function used: '%s'${NC}\n" "$func"
        DRIFT_DETECTED=1
    fi
done

if [ $DRIFT_DETECTED -eq 1 ]; then
    printf "\n${RED}❌ Template Function Drift Detected! Please add missing functions to %s${NC}\n" "$RENDERER_FILE"
    exit 1
else
    printf "\n${GREEN}✅ All template functions are in sync.${NC}\n"
    exit 0
fi
