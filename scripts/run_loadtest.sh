#!/bin/bash
set -e

# Increase file descriptor limit to handle high concurrency
ulimit -n 250000 || echo "Warning: Could not set ulimit -n to 250000. Load test might fail due to "too many open files"."

DB_PATH=".tester/data/benchmark.db"

# Cleanup any previous benchmark DB
rm -f "$DB_PATH"
rm -f "${DB_PATH}-shm"
rm -f "${DB_PATH}-wal"

echo "=== Compiling CLI ==="
mkdir -p .tester/db_tests
go build -o .tester/db_tests/tmp_harness main.go

echo "=== Seeding 100,000 listings ==="
./.tester/db_tests/tmp_harness stress -c 100000 "$DB_PATH"

echo "=== Starting agbalumo server ==="
export DATABASE_URL="$DB_PATH"
./.tester/db_tests/tmp_harness serve > server.log 2>&1 &
SERVER_PID=$!

echo "Server started with PID: $SERVER_PID. Waiting for it to become ready..."
sleep 3

echo "=== Running k6 Load Test ==="
# Notice we disable TLS verify since it might be self-signed on localhost:8443
# TARGET_URL="https://localhost:8443" k6 run --insecure-skip-tls-verify scripts/benchmark_users.js

echo "=== Cleaning up ==="
kill $SERVER_PID || true
rm -f .tester/db_tests/tmp_harness
echo "Done! Check server.log for any server-side errors if needed."
