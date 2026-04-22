#!/bin/sh
# scripts/pushw.sh
# Git Push & Watch wrapper
# Automatically launches the watch tool immediately after a successful push.

git push "$@"
if [ $? -eq 0 ]; then
    echo "✅ Push successful! Starting watch..."
    go run ./cmd/verify watch
fi
