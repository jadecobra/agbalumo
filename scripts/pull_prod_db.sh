#!/bin/sh
# scripts/pull_prod_db.sh
# Pulls a production SQLite snapshot from Fly.io volume to local .tester/data/prod_snapshot.db
# Used by `agbalumo social draft` to maintain prod-parity UUIDs for social deep links.

set -e

DEST="${1:-.tester/data/prod_snapshot.db}"
mkdir -p "$(dirname "$DEST")"

echo "📥 Pulling production SQLite database from Fly.io (/data/agbalumo.db) to ${DEST}..."
fly ssh sftp get /data/agbalumo.db "${DEST}"
echo "✅ Successfully synced production database snapshot to ${DEST}."
