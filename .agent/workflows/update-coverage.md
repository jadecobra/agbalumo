---
description: Update test coverage threshold
---

# Update Test Coverage Workflow

Follow these steps to update the test coverage threshold for the project.

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
export PATH=/opt/homebrew/bin:/usr/local/go/bin:$PATH && ./scripts/pre-commit.sh
```

// turbo
4. Commit the change using:
```sh
git add .agent/coverage-threshold scripts/pre-commit.sh && git commit -m "build: update test coverage threshold to <NEW_PERCENTAGE>%"
```
