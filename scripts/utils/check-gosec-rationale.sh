#!/bin/bash
# check-gosec-rationale.sh: Verify that all // #nosec directives include a rationale comment.
# Rationale is expected to be preceded by a hyphen (-) or double-hyphen (--).

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default to current directory if no path provided
SEARCH_PATH="${1:-.}"

echo "Checking for mandatory rationale in // #nosec directives..."

# Find all // #nosec directives that don't contain a hyphen or double-hyphen for rationale
# We use grep -r and look for patterns like // #nosec G101 without a following - or --
# A valid one looks like: // #nosec G101 - rationale or // #nosec G101 -- rationale
# An invalid one looks like: // #nosec G101 or // #nosec G101,G102

# Regular expression explained:
# ^.*//\s*#nosec\s+G[0-9]+(\s+G[0-9]+)*\s*$  -> Matches bare #nosec with rule IDs but nothing after
# We also want to catch ones that have something after but no hyphen
# So we search for all #nosec and then filter out those that have a hyphen

INVALID_DIRECTIVES=$(grep -rE "//\s*#nosec" "$SEARCH_PATH" --include="*.go" --exclude-dir=".tester" --exclude-dir=".go" --exclude-dir="tmp" --exclude-dir="vendor" | grep -vE " - | -- " || true)

if [ -n "$INVALID_DIRECTIVES" ]; then
    echo -e "${RED}❌ Error: Found // #nosec directives without a mandatory rationale comment.${NC}"
    echo -e "${YELLOW}Rationale must be preceded by a hyphen (-) or double-hyphen (--).${NC}"
    echo -e "${YELLOW}Example: // #nosec G304 - Internal file reading${NC}"
    echo ""
    echo "$INVALID_DIRECTIVES"
    exit 1
fi

echo -e "${GREEN}✅ All // #nosec directives have rationales.${NC}"
exit 0
