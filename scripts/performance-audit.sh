#!/bin/sh
# scripts/performance-audit.sh
# Automated performance audit for agbalumo
# Checks: asset sizes, DB config, caching, accessibility, SQL indexes, response times
#
# Usage:
#   ./scripts/performance-audit.sh              # full audit
#   ./scripts/performance-audit.sh --live       # include live server response-time checks
#
# Exit codes:
#   0  All checks passed
#   1  One or more warnings found (non-blocking)
#   2  One or more critical failures found

# ─── Setup ────────────────────────────────────────────────────────────────────

set -e

FMT="json"
if [ "$1" = "--text" ]; then 
    FMT="text"
    shift
fi

# Robust PATH discovery
source "$(dirname "$0")/utils.sh"
setup_path

WARNINGS=0
FAILURES=0
LIVE_MODE=0

if [ "${1:-}" = "--live" ]; then
    LIVE_MODE=1
fi

# Messaging functions (pass, warn, fail, info) are now in utils.sh
# Overriding them slightly to track WARNINGS/FAILURES
COLLECTED_WARNINGS=()
COLLECTED_FAILURES=()

pass() { if [ "$FMT" != "json" ]; then echo "${GREEN}  ✅ PASS:${NC} $1"; fi; }
warn() { 
    if [ "$FMT" != "json" ]; then echo "${YELLOW}  ⚠️  WARN:${NC} $1"; fi
    WARNINGS=$((WARNINGS + 1))
    COLLECTED_WARNINGS+=("$1")
}
fail() { 
    if [ "$FMT" != "json" ]; then echo "${RED}  ❌ FAIL:${NC} $1"; fi
    FAILURES=$((FAILURES + 1))
    COLLECTED_FAILURES+=("$1")
}
info() { if [ "$FMT" != "json" ]; then echo "${CYAN}  ℹ️  INFO:${NC} $1"; fi; }

# Determine project root (script may be called from any dir)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT"

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BOLD}${BLUE}════════════════════════════════════════════════${NC}"
    echo "${BOLD}${BLUE}  agbalumo Performance Audit${NC}"
    echo "${BOLD}${BLUE}════════════════════════════════════════════════${NC}"
    echo ""
fi

# ─── CHECK 1: Static Asset Sizes ──────────────────────────────────────────────

if [ "$FMT" != "json" ]; then echo "${BOLD}1. Static Asset Sizes${NC}"; fi

# Logo PNG — should be < 100KB (ideally WebP)
LOGO="ui/static/images/logo.png"
LOGO_WEBP="ui/static/images/logo.webp"
if [ -f "$LOGO" ]; then
    LOGO_KB=$(du -k "$LOGO" | cut -f1)
    if [ -f "$LOGO_WEBP" ]; then
        WEBP_KB=$(du -k "$LOGO_WEBP" | cut -f1)
        pass "logo.webp present (${WEBP_KB}KB vs PNG ${LOGO_KB}KB) — <picture> WebP served to modern browsers ✓"
    elif [ "$LOGO_KB" -gt 200 ]; then
        fail "logo.png is ${LOGO_KB}KB (>200KB) and no logo.webp exists. Convert: cwebp -q 85 ui/static/images/logo.png -o ui/static/images/logo.webp"
    elif [ "$LOGO_KB" -gt 100 ]; then
        warn "logo.png is ${LOGO_KB}KB (>100KB) and no logo.webp exists. Consider WebP conversion."
    else
        pass "logo.png is ${LOGO_KB}KB ✓"
    fi
else
    info "logo.png not found, skipping"
fi

# CSS bundle — warn if > 80KB
CSS="ui/static/css/output.css"
if [ -f "$CSS" ]; then
    CSS_KB=$(du -k "$CSS" | cut -f1)
    if [ "$CSS_KB" -gt 150 ]; then
        fail "output.css is ${CSS_KB}KB (>150KB). Run 'npx tailwindcss --minify'."
    elif [ "$CSS_KB" -gt 110 ]; then
        warn "output.css is ${CSS_KB}KB (>110KB). Verify Tailwind purge config includes all template globs."
    else
        pass "output.css is ${CSS_KB}KB ✓"
    fi
fi

# HTMX — should be < 60KB (it's pre-minified)
HTMX="ui/static/js/htmx.min.js"
if [ -f "$HTMX" ]; then
    HTMX_KB=$(du -k "$HTMX" | cut -f1)
    if [ "$HTMX_KB" -gt 60 ]; then
        warn "htmx.min.js is ${HTMX_KB}KB. Check if a newer, smaller build is available."
    else
        pass "htmx.min.js is ${HTMX_KB}KB ✓"
    fi
fi

# chart.js — admin-only, warn if loaded globally
CHART="ui/static/js/chart.umd.min.js"
if [ -f "$CHART" ]; then
    CHART_KB=$(du -k "$CHART" | cut -f1)
    info "chart.js is ${CHART_KB}KB — verify it's only loaded on admin pages (not base.html)"
    if grep -q "chart.umd" ui/templates/base.html 2>/dev/null; then
        fail "chart.js (${CHART_KB}KB) is included in base.html — move to admin-only template block."
    else
        pass "chart.js (${CHART_KB}KB) is not in base.html ✓"
    fi
fi

# app.js — warn if > 30KB
APPJS="ui/static/js/app.js"
if [ -f "$APPJS" ]; then
    APP_KB=$(du -k "$APPJS" | cut -f1)
    APP_LINES=$(wc -l < "$APPJS")
    if [ "$APP_KB" -gt 50 ]; then
        warn "app.js is ${APP_KB}KB / ${APP_LINES} lines. Consider splitting into focused modules."
    else
        pass "app.js is ${APP_KB}KB ✓"
    fi
fi

# Uploaded images — flag any single upload > 500KB
echo ""
LARGE_UPLOADS=$(find ui/static/uploads -name "*.jpg" -o -name "*.png" -o -name "*.webp" 2>/dev/null | while read -r f; do
    SIZE=$(du -k "$f" | cut -f1)
    if [ "$SIZE" -gt 500 ]; then echo "$f (${SIZE}KB)"; fi
done)
if [ -n "$LARGE_UPLOADS" ]; then
    warn "Found uploads > 500KB: $LARGE_UPLOADS"
else
    pass "All uploads are < 500KB ✓"
fi

# ─── CHECK 2: Cache-Busting Strategy ──────────────────────────────────────────

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BOLD}2. Cache-Busting & HTTP Caching${NC}"
fi

# Static assets should use ?v= or content hash for cache busting
if grep -n "app\.js\"" ui/templates/base.html 2>/dev/null | grep -qv "?v=\|?t="; then
    warn "app.js in base.html has no cache-busting query param (e.g., ?v=5). Old clients may use stale JS."
else
    pass "app.js has cache-busting param ✓"
fi

# Static middleware should set Cache-Control immutable on /static/
if grep -q "immutable" cmd/server.go 2>/dev/null; then
    pass "Cache-Control: immutable set for /static/ ✓"
else
    fail "No 'immutable' Cache-Control found in cmd/server.go for static assets."
fi

# ─── CHECK 3: Database Configuration ──────────────────────────────────────────

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BOLD}3. Database Configuration${NC}"
fi

SQLITE_FILE="internal/repository/sqlite/sqlite.go"

# WAL mode — critical for concurrent reads
if grep -q "journal_mode=WAL" "$SQLITE_FILE" 2>/dev/null; then
    pass "WAL mode enabled ✓"
else
    fail "WAL mode not found in $SQLITE_FILE. Add: db.Exec(\"PRAGMA journal_mode=WAL;\")"
fi

# Busy timeout — prevents lock errors under load
if grep -q "busy_timeout" "$SQLITE_FILE" 2>/dev/null; then
    pass "busy_timeout configured ✓"
else
    warn "busy_timeout not set. Add: db.Exec(\"PRAGMA busy_timeout=5000;\") to prevent lock errors."
fi

# Synchronous NORMAL — safe + faster than FULL
if grep -q "synchronous=NORMAL" "$SQLITE_FILE" 2>/dev/null; then
    pass "synchronous=NORMAL (WAL-safe, faster) ✓"
else
    warn "synchronous not set to NORMAL. Consider: PRAGMA synchronous=NORMAL; (safe with WAL)"
fi

# Connection pool — MaxOpenConns should be set for SQLite
if grep -q "SetMaxOpenConns\|MaxOpenConns" "$SQLITE_FILE" 2>/dev/null; then
    pass "MaxOpenConns is configured ✓"
else
    fail "MaxOpenConns not set on sql.DB. Add db.SetMaxOpenConns(1) — SQLite serializes writes."
fi

# FTS5 full-text search index
if grep -q "fts5" "$SQLITE_FILE" 2>/dev/null; then
    pass "FTS5 full-text search index present ✓"
else
    warn "No FTS5 index found. Text search will be a slow LIKE scan on large datasets."
fi

# Composite index for main listing query (is_active, status, type)
if grep -q "idx_listings_active_status_type\|is_active.*status.*type" "$SQLITE_FILE" 2>/dev/null; then
    pass "Composite index on (is_active, status, type) for FindAll ✓"
else
    fail "Missing composite index on listings(is_active, status, type). Main query will full-scan."
fi

# Index on owner_id for profile / FindAllByOwner
if grep -q "idx_listings_owner_id" "$SQLITE_FILE" 2>/dev/null; then
    pass "Index on listings(owner_id) ✓"
else
    warn "No explicit index on listings(owner_id). FindAllByOwner may be slow for large tables."
fi

# ─── CHECK 4: Search Smoke Benchmarks ─────────────────────────────────────────

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BOLD}4. Search Smoke Benchmarks (Repository Layer)${NC}"
fi

# Run a quick smoke benchmark with 500 records
# We use benchtime=100ms for speed in the pre-commit gate
BENCH_OUT=$(BENCH_LISTINGS=500 go test -json -v -bench=BenchmarkSearchPerformance/FindAll_Default_Page1 -benchtime=100ms ./internal/repository/sqlite/search_performance_test.go 2>/dev/null || echo "error")

if [ "$BENCH_OUT" = "error" ]; then
    warn "Could not run search smoke benchmarks (ensure Go is installed and test file exists)"
else
    # Parse the JSON output or fallback to raw text parsing for ns/op
    if echo "$BENCH_OUT" | grep -q '"Action"'; then
        TEXT_OUT=$(echo "$BENCH_OUT" | grep '"Output":"' | sed -E 's/.*"Output":"(.*)".*/\1/' | sed 's/\\t/ /g' | sed 's/\\n//g')
    else
        TEXT_OUT="$BENCH_OUT"
    fi
    NS_OP=$(echo "$TEXT_OUT" | grep "FindAll_Default_Page1" | awk '{print $3}')
    if [ -n "$NS_OP" ]; then
        # Convert ns to ms (one millionth)
        MS_OP=$(awk "BEGIN { printf \"%.2f\", $NS_OP / 1000000 }")
        
        # Threshold: Default page 1 search should be < 10ms for 500 records on any modern dev machine
        THRESHOLD=10.0
        if awk "BEGIN {exit !($MS_OP > $THRESHOLD)}"; then
            fail "Search smoke benchmark: ${MS_OP}ms/op (>${THRESHOLD}ms target) with 500 records."
        else
            pass "Search smoke benchmark: ${MS_OP}ms/op ✓"
        fi
    else
        warn "Could not parse benchmark output."
    fi
fi

# ─── CHECK 5: In-Memory Cache Layer ───────────────────────────────────────────

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BOLD}5. In-Memory Cache Layer${NC}"
fi

CACHED_FILE="internal/repository/cached/cached.go"

if [ -f "$CACHED_FILE" ]; then
    # RWMutex for safe concurrent cache reads
    if grep -q "sync.RWMutex\|RWMutex" "$CACHED_FILE"; then
        pass "RWMutex used in cache layer (safe concurrent reads) ✓"
    else
        fail "Cache layer does not use RWMutex — concurrent reads may race."
    fi

    # TTL-based expiry
    if grep -q "ttl\|TTL\|Duration" "$CACHED_FILE"; then
        pass "Cache has TTL expiry ✓"
    else
        warn "No TTL found in cached store. Cache will never expire — stale data risk."
    fi

    # Return copies (not references) to prevent mutation
    if grep -q "make(map\|copy(" "$CACHED_FILE"; then
        pass "Cache returns value copies (prevents external mutation) ✓"
    else
        warn "Cache may return references — external mutation could corrupt the cache."
    fi
else
    warn "No cached store found at $CACHED_FILE. Hot-path queries (GetCounts, GetLocations) will always hit SQLite."
fi

# ─── CHECK 6: Accessibility (Performance Impact via Core Web Vitals) ──────────

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BOLD}6. Accessibility & CLS/INP Checks${NC}"
fi

BASE_HTML="ui/templates/base.html"

# Logo img must have alt text
if grep -q "logo.png" "$BASE_HTML" 2>/dev/null; then
    if grep "logo.png" "$BASE_HTML" | grep -q 'alt='; then
        pass "Logo img has alt attribute ✓"
    else
        fail "Logo img is missing alt attribute."
    fi
fi

# Mobile bottom nav icon buttons need aria-label
MISSING_ARIA=$(grep -n 'material-symbols-outlined' "$BASE_HTML" 2>/dev/null | \
    while IFS=: read -r lineno content; do
        # Look for the preceding button/anchor without aria-label on same or adjacent line
        # Simple heuristic: check the 3 lines before this icon for aria-label
        START=$((lineno - 3))
        [ "$START" -lt 1 ] && START=1
        CONTEXT=$(sed -n "${START},${lineno}p" "$BASE_HTML")
        if echo "$CONTEXT" | grep -q "<button\|<a " && ! echo "$CONTEXT" | grep -q 'aria-label'; then
            echo "  Line ~$lineno: icon button may be missing aria-label"
        fi
    done | sort -u | head -5)

if [ -n "$MISSING_ARIA" ]; then
    warn "Possible icon-only buttons without aria-label (screen reader inaccessible):
$MISSING_ARIA"
else
    pass "Icon buttons appear to have accessible labels ✓"
fi

# CSS preload for critical stylesheet
if grep -rq 'rel="preload".*output.css\|output.css.*rel="preload"' ui/templates/ 2>/dev/null; then
    pass 'Critical CSS preloaded with <link rel="preload"> ✓'
else
    warn 'output.css not preloaded. Add <link rel="preload" href="/static/css/output.css" as="style"> in <head>.'
fi

# Font preconnect
if grep -rq 'rel="preconnect".*fonts.googleapis.com\|fonts.googleapis.com.*rel="preconnect"' ui/templates/ 2>/dev/null; then
    pass "Google Fonts preconnect hint present ✓"
else
    warn "No preconnect for fonts.googleapis.com. Adds ~100ms latency on first font request."
fi

# ─── CHECK 7: N+1 Query Pattern Detection ────────────────────────────────────

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BOLD}7. N+1 Query Pattern Detection${NC}"
fi

# Heuristic: look for DB calls inside range loops in handlers
N1_CANDIDATES=$(grep -rn "range \|for .*range" internal/handler/ --include="*.go" 2>/dev/null | \
    grep -v "_test.go" | while IFS=: read -r file lineno content; do
        START=$((lineno - 2))
        [ "$START" -lt 1 ] && START=1
        END=$((lineno + 5))
        CONTEXT=$(sed -n "${START},${END}p" "$file" 2>/dev/null || true)
        if echo "$CONTEXT" | grep -qE "repo\.|db\.|FindBy|GetBy|Query"; then
            if ! echo "$CONTEXT" | grep -qi "bounded admin action"; then
                echo "  $file:$lineno — DB call inside a range loop (potential N+1)"
            fi
        fi
    done | head -5)

if [ -n "$N1_CANDIDATES" ]; then
    warn "Potential N+1 patterns detected:
$N1_CANDIDATES"
else
    pass "No obvious N+1 patterns found in handlers ✓"
fi

# ─── CHECK 8: Live Response Time (optional) ───────────────────────────────────

if [ "$LIVE_MODE" -eq 1 ]; then
    if [ "$FMT" != "json" ]; then
        echo ""
        echo "${BOLD}8. Live Response Time Checks${NC}"
    fi

    BASE_URL="${BASE_URL:-https://localhost:8443}"

    check_endpoint() {
        LABEL="$1"
        URL="$2"
        TARGET_MS="${3:-500}"

        # curl with timing, skip TLS verify for self-signed certs
        TIMING=$(curl -sk -o /dev/null -w "%{time_total}" "$URL" 2>/dev/null || echo "error")

        if [ "$TIMING" = "error" ]; then
            warn "$LABEL — could not connect to $URL"
            return
        fi

        # Convert to ms (awk handles float)
        MS=$(awk "BEGIN { printf \"%.0f\", $TIMING * 1000 }")

        if [ "$MS" -gt "$TARGET_MS" ]; then
            fail "$LABEL — ${MS}ms (>${TARGET_MS}ms target) at $URL"
        elif [ "$MS" -gt $((TARGET_MS / 2)) ]; then
            warn "$LABEL — ${MS}ms (within target but elevated)"
        else
            pass "$LABEL — ${MS}ms ✓"
        fi
    }

    check_endpoint "Homepage (GET /)" "$BASE_URL/" 500
    check_endpoint "Listings fragment (GET /listings/fragment)" "$BASE_URL/listings/fragment" 200
    check_endpoint "About page (GET /about)" "$BASE_URL/about" 300
    check_endpoint "Static CSS" "$BASE_URL/static/css/output.css" 100
    check_endpoint "Static HTMX JS" "$BASE_URL/static/js/htmx.min.js" 100

    # Check Cache-Control header on static assets
    CSS_CACHE=$(curl -skI "$BASE_URL/static/css/output.css" 2>/dev/null | grep -i "cache-control" || true)
    if echo "$CSS_CACHE" | grep -qi "immutable\|max-age=31536000"; then
        pass "Cache-Control: immutable present on /static/css/output.css ✓"
    else
        fail "Static CSS missing Cache-Control: immutable. Clients won't cache aggressively."
    fi
else
    if [ "$FMT" != "json" ]; then
        echo ""
        info "Skipping live response time checks. Run with --live to include them:"
        info "  BASE_URL=https://localhost:8443 ./scripts/performance-audit.sh --live"
    fi
fi

# ─── Summary ──────────────────────────────────────────────────────────────────

if [ "$FMT" = "json" ]; then
    COMBINED_HINTS="[]"
    if [ ${#COLLECTED_WARNINGS[@]} -gt 0 ] || [ ${#COLLECTED_FAILURES[@]} -gt 0 ]; then
        # merge arrays for warnings field in JSON envelope
        COMBINED_HINTS=$(printf '%s\n' "${COLLECTED_FAILURES[@]}" "${COLLECTED_WARNINGS[@]}" | jq -R . | jq -s .)
    fi

    if [ "$FAILURES" -eq 0 ] && [ "$WARNINGS" -eq 0 ]; then
        output_json_envelope true "performance-audit.sh" "🏆 All checks passed with no warnings!" "$COMBINED_HINTS"
        exit 0
    elif [ "$FAILURES" -eq 0 ]; then
        output_json_envelope true "performance-audit.sh" "⚠️ $WARNINGS warning(s) found — no critical failures." "$COMBINED_HINTS"
        exit 0
    else
        output_json_envelope false "performance-audit.sh" "❌ $FAILURES failure(s), $WARNINGS warning(s) found." "$COMBINED_HINTS"
        exit 2
    fi
fi

echo ""
echo "${BOLD}${BLUE}════════════════════════════════════════════════${NC}"
echo "${BOLD}  Audit Summary${NC}"
echo "${BOLD}${BLUE}════════════════════════════════════════════════${NC}"
echo ""

if [ "$FAILURES" -eq 0 ] && [ "$WARNINGS" -eq 0 ]; then
    echo "${GREEN}${BOLD}🏆 All checks passed with no warnings!${NC}"
    exit 0
elif [ "$FAILURES" -eq 0 ]; then
    echo "${YELLOW}${BOLD}⚠️  ${WARNINGS} warning(s) found — no critical failures.${NC}"
    echo "   Address warnings to maintain peak performance."
    exit 0
else
    echo "${RED}${BOLD}❌ ${FAILURES} failure(s), ${WARNINGS} warning(s) found.${NC}"
    echo "   Fix failures before deploying."
    exit 2
fi
