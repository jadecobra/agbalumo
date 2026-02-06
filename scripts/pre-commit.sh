#!/bin/sh
set -e

echo "Running 10x Engineer Quality Checks..."

echo "1. Running Go Vet..."
go vet ./...

echo "2. Running Tests..."
go test -count=1 ./...

echo "Quality Check Passed! 🚀"
