#!/bin/sh
# scripts/security-check.sh
# Security validation script to maintain Grade A per docs/SECURITY_AUDIT.md
# This script checks for security violations that would lower the security grade

set -e

# Robust PATH discovery
. "$(dirname "$0")/utils.sh"
setup_path

FAILED=0
FULL_SCAN=0
if [ "$1" = "--full" ]; then
    FULL_SCAN=1
    echo "Running FULL Security Audit (all files)..."
else
    echo "Running Security Checks (staged changes)..."
fi

# Helper to check if a pattern exists in files (staged or all)
# Usage: check_pattern "pattern" "file_glob" "error_message"
check_pattern() {
    PATTERN="$1"
    GLOB="$2"
    MSG="$3"
    
    if [ $FULL_SCAN -eq 1 ]; then
        # Search all tracked files matching glob
        # We use git grep to avoid searching ignored files/vendor
        if git grep -E "$PATTERN" -- "$GLOB" > /dev/null 2>&1; then
            echo "${RED}❌ FAIL: $MSG${NC}"
            if [ "$4" != "quiet" ]; then
                git grep -n -E "$PATTERN" -- "$GLOB"
            fi
            return 1
        fi
    else
        # Search only added/modified lines in staged changes
        # We filter for lines starting with +
        if git diff --cached -- "$GLOB" | grep "^+" | grep -qE "$PATTERN"; then
            echo "${RED}❌ FAIL: $MSG${NC}"
            return 1
        fi
    fi
    return 0
}

# Helper to check if a file exists (staged or all) and matches a regex
# Usage: check_file_path "regex" "error_message"
check_file_path() {
    REGEX="$1"
    MSG="$2"
    
    if [ $FULL_SCAN -eq 1 ]; then
        if git ls-files | grep -E "$REGEX" > /dev/null 2>&1; then
            echo "${RED}❌ Secret/Artifact Leak: $MSG${NC}"
            git ls-files | grep -E "$REGEX"
            return 1
        fi
    else
        if git diff --cached --name-only | grep -E "$REGEX" > /dev/null 2>&1; then
            echo "${RED}❌ Secret/Artifact Leak: $MSG${NC}"
            git diff --cached --name-only | grep -E "$REGEX"
            return 1
        fi
    fi
    return 0
}

# ============================================
# CHECK 1: No inline scripts in templates
# ============================================
echo ""
echo "1. Checking for inline <script> tags in templates..."

INLINE_SCRIPT_FAILED=0
# Check all staged HTML files for inline script tags
# We ignore lines that start with <script src= (external scripts)
if [ $FULL_SCAN -eq 1 ]; then
    FILES=$(git ls-files 'ui/templates/**/*.html')
else
    FILES=$(git diff --cached --name-only -- 'ui/templates/**/*.html')
fi

for file in $FILES; do
    if [ $FULL_SCAN -eq 1 ]; then
        CONTENT=$(cat "$file")
    else
        CONTENT=$(git show :"$file")
    fi
    if echo "$CONTENT" | grep -v "<script src=" | grep -q "<script"; then
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
if check_pattern "onclick=" "*.go *.html" "Found onclick= handlers. Use hx-on:click instead."; then
    echo "${GREEN}✅ PASS: No onclick handlers found${NC}"
else
    FAILED=1
fi

# ============================================
# CHECK 3: No forbidden CDN domains
# ============================================
echo ""
echo "3. Checking for forbidden CDN domains..."

FORBIDDEN_CDNS="unpkg.com|cdn.jsdelivr.net|cdn.tailwindcss.com|jsdelivr.net"

if check_pattern "https?://($FORBIDDEN_CDNS)" "*.go *.html" "Forbidden CDN domains found. Self-host scripts in ui/static/js/"; then
    echo "${GREEN}✅ PASS: No forbidden CDN domains${NC}"
else
    FAILED=1
fi

# ============================================
# CHECK 4: CSP allows only permitted domains
# ============================================
echo ""
echo "4. Checking CSP configuration..."

CSP_FILE="internal/middleware/security.go"
if [ $FULL_SCAN -eq 1 ] || git diff --cached --name-only | grep -q "$CSP_FILE"; then
    if [ $FULL_SCAN -eq 1 ]; then
        CSP_CONTENT=$(cat "$CSP_FILE" 2>/dev/null)
    else
        CSP_CONTENT=$(git show :"$CSP_FILE" 2>/dev/null)
    fi
    
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

if check_pattern "(eval\(|Function\(|innerHTML\s*=)" "*.js *.html" "Found dangerous JavaScript patterns (eval, Function, innerHTML)"; then
    echo "${GREEN}✅ PASS: No dangerous JS patterns${NC}"
else
    FAILED=1
fi

# ============================================
# CHECK 6: SQL Injection Prevention
# ============================================
echo ""
echo "6. Checking for SQL Injection patterns..."

SQLI_FAILED=0
if ! check_pattern "Sprintf.*\b(SELECT|INSERT|UPDATE|DELETE)\b" "*.go" "Potential SQL Injection found using string formatting (Sprintf)"; then
    SQLI_FAILED=1
fi

if ! check_pattern "(Query|QueryRow|Exec)\(.*[+]" "*.go" "Potential SQL Injection found using string concatenation (+)"; then
    SQLI_FAILED=1
fi

if [ $SQLI_FAILED -eq 1 ]; then
    echo "   Use parameterized queries (?) instead of formatting strings into SQL."
    FAILED=1
else
    echo "${GREEN}✅ PASS: No obvious SQL injection patterns found${NC}"
fi

# ============================================
# CHECK 7: Cross-Site Scripting (XSS) in Go
# ============================================
echo ""
echo "7. Checking for XSS vulnerabilities (template.HTML)..."

if [ $FULL_SCAN -eq 1 ]; then
    # In full scan, we only warn if new patterns are found? 
    # Actually, let's keep it consistent with the original logic.
    if git grep -q 'template\.HTML(' -- "*.go"; then
        echo "${YELLOW}⚠️  WARNING: template.HTML() found in codebase${NC}"
    else
        echo "${GREEN}✅ PASS: No dangerous unescaped HTML patterns (template.HTML) found${NC}"
    fi
else
    if git diff --cached -- '*.go' | grep "^+" | grep -q 'template\.HTML('; then
        echo "${YELLOW}⚠️  WARNING: template.HTML() found in staged changes${NC}"
        echo "   Verify that the input is strictly sanitized. Only use this for trusted content like system icons."
    else
        echo "${GREEN}✅ PASS: No dangerous unescaped HTML patterns (template.HTML) found${NC}"
    fi
fi

# ============================================
# CHECK 8: CSRF in OAuth State
# ============================================
echo ""
echo "8. Checking for hardcoded OAuth state (Login CSRF)..."

if check_pattern 'GetAuthCodeURL\("[^"]+"' "*.go :!*_test.go" "Hardcoded static state found in GetAuthCodeURL"; then
    echo "${GREEN}✅ PASS: No hardcoded OAuth state found${NC}"
else
    FAILED=1
fi

# ============================================
# CHECK 9: Secret Scanner (Moved from pre-commit.sh)
# ============================================
echo ""
echo "9. Running Secret Scanner..."

# 1. Check filenames
if ! check_file_path "\.env$|\.pem$|\.key$|\.db$|\.db-shm$|\.db-wal$" "Sensitive file extension detected!"; then
    FAILED=1
fi

# 2. Check content
SECRET_PATTERN="[B]EGIN PRIVATE KEY|sk-[a-zA-Z0-9]{20,}|[A]Iza[a-zA-Z0-9_-]{30,}"
if [ $FULL_SCAN -eq 1 ]; then
    if git grep -I -n -E "$SECRET_PATTERN" -- '.' > /dev/null 2>&1; then
        echo "${RED}❌ Secret Leak: Sensitive pattern detected in codebase!${NC}"
        git grep -I -n -E "$SECRET_PATTERN" -- '.'
        FAILED=1
    fi
else
    if git grep --cached -I -n -E "$SECRET_PATTERN" -- '.' > /dev/null 2>&1; then
        echo "${RED}❌ Secret Leak: Sensitive pattern detected in staged content!${NC}"
        git grep --cached -I -n -E "$SECRET_PATTERN" -- '.'
        FAILED=1
    fi
fi

if [ $FAILED -eq 0 ]; then
    echo "${GREEN}✅ PASS: No obvious secrets or sensitive files found${NC}"
fi

# ============================================
# CHECK 10: No hardcoded secrets in Go/JS
# ============================================
echo ""
echo "10. Checking for hardcoded secrets (less strict, for general patterns)..."

# Check for common secret patterns
if check_pattern "(password|secret|key|token).*=\s*[\"'][^\"']{8,}[\"']" "*.go *.js *.html" "Possible hardcoded secret detected" "quiet"; then
    echo "${GREEN}✅ PASS: No obvious hardcoded secrets (general patterns)${NC}"
else
    echo "${YELLOW}⚠️  WARNING: Possible hardcoded secrets. Use environment variables for secrets.${NC}"
fi

# ============================================
# CHECK 11: Container Scan (if Docker modified)
# ============================================
echo ""
echo "11. Checking for container vulnerabilities..."

# Only run if Dockerfile or scripts are modified, or if forced (FULL_SCAN)
if [ $FULL_SCAN -eq 1 ] || git diff --cached --name-only | grep -qE "Dockerfile|scripts/|go.mod|go.sum"; then
    if command -v docker >/dev/null 2>&1; then
        echo "${YELLOW}Building and scanning container image...${NC}"
        # We reuse the repro script's logic but more integrated
        IMAGE_NAME="agbalumo-security-check"
        if docker build -q -t "$IMAGE_NAME" . >/dev/null 2>&1; then
            if docker run --rm \
                -v /var/run/docker.sock:/var/run/docker.sock \
                -v "$(pwd)/.cache/trivy:/root/.cache/" \
                aquasec/trivy:latest \
                image \
                --exit-code 1 \
                --severity CRITICAL,HIGH \
                --ignore-unfixed \
                "$IMAGE_NAME" >/dev/null 2>&1; then
                echo "${GREEN}✅ PASS: No critical/high vulnerabilities in container${NC}"
            else
                echo "${RED}❌ FAIL: Container vulnerabilities detected! Run scripts/repro_ci_failure.sh for details.${NC}"
                FAILED=1
            fi
        else
            echo "${YELLOW}⚠️  WARNING: Container build failed during security check, skipping scan${NC}"
        fi
    else
        echo "${YELLOW}⚠️  WARNING: Docker not found, skipping container scan${NC}"
    fi
else
    echo "   Skipping container scan (no relevant changes detected)"
fi

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
