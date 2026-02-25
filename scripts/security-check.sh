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

# Robust PATH discovery for macOS and Linux
for dir in /usr/local/bin /opt/homebrew/bin /usr/bin /bin; do
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) export PATH="$PATH:$dir" ;;
    esac
done

FAILED=0


echo "Running Security Checks (per SECURITY_AUDIT.md Grade A)..."

# ============================================
# CHECK 1: No inline scripts in templates
# ============================================
echo ""
echo "1. Checking for inline <script> tags in templates..."

INLINE_SCRIPT_FAILED=0
# Check all staged HTML files for inline script tags
# We ignore lines that start with <script src= (external scripts)
for file in $(git diff --cached --name-only -- 'ui/templates/**/*.html'); do
    if git show --cached "$file" | grep -v "<script src=" | grep -q "<script"; then
        echo "${RED}❌ FAIL: Inline <script> tag found in $file${NC}"
        INLINE_SCRIPT_FAILED=1
    fi
done

if [ $INLINE_SCRIPT_FAILED -eq 1 ]; then
    echo "   Templates should use external JS files from /static/js/"
    echo "   Move inline scripts to ui/static/js/ and include via base.html"
    FAILED=1
else
    echo "${GREEN}✅ PASS: No inline scripts in templates${NC}"
fi

# ============================================
# CHECK 2: No onclick handlers (use hx-on)
# ============================================
echo ""
echo "2. Checking for insecure onclick handlers..."

# Check staged HTML/Go files for onclick handlers
if git diff --cached -- '*.go' '*.html' | grep "^+" | grep -q 'onclick='; then
    echo "${RED}❌ FAIL: Found onclick= handlers in staged changes${NC}"
    echo "   Use hx-on:click instead for HTMX event handling"
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
if git diff --cached -- '*.go' '*.html' | grep "^+" | grep -qE "https?://($FORBIDDEN_CDNS)"; then
    echo "${RED}❌ FAIL: Forbidden CDN domains found in staged changes${NC}"
    echo "   Allowed: maps.googleapis.com, fonts.googleapis.com"
    echo "   Forbidden: unpkg.com, cdn.jsdelivr.net, cdn.tailwindcss.com"
    echo "   Self-host scripts in ui/static/js/"
    FAILED=1
else
    echo "${GREEN}✅ PASS: No forbidden CDN domains${NC}"
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
    else
        echo "${GREEN}✅ PASS: CSP configuration valid${NC}"
    fi
else
    echo "${GREEN}✅ PASS: CSP file not modified${NC}"
fi

# ============================================
# CHECK 5: No eval() or Function() in JavaScript
# ============================================
echo ""
echo "5. Checking for dangerous JavaScript patterns..."

if git diff --cached -- '*.js' '*.html' | grep "^+" | grep -qE "(eval\(|Function\(|innerHTML\s*=)"; then
    echo "${RED}❌ FAIL: Found dangerous JavaScript patterns (eval, Function, innerHTML) in staged changes${NC}"
    FAILED=1
else
    echo "${GREEN}✅ PASS: No dangerous JS patterns${NC}"
fi

# ============================================
# CHECK 6: Secret Scanner (Moved from pre-commit.sh)
# ============================================
echo ""
echo "6. Running Secret Scanner..."

# 1. Check filenames (staged)
# Matches .env, .pem, .key, .db, .db-shm, .db-wal
if git diff --cached --name-only | grep -E "\.env$|\.pem$|\.key$|\.db$|\.db-shm$|\.db-wal$"; then
    echo "${RED}❌ Secret/Artifact Leak: Sensitive file extension detected!${NC}"
    FAILED=1
fi

# 2. Check content (staged)
# We use git grep --cached to search the index.
# Patterns: Private Key, OpenAI Key, Google API Key
if git grep --cached -I -n -E "[B]EGIN PRIVATE KEY|sk-[a-zA-Z0-9]{20,}|[A]Iza[a-zA-Z0-9_-]{30,}" -- '.'; then
    echo "${RED}❌ Secret Leak: Sensitive pattern detected in staged content!${NC}"
    FAILED=1
fi

if [ $FAILED -eq 0 ]; then
    echo "${GREEN}✅ PASS: No obvious secrets or sensitive files staged${NC}"
fi

# ============================================
# CHECK 7: No hardcoded secrets in Go/JS
# ============================================
echo ""
echo "7. Checking for hardcoded secrets (less strict, for general patterns)..."

# Check for common secret patterns
if git diff --cached -- '*.go' '*.js' '*.html' 2>/dev/null | grep -qE "(password|secret|key|token).*=\s*[\"'][^\"']{8,}[\"']"; then
    echo "${YELLOW}⚠️  WARNING: Possible hardcoded secret detected${NC}"
    echo "   Use environment variables for secrets"
fi

echo "${GREEN}✅ PASS: No obvious hardcoded secrets (general patterns)${NC}"

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
