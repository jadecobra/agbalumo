package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type handlerConcurrencyTestCase struct {
	name          string
	content       string
	matchMessage  string
	wantViolation bool
}

func getHandlerConcurrencyCases() []handlerConcurrencyTestCase {
	return []handlerConcurrencyTestCase{
		{
			name: "unbounded goroutine loop inside handler",
			content: `package sample

import (
	"context"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleListings(c echo.Context) error {
	items := []string{"a", "b", "c"}
	for _, item := range items {
		go func(it string) {
			_ = h.App.DB.FindByID(context.Background(), it)
		}(item)
	}
	return nil
}
`,
			wantViolation: true,
			matchMessage:  "unbounded goroutine spawn in loop",
		},
		{
			name: "goroutine fan-out exceeding concurrency limit",
			content: `package sample

import (
	"sync"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleHome(c echo.Context) error {
	var wg sync.WaitGroup
	wg.Add(5)
	go func() { defer wg.Done() }()
	go func() { defer wg.Done() }()
	go func() { defer wg.Done() }()
	go func() { defer wg.Done() }()
	go func() { defer wg.Done() }()
	wg.Wait()
	return nil
}
`,
			wantViolation: true,
			matchMessage:  "exceeds concurrency limit",
		},
		{
			name: "bounded goroutine fan-out within budget",
			content: `package sample

import (
	"context"
	"sync"
	"github.com/labstack/echo/v4"
)

func (h *Handler) fetchHomeData(ctx context.Context, c echo.Context) error {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		_ = h.App.DB.FindAll(ctx, "", "", "", 0, 0, 0, "", "", false, 10, 0)
	}()
	go func() {
		defer wg.Done()
		_ = h.App.DB.GetFeaturedListings(ctx, "Food", "")
	}()
	go func() {
		defer wg.Done()
		_ = h.App.DB.GetLocations(ctx)
	}()
	wg.Wait()
	return nil
}
`,
			wantViolation: false,
		},
		{
			name: "standard sequential handler without goroutines",
			content: `package sample

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) HandleDetail(c echo.Context) error {
	id := c.Param("id")
	_, err := h.App.DB.FindByID(c.Request().Context(), id)
	return err
}
`,
			wantViolation: false,
		},
	}
}

func runHandlerConcurrencyTestCase(t *testing.T, tc handlerConcurrencyTestCase) {
	t.Helper()
	tmpDir := t.TempDir()
	moduleDir := filepath.Join(tmpDir, "internal", "module", "sample")
	if err := os.MkdirAll(moduleDir, 0750); err != nil {
		t.Fatalf("failed to create module dir: %v", err)
	}

	filePath := filepath.Join(moduleDir, "handler.go")
	if err := os.WriteFile(filePath, []byte(tc.content), 0600); err != nil {
		t.Fatalf("failed to write fixture file: %v", err)
	}

	violations, err := CheckHandlerConcurrency(tmpDir, HandlerConcurrencyOptions{
		MaxGoroutinesPerFunc: 4,
	})
	if err != nil {
		t.Fatalf("unexpected error from CheckHandlerConcurrency: %v", err)
	}

	assertConcurrencyViolationExpectations(t, tc, violations)
}

func assertConcurrencyViolationExpectations(t *testing.T, tc handlerConcurrencyTestCase, violations []HandlerConcurrencyViolation) {
	t.Helper()
	if !tc.wantViolation {
		if len(violations) > 0 {
			t.Fatalf("expected no violations, got %d: %v", len(violations), violations)
		}
		return
	}

	if len(violations) == 0 {
		t.Fatalf("expected violation, got none")
	}
	if tc.matchMessage != "" && !hasMatchingViolationMessage(violations, tc.matchMessage) {
		t.Errorf("expected violation message containing %q, got: %v", tc.matchMessage, violations)
	}
}

func hasMatchingViolationMessage(violations []HandlerConcurrencyViolation, target string) bool {
	for _, v := range violations {
		if strings.Contains(v.Message, target) {
			return true
		}
	}
	return false
}

func TestCheckHandlerConcurrency(t *testing.T) {
	cases := getHandlerConcurrencyCases()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			runHandlerConcurrencyTestCase(t, tc)
		})
	}
}

func TestCheckHandlerConcurrency_RepoRoot(t *testing.T) {
	violations, err := CheckHandlerConcurrency("../..", DefaultHandlerConcurrencyOptions())
	if err != nil {
		t.Fatalf("failed to check repository root: %v", err)
	}

	if len(violations) > 0 {
		t.Fatalf("expected 0 violations in internal/module, found %d: %+v", len(violations), violations)
	}
}
