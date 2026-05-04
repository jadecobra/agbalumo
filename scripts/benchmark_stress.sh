#!/bin/bash

# Load Invariants
INVARIANTS_FILE=".agents/invariants.json"
if [ -f "$INVARIANTS_FILE" ]; then
    PROTOCOL=$(jq -r '.protocol' "$INVARIANTS_FILE")
    INV_PORT=$(jq -r '.port' "$INVARIANTS_FILE")
else
    PROTOCOL="http"
    INV_PORT="8080"
fi

# Accept count as argument, default to 100000
COUNT=${1:-100000}
CONCURRENCY=${2:-50}
REQUESTS=${3:-500}
PORT=${4:-$INV_PORT} # Default to invariant port

DB_PATH=".tester/data/benchmark.db"
LOG_FILE=".tester/benchmark_server.log"

# Cleanup any previous benchmark DB and logs
rm -f "$DB_PATH"
rm -f "${DB_PATH}-shm"
rm -f "${DB_PATH}-wal"
rm -f "$LOG_FILE"

echo "=== System Invariants ==="
echo "Protocol: $PROTOCOL"
echo "Target Port: $PORT"

echo "=== Compiling CLI ==="
mkdir -p .tester/db_tests
go build -o .tester/db_tests/tmp_harness main.go || exit 1

echo "=== Write Benchmark (${COUNT} Inserts) ==="
./.tester/db_tests/tmp_harness stress -c $COUNT "$DB_PATH" || exit 1

echo "=== Reconnect & Query DB ==="
echo "Total Listings Generated:"
sqlite3 "$DB_PATH" "SELECT count(*) FROM listings;"

echo "=== Basic Read Benchmark ==="
echo "Querying Page 1 (No Filters) - Top 20:"
time sqlite3 "$DB_PATH" "SELECT id FROM listings ORDER BY created_at DESC LIMIT 20 OFFSET 0;" > /dev/null

echo "=== HTTP Benchmark Setup ==="
if ! command -v ab &> /dev/null; then
    echo "Error: ab (Apache Benchmark) is not installed."
    exit 1
fi

# Check if port is in use
if lsof -i :$PORT > /dev/null; then
    echo "Warning: Port $PORT is already in use. Stress test may fail or hit the wrong instance."
    # We'll try to proceed but if it fails, it's expected.
fi

echo "Starting server on port $PORT (Protocol: $PROTOCOL)..."
# Use development env to allow local cert loading if HTTPS
DATABASE_URL="$DB_PATH" AGBALUMO_ENV=development PORT=$PORT ./.tester/db_tests/tmp_harness serve > "$LOG_FILE" 2>&1 &
SERVER_PID=$!

# Wait for server to be ready
echo "Waiting for server to respond on $PROTOCOL://127.0.0.1:$PORT/ ..."
RETRY_COUNT=0
MAX_RETRIES=20
# Use -k for curl to ignore self-signed cert issues if protocol is https
CURL_OPTS="-s"
if [ "$PROTOCOL" == "https" ]; then
    CURL_OPTS="-sk"
fi

while ! curl $CURL_OPTS $PROTOCOL://127.0.0.1:$PORT/ > /dev/null; do
    RETRY_COUNT=$((RETRY_COUNT+1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        echo "Server failed to start after $MAX_RETRIES seconds. Check $LOG_FILE"
        kill $SERVER_PID
        exit 1
    fi
    echo "Attempt $RETRY_COUNT/$MAX_RETRIES..."
    sleep 1
done

echo "Server is UP!"

# Setup ab options for HTTPS if needed
AB_OPTS=""
# ab on macOS sometimes needs specific flags for local certs or -f for protocols
# but based on my manual test, it worked without flags if the URL is correct.

echo "=== HTTP Benchmark: Home Page SSR ($REQUESTS requests, $CONCURRENCY concurrency) ==="
ab $AB_OPTS -n $REQUESTS -c $CONCURRENCY "$PROTOCOL://127.0.0.1:$PORT/" | grep -E "Requests per second|Time per request|99%"

echo "=== HTTP Benchmark: Search Logic (Jollof) ($REQUESTS requests, $CONCURRENCY concurrency) ==="
ab $AB_OPTS -n $REQUESTS -c $CONCURRENCY "$PROTOCOL://127.0.0.1:$PORT/?q=Jollof" | grep -E "Requests per second|Time per request|99%"

echo "=== Cleanup ==="
kill $SERVER_PID
rm -f .tester/db_tests/tmp_harness
echo "Done!"
