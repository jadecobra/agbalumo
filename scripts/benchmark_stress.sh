#!/bin/bash
set -e

DB_PATH=".tester/data/benchmark.db"

# Cleanup any previous benchmark DB
rm -f "$DB_PATH"
rm -f "${DB_PATH}-shm"
rm -f "${DB_PATH}-wal"

echo "=== Compiling CLI ==="
go build -o tmp_harness main.go

echo "=== Write Benchmark (100k Inserts) ==="
time ./tmp_harness stress -c 100000 "$DB_PATH"

echo "=== Reconnect & Query DB ==="
echo "Total Listings Generated:"
sqlite3 "$DB_PATH" "SELECT count(*) FROM listings;"

echo "=== Basic Read Benchmark ==="
echo "Querying Page 1 (No Filters) - Top 20:"
time sqlite3 "$DB_PATH" "SELECT id FROM listings ORDER BY created_at DESC LIMIT 20 OFFSET 0;" > /dev/null

echo "Querying Page 500 (No Filters) - 20 items deep pagination:"
time sqlite3 "$DB_PATH" "SELECT id FROM listings ORDER BY created_at DESC LIMIT 20 OFFSET 10000;" > /dev/null

echo "Querying 'Business' Category - Top 20:"
time sqlite3 "$DB_PATH" "SELECT id FROM listings WHERE type = 'Business' ORDER BY created_at DESC LIMIT 20 OFFSET 0;" > /dev/null

echo "=== Cleanup ==="
rm tmp_harness
echo "Done!"
