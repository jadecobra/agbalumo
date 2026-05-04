#!/bin/bash

# Accept count as argument, default to 100000
COUNT=${1:-100000}
CONCURRENCY=${2:-50}
REQUESTS=${3:-500}

DB_PATH=".tester/data/benchmark.db"
PORT=8081 
LOG_FILE=".tester/benchmark_server.log"

# Cleanup any previous benchmark DB and logs
rm -f "$DB_PATH"
rm -f "${DB_PATH}-shm"
rm -f "${DB_PATH}-wal"
rm -f "$LOG_FILE"

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

echo "Starting server on port $PORT (forcing HTTP via AGBALUMO_ENV=production)..."
DATABASE_URL="$DB_PATH" AGBALUMO_ENV=production PORT=$PORT ./.tester/db_tests/tmp_harness serve > "$LOG_FILE" 2>&1 &
SERVER_PID=$!

# Wait for server to be ready
echo "Waiting for server to respond on http://127.0.0.1:$PORT/ ..."
RETRY_COUNT=0
MAX_RETRIES=20
while ! curl -s http://127.0.0.1:$PORT/ > /dev/null; do
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

echo "=== HTTP Benchmark: Home Page SSR ($REQUESTS requests, $CONCURRENCY concurrency) ==="
ab -n $REQUESTS -c $CONCURRENCY "http://127.0.0.1:$PORT/" | grep -E "Requests per second|Time per request|99%"

echo "=== HTTP Benchmark: Search Logic (Jollof) ($REQUESTS requests, $CONCURRENCY concurrency) ==="
ab -n $REQUESTS -c $CONCURRENCY "http://127.0.0.1:$PORT/?q=Jollof" | grep -E "Requests per second|Time per request|99%"

echo "=== Cleanup ==="
kill $SERVER_PID
rm -f .tester/db_tests/tmp_harness
echo "Done!"
