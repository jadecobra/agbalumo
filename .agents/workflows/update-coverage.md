---
description: Update test coverage threshold (ONLY INCREASE, NEVER LOWER)
---

# Update Test Coverage Workflow

> **CRITICAL RULE: NEVER lower the test coverage below its current value in `.agent/coverage-threshold`. If test coverage drops due to your changes, you MUST write new tests to cover the new or modified code until the threshold is met again. Lowering the threshold is strictly forbidden.**

Follow these steps to update the test coverage threshold for the project ONLY when you have ADDED tests and the coverage has INCREASED.

1. Find the current total test coverage percentage by executing the following command:
```sh
export PATH=/opt/homebrew/bin:/usr/local/go/bin:$PATH && go test -race -coverprofile=.tester/coverage/coverage.out ./... && go tool cover -func=.tester/coverage/coverage.out | grep -v "mock" | grep total | awk '{print substr($3, 1, length($3)-1)}'
```
*(Extract the coverage percentage from the output).*

// turbo
2. Update the file `.agent/coverage-threshold` with the new coverage percentage.

// turbo
3. Run the pre-commit script to verify the threshold is correctly set and all tests pass:
```sh
export PATH=/opt/homebrew/bin:/usr/local/go/bin:$PATH && task pre-commit
```

// turbo
4. Commit the change using:
```sh
git add .agent/coverage-threshold scripts/pre-commit.sh && git commit -m "build: update test coverage threshold to <NEW_PERCENTAGE>%"
```
