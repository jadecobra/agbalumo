#!/bin/bash
# refactor.sh: orchestrates the Refactor-with-Teeth loop.

set -e

GOBIN="${GOBIN:-$(pwd)/.tester/tmp/go/bin}"
STATS_FILE=".tester/tmp/refactor_stats.json"

case "$1" in
  init)
    echo "📊 Capturing Refactor Baseline..."
    # Calculate current complexity for all staged files
    FILES=$(git diff --staged --name-only --diff-filter=ACMR | grep '\.go$' || true)
    if [ -z "$FILES" ]; then
      FILES=$(git diff --name-only --diff-filter=ACMR | grep '\.go$' || true)
    fi
    if [ -z "$FILES" ]; then
      echo "No Go files modified to baseline."
      exit 0
    fi
    # Store initial gocognit scores
    "$GOBIN/gocognit" $FILES | awk '{print $4 " " $1}' | sort > "$STATS_FILE"
    echo "✅ Baseline captured in $STATS_FILE"
    ;;
    
  verify)
    echo "⚖️  Verifying Refactor Quality..."
    FILES=$(git diff --staged --name-only --diff-filter=ACMR | grep '\.go$' || true)
    if [ -z "$FILES" ]; then
      FILES=$(git diff --name-only --diff-filter=ACMR | grep '\.go$' || true)
    fi
    if [ -z "$FILES" ]; then
      echo "No Go files to verify."
      exit 0
    fi

    echo "--- Complexity Audit ---"
    # Mandatory threshold
    FAILED=0
    while read -r line; do
        score=$(echo $line | awk '{print $1}')
        file=$(echo $line | awk '{print $4}')
        if [ "$score" -gt 12 ]; then
            echo "❌ CRITICAL: $file has cognitive complexity $score (Threshold: 12)"
            FAILED=1
        fi
    done < <("$GOBIN/gocognit" $FILES)

    # Drift check
    if [ -f "$STATS_FILE" ]; then
        echo "--- Drift Check ---"
        "$GOBIN/gocognit" $FILES | awk '{print $4 " " $1}' | sort > "${STATS_FILE}.new"
        # Compare scores
        while read -r file score; do
            old_score=$(grep "^$file" "$STATS_FILE" | awk '{print $2}' || true)
            if [ -n "$old_score" ]; then
                if [ "$score" -gt "$old_score" ]; then
                    echo "❌ REGRESSION: $file complexity increased from $old_score to $score"
                    FAILED=1
                elif [ "$score" -lt "$old_score" ]; then
                    echo "❇️  IMPROVEMENT: $file complexity decreased from $old_score to $score"
                fi
            else
                echo "ℹ️  NEW FUNCTION: $file introduced with complexity $score"
            fi
        done < "${STATS_FILE}.new"
    fi

    # Duplication check
    echo "--- Duplication Audit ---"
    DUPS=$("$GOBIN/dupl" -threshold 15 $FILES | grep -v "Found total 0 clone groups." || true)
    if [ -n "$DUPS" ]; then
        echo "❌ DUPLICATION FOUND:"
        echo "$DUPS"
        FAILED=1
    else
        echo "✅ No new duplication detected."
    fi

    if [ $FAILED -ne 0 ]; then
        echo "⛔ Refactor rejected. Fix the issues and try again."
        exit 1
    fi
    echo "🏆 10x Refactor Passed!"
    ;;

  *)
    echo "Usage: $0 {init|verify}"
    exit 1
    ;;
esac
