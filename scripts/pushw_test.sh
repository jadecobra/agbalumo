#!/bin/sh
# scripts/pushw_test.sh
# Reproduction and regression test suite for scripts/pushw.sh
# Verifies non-interactive batch mode and immediate exit when no open PR exists.

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PUSHW_SCRIPT="$SCRIPT_DIR/pushw.sh"

TMP_DIR=$(mktemp -d /tmp/pushw_test_XXXXXX)
trap 'rm -rf "$TMP_DIR"' EXIT INT TERM

MOCK_BIN="$TMP_DIR/bin"
mkdir -p "$MOCK_BIN"

# Create mock git
cat << 'EOF' > "$MOCK_BIN/git"
#!/bin/sh
case "$1" in
    rev-parse)
        if [ "$2" = "--abbrev-ref" ] && [ "$3" = "HEAD" ]; then
            echo "${MOCK_BRANCH:-feat/test-branch}"
            exit 0
        elif [ "$2" = "HEAD" ]; then
            echo "mockcommitsha123456"
            exit 0
        fi
        ;;
    push)
        exit 0
        ;;
esac
echo "Unknown git invocation: $@" >&2
exit 1
EOF
chmod +x "$MOCK_BIN/git"

# Create mock gh
cat << 'EOF' > "$MOCK_BIN/gh"
#!/bin/sh
if [ "$1" = "run" ] && [ "$2" = "list" ]; then
    if [ "$MOCK_HAS_RUN" = "1" ]; then
        echo "999888"
        exit 0
    fi
    echo ""
    exit 0
fi

if [ "$1" = "pr" ] && [ "$2" = "view" ]; then
    case "$MOCK_PR_STATE" in
        "NONE")
            echo "no pull requests found for branch" >&2
            exit 1
            ;;
        "MERGED")
            # If queried with jq expression selecting only OPEN state
            if echo "$*" | grep -q 'select(.state == "OPEN")'; then
                echo ""
                exit 0
            fi
            # If queried simply for .number without state filtering
            if echo "$*" | grep -q '.number'; then
                echo "42"
                exit 0
            fi
            echo '{"number": 42, "state": "MERGED"}'
            exit 0
            ;;
        "OPEN")
            if echo "$*" | grep -q '.number'; then
                echo "101"
                exit 0
            fi
            echo '{"number": 101, "state": "OPEN"}'
            exit 0
            ;;
    esac
fi

if [ "$1" = "run" ] && [ "$2" = "watch" ]; then
    if [ -z "$3" ]; then
        if [ -n "$GH_PROMPT_DISABLED" ]; then
            echo "run ID required when not running interactively" >&2
            exit 1
        else
            echo "? Select a workflow run" >&2
            # Simulate interactive hanging by sleeping if not interrupted
            exit 2
        fi
    fi
    echo "Watched run $3"
    exit 0
fi

echo "Unknown gh invocation: $@" >&2
exit 1
EOF
chmod +x "$MOCK_BIN/gh"

export PATH="$MOCK_BIN:$PATH"

FAILED=0

# Test 1: GH_PROMPT_DISABLED environment variable must be exported
echo "--- Test 1: Verify GH_PROMPT_DISABLED export ---"
if grep -q "export GH_PROMPT_DISABLED=1" "$PUSHW_SCRIPT" || grep -q "GH_PROMPT_DISABLED=1" "$PUSHW_SCRIPT"; then
    echo "✅ PASS: GH_PROMPT_DISABLED is configured in pushw.sh"
else
    echo "❌ FAIL: GH_PROMPT_DISABLED is not set in pushw.sh"
    FAILED=1
fi

# Test 2: When no open PR exists, must exit cleanly without 60s delay
echo "--- Test 2: Immediate exit when no PR exists ---"
export MOCK_BRANCH="feat/feature-no-pr"
export MOCK_PR_STATE="NONE"
export MOCK_HAS_RUN="0"

START_TIME=$(date +%s)
OUTPUT=$(sh "$PUSHW_SCRIPT" 2>&1) || EXIT_CODE=$?
EXIT_CODE=${EXIT_CODE:-0}
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

if [ "$EXIT_CODE" -ne 0 ]; then
    echo "❌ FAIL: Expected exit code 0 when no PR exists, got $EXIT_CODE"
    echo "$OUTPUT"
    FAILED=1
elif [ "$DURATION" -ge 5 ]; then
    echo "❌ FAIL: Script took $DURATION seconds (expected < 5s); entered polling loop before PR check"
    FAILED=1
elif ! echo "$OUTPUT" | grep -q "No open PR found"; then
    echo "❌ FAIL: Expected 'No open PR found' in output, got:"
    echo "$OUTPUT"
    FAILED=1
else
    echo "✅ PASS: Exited cleanly in ${DURATION}s with no open PR"
fi

# Test 3: When branch PR is MERGED/CLOSED, must treat as no open PR
echo "--- Test 3: Clean exit when PR state is MERGED ---"
export MOCK_BRANCH="feat/feature-merged-pr"
export MOCK_PR_STATE="MERGED"
export MOCK_HAS_RUN="0"

OUTPUT=$(sh "$PUSHW_SCRIPT" 2>&1) || EXIT_CODE=$?
EXIT_CODE=${EXIT_CODE:-0}

if [ "$EXIT_CODE" -ne 0 ]; then
    echo "❌ FAIL: Expected exit code 0 when PR is MERGED, got $EXIT_CODE (fell back to error exit 1)"
    echo "$OUTPUT"
    FAILED=1
elif ! echo "$OUTPUT" | grep -q "No open PR found"; then
    echo "❌ FAIL: Expected 'No open PR found' for MERGED PR, got:"
    echo "$OUTPUT"
    FAILED=1
else
    echo "✅ PASS: Correctly identified MERGED PR as not open"
fi

if [ "$FAILED" -ne 0 ]; then
    echo "❌ Reproduction test failed as expected on unfixed pushw.sh"
    exit 1
fi

echo "🎉 All pushw.sh tests passed!"
exit 0
