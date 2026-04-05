//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/jgautheron/goconst/cmd/goconst"
	_ "github.com/mibk/dupl"
	_ "github.com/uudashr/gocognit/cmd/gocognit"
	_ "github.com/zricethezav/gitleaks/v8"
	_ "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
