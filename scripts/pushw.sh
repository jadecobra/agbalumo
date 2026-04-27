#!/bin/sh
# scripts/pushw.sh
# Git Push & Watch wrapper
# Automatically monitors the remote CI run immediately after a successful push.

git push "$@"
if [ $? -eq 0 ]; then
    echo "✅ Push successful! Waiting for CI to register..."
    sleep 5
    gh run watch --exit-status
fi
