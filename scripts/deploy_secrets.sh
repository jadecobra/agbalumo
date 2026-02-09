#!/bin/bash
set -e

# Default to environment variables if set, otherwise prompt
CLIENT_ID=${GOOGLE_CLIENT_ID}
CLIENT_SECRET=${GOOGLE_CLIENT_SECRET}

if [ -z "$CLIENT_ID" ]; then
    echo "Enter your Google Client ID:"
    read -r CLIENT_ID
fi

if [ -z "$CLIENT_SECRET" ]; then
    echo "Enter your Google Client Secret:"
    read -r CLIENT_SECRET
fi

if [ -z "$CLIENT_ID" ] || [ -z "$CLIENT_SECRET" ]; then
    echo "Error: Both Client ID and Client Secret are required."
    exit 1
fi

echo "Setting secrets in Fly.io..."
fly secrets set GOOGLE_CLIENT_ID="$CLIENT_ID" GOOGLE_CLIENT_SECRET="$CLIENT_SECRET"

echo "Secrets set successfully! The application will restart automatically."
