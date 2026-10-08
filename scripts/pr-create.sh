#!/bin/sh
# scripts/pr-create.sh
# Deterministic PR creation with sync verification and commit assertion.

set -e

usage() {
    echo "Usage: $0 --title <title> (--body-file <path> | --body <text>) [--base <base-branch>]"
    echo "   or: $0 <title> (<path-to-file> | --body <text>) [<base-branch>]"
    exit "${1:-1}"
}

TITLE=""
BODY_FILE=""
BODY_TEXT=""
BASE_BRANCH="main"
CLEANUP_FILE=""

trap 'if [ -n "$CLEANUP_FILE" ] && [ -f "$CLEANUP_FILE" ]; then rm -f "$CLEANUP_FILE"; fi' EXIT INT TERM

while [ $# -gt 0 ]; do
    case "$1" in
        --title)
            TITLE="$2"
            shift 2
            ;;
        --body-file)
            BODY_FILE="$2"
            shift 2
            ;;
        --body)
            BODY_TEXT="$2"
            shift 2
            ;;
        --base)
            BASE_BRANCH="$2"
            shift 2
            ;;
        -h|--help)
            usage 0
            ;;
        *)
            if [ -z "$TITLE" ]; then
                TITLE="$1"
                shift
            elif [ -z "$BODY_FILE" ] && [ -z "$BODY_TEXT" ]; then
                BODY_FILE="$1"
                shift
            elif [ "$BASE_BRANCH" = "main" ]; then
                BASE_BRANCH="$1"
                shift
            else
                usage
            fi
            ;;
    esac
done

if [ -n "$BODY_TEXT" ] && [ -z "$BODY_FILE" ]; then
    mkdir -p tmp
    BODY_FILE="tmp/pr_body_$$.md"
    printf "%s\n" "$BODY_TEXT" > "$BODY_FILE"
    CLEANUP_FILE="$BODY_FILE"
fi

if [ -z "$TITLE" ] || [ -z "$BODY_FILE" ]; then
    echo "❌ Error: Title and either --body or --body-file are required."
    usage
fi

if [ ! -r "$BODY_FILE" ]; then
    echo "❌ Error: Body file not found or not readable at: $BODY_FILE"
    exit 1
fi

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" = "HEAD" ] || [ -z "$CURRENT_BRANCH" ]; then
    echo "❌ Error: Detached HEAD state. Please checkout a feature branch."
    exit 1
fi

if [ "$CURRENT_BRANCH" = "$BASE_BRANCH" ]; then
    echo "❌ Error: Cannot create PR from '$BASE_BRANCH' into '$BASE_BRANCH'."
    exit 1
fi

echo "🔍 Verifying remote synchronization for branch '$CURRENT_BRANCH'..."
LOCAL_SHA=$(git rev-parse HEAD)

# Verify upstream tracking branch is configured
UPSTREAM_REF=$(git rev-parse --abbrev-ref --symbolic-full-name '@{u}' 2>/dev/null || true)
if [ -z "$UPSTREAM_REF" ]; then
    echo "❌ Error: Upstream tracking branch is not configured for '$CURRENT_BRANCH'."
    echo "   Run ./scripts/pushw.sh -u origin $CURRENT_BRANCH first to push and synchronize."
    exit 1
fi

UPSTREAM_SHA=$(git rev-parse '@{u}' 2>/dev/null || true)
if [ -z "$UPSTREAM_SHA" ]; then
    echo "❌ Error: Could not resolve SHA for upstream '$UPSTREAM_REF'."
    echo "   Run ./scripts/pushw.sh to synchronize."
    exit 1
fi

if [ "$LOCAL_SHA" != "$UPSTREAM_SHA" ]; then
    echo "❌ Error: Remote branch is out of sync with local HEAD."
    echo "   Local SHA:    $LOCAL_SHA"
    echo "   Upstream SHA: $UPSTREAM_SHA"
    echo "   Run ./scripts/pushw.sh to synchronize before creating PR."
    exit 1
fi

echo "✅ Remote branch is synchronized ($LOCAL_SHA)."
echo "🚀 Creating pull request into '$BASE_BRANCH'..."

PR_URL=$(gh pr create --base "$BASE_BRANCH" --head "$CURRENT_BRANCH" --title "$TITLE" --body-file "$BODY_FILE")

echo "✅ Pull request created: $PR_URL"
echo "🔍 Validating PR commits and metadata..."

PR_DATA=$(gh pr view "$CURRENT_BRANCH" --json number,url,commits)
PR_NUMBER=$(echo "$PR_DATA" | jq -r '.number')
COMMITS_COUNT=$(echo "$PR_DATA" | jq -r '.commits | length')

if [ "$COMMITS_COUNT" -eq 0 ]; then
    echo "❌ Error: PR #$PR_NUMBER was created but contains 0 commits between '$CURRENT_BRANCH' and '$BASE_BRANCH'."
    exit 1
fi

echo "🎉 PR #$PR_NUMBER verified with $COMMITS_COUNT commit(s): $PR_URL"
