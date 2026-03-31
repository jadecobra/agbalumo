#!/bin/zsh

# Check if aglog binary exists, if not build it
BINARY="bin/aglog"
if [[ ! -f "$BINARY" ]]; then
    # Source sandbox if it exists and we're in a restricted environment
    if [[ -f "scripts/sandbox.env" ]]; then
        source scripts/sandbox.env
    fi
    echo "Building aglog..."
    go build -o "$BINARY" cmd/aglog/main.go
fi

# Initialize arguments
FEATURE=""
ARCH=""
PO=""
SDET=""
BE=""
SUMMARY=""

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --feature) FEATURE="$2"; shift ;;
        --arch) ARCH="$2"; shift ;;
        --po) PO="$2"; shift ;;
        --sdet) SDET="$2"; shift ;;
        --be) BE="$2"; shift ;;
        --summary) SUMMARY="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

if [[ -z "$FEATURE" ]]; then
    echo "Error: --feature is required"
    exit 1
fi

# Run aglog and capture output (filepath)
FILEPATH=$($BINARY --feature "$FEATURE" --arch "$ARCH" --po "$PO" --sdet "$SDET" --be "$BE" --summary "$SUMMARY")

if [[ $? -ne 0 ]]; then
    echo "Error running aglog"
    echo "$FILEPATH"
    exit 1
fi

# Output the file path
echo "$FILEPATH"

# Prepare JSON data for the sync marker
# We use the same keys as SquadDecision struct
JSON_DATA=$(cat <<EOF
{
  "FeatureName": "$FEATURE",
  "SystemsArchitect": "$ARCH",
  "ProductOwner": "$PO",
  "SDET": "$SDET",
  "BackendEngineer": "$BE",
  "DecisionSummary": "$SUMMARY"
}
EOF
)

# Output marker for agent synchronization
echo "[SQUAD-MEMORY-SYNC] $JSON_DATA"
