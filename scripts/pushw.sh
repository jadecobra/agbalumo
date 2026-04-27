#!/bin/sh
# scripts/pushw.sh
# Git Push & Watch wrapper
# Automatically monitors the remote CI run for the pushed commit.

git push "$@"
if [ $? -eq 0 ]; then
    COMMIT_SHA=$(git rev-parse HEAD)
    echo "✅ Push successful! Waiting for CI run to register for commit ${COMMIT_SHA}..."
    
    RUN_ID=""
    for i in $(seq 1 30); do
        RUN_ID=$(gh run list --commit "$COMMIT_SHA" --limit 1 --json databaseId --jq '.[0].databaseId' 2>/dev/null)
        if [ -n "$RUN_ID" ] && [ "$RUN_ID" != "null" ]; then
            break
        fi
        sleep 2
    done

    if [ -z "$RUN_ID" ] || [ "$RUN_ID" = "null" ]; then
        echo "⚠️ Could not find CI run for commit ${COMMIT_SHA}. Falling back to default watch..."
        gh run watch --exit-status
    else
        echo "🔍 Found CI run ${RUN_ID}. Monitoring progress..."
        gh run watch "$RUN_ID" --exit-status
    fi
fi
