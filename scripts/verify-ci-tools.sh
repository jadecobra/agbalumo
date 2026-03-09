#!/bin/bash
# scripts/verify-ci-tools.sh
# Verifies that CI workflow tools are open-source or properly authenticated.

set -e

RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
NC=$(printf '\033[0m')

CI_FILE=".github/workflows/ci.yml"

echo "Verifying CI tools in $CI_FILE..."

# Check for Docker Scout (known to fail without entitlement/auth)
if grep -q "docker/scout-action" "$CI_FILE"; then
    echo "${RED}❌ FAIL: Proprietary tool 'docker/scout-action' found without confirmed authentication.${NC}"
    echo "   Docker Scout requires Docker Hub authentication and entitlement."
    exit 1
fi

# Confirm Trivy is used (our preferred open-source alternative)
if grep -q "aquasecurity/trivy-action" "$CI_FILE"; then
    echo "${GREEN}✅ PASS: Using Trivy for container scanning (Open Source, local-friendly).${NC}"
else
    echo "${YELLOW}⚠️  WARNING: No container scanner detected in CI (expected Trivy).${NC}"
fi

echo "${GREEN}✅ CI Toolset Verification Passed${NC}"
exit 0
