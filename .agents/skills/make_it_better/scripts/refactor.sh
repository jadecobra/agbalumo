#!/usr/bin/env bash
set -e

echo "[BackendEngineer] Executing make_it_better skill..."
go test -race ./...
task pre-commit
