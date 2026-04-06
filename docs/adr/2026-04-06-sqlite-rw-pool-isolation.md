# ADR [013]: SQLite Read-Write Pool Isolation
**Date**: 2026-04-06 **Status**: Accepted

## Context
SQLite with WAL mode allows concurrent readers but only one concurrent writer. Our current single-pool implementation throttles read throughput when write locks are held, as the entire repository is restricted to `SetMaxOpenConns(1)` or suffers from locking contention in high-concurrency environments.

## Decision
We will isolate SQLite connection pools into a dedicated `writeDB` (highly constrained to 1 connection) and a `readDB` (liberal concurrency, e.g., 100 connections). Every repository method will be audited to target the correct pool based on whether it performs a mutation or a query.

## Consequences
This dramatically increases read throughput for the homepage and search feeds while maintaining SQLite's single-writer integrity. The codebase becomes slightly more complex by requiring explicit pool selection, and we must ensure both pools use consistent WAL pragmas to maintain data consistency.
