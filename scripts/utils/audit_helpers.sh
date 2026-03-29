#!/bin/sh
# scripts/utils/audit_helpers.sh
# Abstracted reporting and summary functions for audit scripts

# These functions expect FMT to be defined by the sourcing script
# and also expect output_json_envelope to be defined from utils.sh

# Overriding them slightly to track WARNINGS/FAILURES
WARNINGS=0
FAILURES=0
COLLECTED_WARNINGS=()
COLLECTED_FAILURES=()

pass() { if [ "$FMT" != "json" ] && [ "${VERBOSE:-0}" -eq 1 ]; then echo "${GREEN}  ✅ PASS:${NC} $1"; fi; }
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

# audit_summary
# Prints the final summary of the audit and exits with the correct status code.
# Usage: audit_summary "script_name.sh"
audit_summary() {
    local script_name="$1"

    if [ "$FMT" = "json" ]; then
        local combined_hints="[]"
        if [ ${#COLLECTED_WARNINGS[@]} -gt 0 ] || [ ${#COLLECTED_FAILURES[@]} -gt 0 ]; then
            # merge arrays for warnings field in JSON envelope
            combined_hints=$(printf '%s\n' "${COLLECTED_FAILURES[@]}" "${COLLECTED_WARNINGS[@]}" | jq -R . | jq -s .)
        fi

        if [ "$FAILURES" -eq 0 ] && [ "$WARNINGS" -eq 0 ]; then
            output_json_envelope true "$script_name" "🏆 All checks passed with no warnings!" "$combined_hints"
            exit 0
        elif [ "$FAILURES" -eq 0 ]; then
            output_json_envelope true "$script_name" "⚠️ $WARNINGS warning(s) found — no critical failures." "$combined_hints"
            exit 0
        else
            output_json_envelope false "$script_name" "❌ $FAILURES failure(s), $WARNINGS warning(s) found." "$combined_hints"
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
}
