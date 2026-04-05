#!/bin/bash
set -e

# Check for format flag
FMT="json"
if [ "$1" = "--text" ]; then 
    FMT="text"
    shift
fi

source "$(dirname "$0")/utils.sh"

# Default to environment variables if set, otherwise prompt
CLIENT_ID=${GOOGLE_CLIENT_ID}
CLIENT_SECRET=${GOOGLE_CLIENT_SECRET}
ADMIN_CODE=${ADMIN_CODE}
SESSION_SECRET=${SESSION_SECRET}

if [ -z "$CLIENT_ID" ] && [ "$FMT" = "text" ]; then
    echo "Enter your Google Client ID:"
    read -r CLIENT_ID
fi

if [ -z "$CLIENT_SECRET" ] && [ "$FMT" = "text" ]; then
    echo "Enter your Google Client Secret:"
    read -r CLIENT_SECRET
fi

if [ -z "$ADMIN_CODE" ] && [ "$FMT" = "text" ]; then
    echo "Enter your Admin Access Code:"
    read -r ADMIN_CODE
fi

if [ -z "$SESSION_SECRET" ] && [ "$FMT" = "text" ]; then
    echo "Enter your Session Secret (for cookie signing):"
    read -r SESSION_SECRET
fi

if [ -z "$CLIENT_ID" ] || [ -z "$CLIENT_SECRET" ] || [ -z "$ADMIN_CODE" ] || [ -z "$SESSION_SECRET" ]; then
    if [ "$FMT" = "text" ]; then
        echo "Error: All secrets (Google Client ID/Secret, Admin Code, Session Secret) are required."
    else
        output_json_envelope false "deploy_secrets.sh" "All secrets are required if not provided via environment." "[]"
    fi
    exit 1
fi

if [ "$FMT" = "text" ]; then echo "Setting secrets in Fly.io..."; fi
# Capture output to return in JSON
OUT=$(fly secrets set \
    GOOGLE_CLIENT_ID="$CLIENT_ID" \
    GOOGLE_CLIENT_SECRET="$CLIENT_SECRET" \
    ADMIN_CODE="$ADMIN_CODE" \
    SESSION_SECRET="$SESSION_SECRET" \
    2>&1) || {
    if [ "$FMT" = "text" ]; then
        echo "Error: Failed to set secrets in Fly.io"
        echo "$OUT"
    else
        output_json_envelope false "deploy_secrets.sh" "Failed to set secrets in Fly.io: $OUT" "[]"
    fi
    exit 1
}

if [ "$FMT" = "text" ]; then
    echo "$OUT"
    echo "Secrets set successfully! The application will restart automatically."
else
    output_json_envelope true "deploy_secrets.sh" "Secrets set successfully! The application will restart automatically." "[]"
fi
