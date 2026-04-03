#!/bin/sh
set -e

# verify-action-shas.sh - Enforce SHA pinning for GitHub Actions
# 10x Engineering Protocol: Verified Infrastructure

# Find project root for robustness
ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
cd "$ROOT"

EXIT_CODE=0
# Find all workflow and action files
FILES=$(find .github/workflows -name "*.yml"; find .github/actions -name "action.yml")

for FILE in $FILES; do
  # Skip if file not found
  [ -f "$FILE" ] || continue
  
  LINE_NUM=0
  while IFS= read -r LINE; do
    LINE_NUM=$((LINE_NUM + 1))
    
    # Check if the line has 'uses:' and is not a local action or a comment
    if echo "$LINE" | grep -q "uses:" && ! echo "$LINE" | grep -q "uses: \./" && ! echo "$LINE" | grep -q "^[[:space:]]*#"; then
      # Extract the action part (handle cases with or without quotes)
      ACTION_SPEC=$(echo "$LINE" | sed -E "s/.*uses:[[:space:]]+(['\"]?)([^'\"[:space:]#]+).*/\2/")
      
      # Check if it has an '@' followed by a 40-character SHA
      if ! echo "$ACTION_SPEC" | grep -qE "@[0-9a-f]{40}$"; then
        echo "❌ Error in $FILE (Line $LINE_NUM): Action '$ACTION_SPEC' must be pinned to a 40-character SHA."
        EXIT_CODE=1
      fi
      
      # Verify presence of version comment (# vX.Y.Z) for 10x clarity
      if ! echo "$LINE" | grep -q "# v"; then
        echo "⚠️  Warning in $FILE (Line $LINE_NUM): Action '$ACTION_SPEC' is missing a version comment (e.g. # v1.0.0)."
      fi
    fi
  done < "$FILE"
done

if [ $EXIT_CODE -eq 0 ]; then
  echo "✅ All GitHub Actions are correctly pinned to SHAs."
else
  echo "❌ Infrastructure drift detected. Fail."
fi

exit $EXIT_CODE
