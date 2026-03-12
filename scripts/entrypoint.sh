#!/bin/sh
set -e

# Ensure the /data directory exists and is owned by the app user.
mkdir -p /data

# --- Restore from replica if database does not exist ---
# This is the disaster recovery path. On a fresh machine with an empty volume,
# Litestream restores the latest snapshot from Cloudflare R2 before the app starts.
if [ ! -f /data/agbalumo.db ]; then
    echo "[entrypoint] No database found at /data/agbalumo.db. Attempting restore from R2..."
    litestream restore -config /etc/litestream.yml /data/agbalumo.db && \
        echo "[entrypoint] Restore successful." || \
        echo "[entrypoint] No remote data found (first run or empty bucket). Starting fresh."
else
    echo "[entrypoint] Existing database found at /data/agbalumo.db. Skipping restore."
fi

# --- Run City Backfill ---
# This ensures that any listings with missing cities (from bulk uploads or legacy data)
# are geocoded before the server starts.
echo "[entrypoint] Running city backfill..."
./server listing backfill-cities || echo "[entrypoint] Backfill failed or skipped."

# --- Start the app under Litestream replication ---
# Litestream acts as a supervisor: it starts the app process and continuously
# replicates WAL changes to Cloudflare R2. If the app crashes, Litestream exits too.
echo "[entrypoint] Starting Litestream replication and app server..."
exec litestream replicate -config /etc/litestream.yml -exec "./server serve"
