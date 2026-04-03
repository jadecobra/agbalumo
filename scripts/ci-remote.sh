#!/bin/bash
# scripts/ci-remote.sh: Checks production CI status using GitHub CLI.
# Fast and reliable status check for production CI workflows.

set -e

BRANCH=$(git branch --show-current)
echo "🔍 Checking production CI status for branch: $BRANCH"

# Use gh cli to get the status of the latest run on the current branch
# Filter by the 'CI' workflow if applicable, or just get the most recent.
LAST_RUN=$(gh run list --branch "$BRANCH" --limit 1 --json databaseId,status,conclusion,url)

if [ -z "$LAST_RUN" ] || [ "$LAST_RUN" == "[]" ]; then
  echo "⚠️  No CI runs found for branch $BRANCH."
  exit 1
fi

ID=$(echo "$LAST_RUN" | jq -r '.[0].databaseId')
STATUS=$(echo "$LAST_RUN" | jq -r '.[0].status')
CONCLUSION=$(echo "$LAST_RUN" | jq -r '.[0].conclusion')
URL=$(echo "$LAST_RUN" | jq -r '.[0].url')

echo "Run ID: $ID"
echo "Status: $STATUS"
echo "Conclusion: $CONCLUSION"
echo "URL: $URL"

if [ "$STATUS" == "completed" ] && [ "$CONCLUSION" == "success" ]; then
  echo "✅ Production CI passed for this branch."
  exit 0
elif [ "$STATUS" == "completed" ]; then
  echo "❌ Production CI failed (conclusion: $CONCLUSION)."
  exit 1
else
  echo "⏳ Production CI is currently $STATUS..."
  exit 1
fi
