#!/bin/sh
# scripts/pushw.sh
# Git Push & Watch wrapper
# Automatically monitors the remote CI run for the pushed commit.

export GH_PROMPT_DISABLED=1

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" = "main" ] && [ "$ALLOW_MAIN_PUSH" != "1" ]; then
    echo "❌ Error: Direct push to 'main' is prohibited. Work on a feature branch and open a PR."
    exit 1
fi

git push "$@"
if [ $? -eq 0 ]; then
    if [ "$CURRENT_BRANCH" != "main" ]; then
        PR_NUM=$(gh pr view "$CURRENT_BRANCH" --json number,state -q 'select(.state == "OPEN") | .number' 2>/dev/null || true)
        if [ -z "$PR_NUM" ]; then
            echo "ℹ️ No open PR found for branch '$CURRENT_BRANCH'. Remote CI triggers on pull_request."
            echo "   Run ./scripts/pr-create.sh to open a PR and start CI."
            exit 0
        fi
    fi

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
        echo "⚠️ Could not find CI run for commit ${COMMIT_SHA}."
        exit 1
    else
        echo "🔍 Found CI run ${RUN_ID}. Monitoring progress..."
        gh run watch "$RUN_ID" --exit-status
    fi
fi
