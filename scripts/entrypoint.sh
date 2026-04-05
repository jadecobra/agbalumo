#!/bin/bash
set -e

# --- Volume Permissions ---
# Fly.io volumes are often mounted as root. We ensure the appuser owns the data directory.
echo "[entrypoint] Ensuring /data and /app permissions..."
chown -R appuser:appuser /data /app

# --- Run Tasks as appuser ---
# We use su-exec to drop privileges and start the app process.
# We run the city backfill in the background as appuser.

echo "[entrypoint] Starting city backfill in background..."
su-exec appuser /bin/bash -c "/app/server listing backfill-cities > /tmp/backfill.log 2>&1 &"

# Restore from replica if database does not exist
if [ ! -f /data/agbalumo.db ]; then
    echo "[entrypoint] No database found at /data/agbalumo.db. Attempting restore from R2..."
    su-exec appuser /usr/local/bin/litestream restore -config /etc/litestream.yml /data/agbalumo.db && \
        echo "[entrypoint] Restore successful." || \
        echo "[entrypoint] No remote data found (first run or empty bucket). Starting fresh."
else
    echo "[entrypoint] Existing database found at /data/agbalumo.db. Skipping restore."
fi

# Start the app under Litestream replication
echo "[entrypoint] Starting Litestream replication and app server..."
exec su-exec appuser /usr/local/bin/litestream replicate -config /etc/litestream.yml -exec "/app/server serve"
