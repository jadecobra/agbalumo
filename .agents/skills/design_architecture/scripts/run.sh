#!/usr/bin/env bash
set -e
FEATURE_NAME="${1:-new-feature}"

echo "[ProductOwner/SystemsArchitect] Initializing design architecture for: $FEATURE_NAME"

# Path to our local skill scripts
SKILL_DIR="$(dirname "$0")"

# 1. Initialize the harness (this helps setup the environment but the skill logic is here)
./scripts/agent-exec.sh workflow init "$FEATURE_NAME" || echo "Note: Initialized harness state"

# 2. Check if the artifact was created
PLAN_FILE="implementation_plan.md"

if [ ! -f "$PLAN_FILE" ]; then
    echo "Wait: Generating initial $PLAN_FILE template..."
    cat <<EOF > "$PLAN_FILE"
# Implementation Plan: $FEATURE_NAME

## Target User Avatar
[Define who this serves - e.g., 'The First-Gen Student']

## Pain Point Mapping
- [Point 1: e.g., High fees]

## Strategic Critique
[Why is this not "dumb"? Push back on bloat!]

## Technical Contract
\`\`\`go
// Use types and interfaces
\`\`\`

## Security STRIDE
- [Boundary Analysis]
EOF
fi

# 3. Call the ChiefCritic Gate
echo "Calling ChiefCritic Audit..."
"$SKILL_DIR/critic-gate.sh" "$PLAN_FILE"

# 4. Verify API Spec if applicable
if [ -f "docs/openapi.yaml" ]; then
    swagger-cli validate docs/openapi.yaml || echo "Warning: Spec validation failed"
    ./scripts/agent-exec.sh verify api-spec
fi
