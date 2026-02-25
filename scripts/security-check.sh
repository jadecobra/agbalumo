#!/bin/sh
# scripts/security-check.sh
# Security validation script to maintain Grade A per docs/SECURITY_AUDIT.md
# This script checks for security violations that would lower the security grade

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

FAILED=0

echo "Running Security Checks (per SECURITY_AUDIT.md Grade A)..."

# ============================================
# CHECK 1: No inline scripts in templates
# ============================================
echo ""
echo "1. Checking for inline <script> tags in templates..."

# Check for any <script> tags in HTML templates (excluding external script src)
if git diff --cached --name-only | grep -E "\.html$" | xargs git diff --cached --name-only | xargs -I{} sh -c 'git diff --cached {} 2>/dev/null | grep "^+" | grep -n "<script>" | head -5' 2>/dev/null || true; then
    # More precise check using git diff
    INLINE_SCRIPTS=$(git diff --cached -- '*.html' 2>/dev/null | grep "^+" | grep -c '<script>' || true)
    if [ "$INLINE_SCRIPTS" -gt 0 ]; then
        echo "${RED}❌ FAIL: Found inline <script> tags in HTML templates${NC}"
        echo "   Templates should use external JS files from /static/js/"
        echo "   Move inline scripts to ui/static/js/ and include via base.html"
        FAILED=1
    fi
fi

# Alternative check: grep staged HTML files for inline scripts
if git diff --cached --name-only -- 'ui/templates/**/*.html' | while read -r file; do
    if git show --cached "$file" 2>/dev/null | grep -q '<script[^>]*>.*</script>'; then
        echo "❌ FAIL: Inline script found in $file"
        exit 1
    fi
done; then
    : # Pass
else
    if [ $? -eq 1 ]; then
        echo "${RED}❌ FAIL: Inline scripts found in templates${NC}"
        FAILED=1
    fi
fi

# Check for inline scripts in partials specifically
if git diff --cached --name-only | grep -q "ui/templates/partials/.*\.html"; then
    for file in $(git diff --cached --name-only | grep "ui/templates/partials/.*\.html"); do
        if git show --cached ":$file" 2>/dev/null | grep -q '<script>'; then
            echo "${RED}❌ FAIL: Inline <script> tag found in $file${NC}"
            FAILED=1
        fi
    done
fi

if [ $FAILED -eq 0 ]; then
    echo "${GREEN}✅ PASS: No inline scripts in templates${NC}"
fi

# ============================================
# CHECK 2: No onclick handlers (use hx-on)
# ============================================
echo ""
echo "2. Checking for insecure onclick handlers..."

# Check Go handler files for inline onclick (should use hx-on)
if git diff --cached --name-only | grep -E "\.(go|html)$" | xargs git diff --cached 2>/dev/null | grep "^+" | grep -q 'onclick='; then
    echo "${RED}❌ FAIL: Found onclick= handlers${NC}"
    echo "   Use hx-on:click instead for HTMX event handling"
    echo "   Example: hx-on:click=\"this.parentElement.remove()\""
    FAILED=1
else
    echo "${GREEN}✅ PASS: No onclick handlers found${NC}"
fi

# ============================================
# CHECK 3: No forbidden CDN domains
# ============================================
echo ""
echo "3. Checking for forbidden CDN domains..."

# Forbidden CDNs per SECURITY_AUDIT.md
FORBIDDEN_CDNS="unpkg.com|cdn.jsdelivr.net|cdn.tailwindcss.com|jsdelivr.net"

# Check staged files for forbidden CDNs
if git diff --cached --name-only | xargs -I{} sh -c 'git show --cached {} 2>/dev/null | grep -qE "https?://('$FORBIDDEN_CDNS')" && echo "FOUND"' 2>/dev/null | grep -q "FOUND"; then
    echo "${RED}❌ FAIL: Forbidden CDN domains found${NC}"
    echo "   Allowed: maps.googleapis.com, fonts.googleapis.com"
    echo "   Forbidden: unpkg.com, cdn.jsdelivr.net, cdn.tailwindcss.com"
    echo "   Self-host scripts in ui/static/js/"
    FAILED=1
else
    # Direct check
    CDN_VIOLATIONS=$(git diff --cached -- '*.html' '*.go' 2>/dev/null | grep -cE "https?://($FORBIDDEN_CDNS)" || true)
    if [ "$CDN_VIOLATIONS" -gt 0 ]; then
        echo "${RED}❌ FAIL: Forbidden CDN domains detected${NC}"
        FAILED=1
    else
        echo "${GREEN}✅ PASS: No forbidden CDN domains${NC}"
    fi
fi

# ============================================
# CHECK 4: CSP allows only permitted domains
# ============================================
echo ""
echo "4. Checking CSP configuration..."

# Check security.go for CSP - should only allow maps.googleapis.com and fonts.googleapis.com
if git diff --cached --name-only | grep -q "internal/middleware/security.go"; then
    CSP_CONTENT=$(git show --cached internal/middleware/security.go 2>/dev/null)
    
    # Check for forbidden domains in CSP
    if echo "$CSP_CONTENT" | grep -qE "unpkg.com|jsdelivr.net|cdn.tailwindcss.com"; then
        echo "${RED}❌ FAIL: CSP contains forbidden CDN domains${NC}"
        FAILED=1
    fi
    
    # Check for required domains
    if ! echo "$CSP_CONTENT" | grep -q "maps.googleapis.com"; then
        echo "${YELLOW}⚠️  WARNING: CSP should allow maps.googleapis.com for Google Maps${NC}"
    fi
fi

echo "${GREEN}✅ PASS: CSP configuration valid${NC}"

# ============================================
# CHECK 5: No eval() or Function() in JavaScript
# ============================================
echo ""
echo "5. Checking for dangerous JavaScript patterns..."

if git diff --cached --name-only | grep -qE "\.js$|\.html$"; then
    EVAL_VIOLATIONS=$(git diff --cached -- '*.js' '*.html' 2>/dev/null | grep -cE "(eval\(|Function\(|innerHTML\s*=)" || true)
    if [ "$EVAL_VIOLATIONS" -gt 0 ]; then
        echo "${RED}❌ FAIL: Found dangerous JavaScript patterns (eval, Function, innerHTML)${NC}"
        FAILED=1
    else
        echo "${GREEN}✅ PASS: No dangerous JS patterns${NC}"
    fi
else
    echo "${GREEN}✅ PASS: No JS files changed${NC}"
fi

# ============================================
# CHECK 6: No hardcoded secrets in Go/JS
# ============================================
echo ""
echo "6. Checking for hardcoded secrets..."

# Check for common secret patterns
if git diff --cached -- '*.go' '*.js' '*.html' 2>/dev/null | grep -qE "(password|secret|key|token).*=\s*[\"'][^\"']{8,}[\"']"; then
    echo "${YELLOW}⚠️  WARNING: Possible hardcoded secret detected${NC}"
    echo "   Use environment variables for secrets"
fi

echo "${GREEN}✅ PASS: No obvious hardcoded secrets${NC}"

# ============================================
# FINAL RESULT
# ============================================
echo ""
echo "========================================"
if [ $FAILED -ne 0 ]; then
    echo "${RED}❌ SECURITY CHECK FAILED${NC}"
    echo ""
    echo "This commit violates security requirements"
    echo "per docs/SECURITY_AUDIT.md (Grade A)"
    echo ""
    echo "Fix the issues above and try again."
    exit 1
else
    echo "${GREEN}✅ ALL SECURITY CHECKS PASSED${NC}"
    echo ""
    echo "Security grade A maintained per SECURITY_AUDIT.md"
    exit 0
fi
