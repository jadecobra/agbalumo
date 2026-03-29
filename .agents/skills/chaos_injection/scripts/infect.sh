#!/usr/bin/env bash
# Chaos Monkey Infestation Tool
# Mandate: Active Sabotage and Fault Injection

set -u

STATE_FILE=".agents/state.json"
TMP_DIR=".tester/tmp"

usage() {
    echo "Usage: $0 [OPTION]"
    echo "Options:"
    echo "  --state-corrupt    Randomly alter a signature in state.json"
    echo "  --env-wipe         Delete the .tester/tmp/ environment"
    echo "  --test-sabotage    Inject a False Positive into a validation test"
    echo "  --dry-run          Show what *would* happen"
    exit 1
}

if [[ $# -eq 0 ]]; then usage; fi

# --- Chaos Events ---

case "$1" in
    --state-corrupt)
        echo "[ChaosMonkey] Injecting State Corruption..."
        if [[ -f "$STATE_FILE" ]]; then
            # Replace a signature character to trigger tampering detection
            sed -i.bak 's/"signature": "/"signature": "X/' "$STATE_FILE"
            echo "[SUCCESS] State file corrupted (Backup created: ${STATE_FILE}.bak)"
            echo "[ACTION] Run 'go run scripts/verify-persona.go' to verify detection."
        else
            echo "[ERROR] State file not found."
            exit 1
        fi
        ;;

    --env-wipe)
        echo "[ChaosMonkey] Performing Environment Wipe..."
        if [[ -d "$TMP_DIR" ]]; then
            rm -rf "$TMP_DIR"
            echo "[SUCCESS] Environment wiped: $TMP_DIR"
            echo "[ACTION] Trigger any 'agent-exec.sh' task to verify auto-rebuild."
        else
            echo "[SKIP] $TMP_DIR already clean."
        fi
        ;;

    --test-sabotage)
        echo "[ChaosMonkey] Sabotaging Tests..."
        # Find a suitable test and inject 'return true'
        # Targeted: any test containing a 'Safe' check
        TARGET=$(grep -rl "SafeOpen" internal/util/ 2>/dev/null | head -n 1)
        if [[ -n "$TARGET" ]]; then
            echo "[INFECT] Modifying $TARGET..."
            cp "$TARGET" "${TARGET}.original"
            sed -i.bak 's/if err != nil {/if false { \/\/ CHAOS INJECTED/g' "$TARGET"
            echo "[SUCCESS] Logic failure injected into $TARGET."
            echo "[ACTION] Run 'go test ./internal/util/...' to verify detection failure."
        else
            echo "[ERROR] No suitable testing target found."
            exit 1
        fi
        ;;

    --dry-run)
        echo "[ChaosMonkey] DRY RUN - Chaos Events Ready:"
        echo " - Would corrupt $STATE_FILE"
        echo " - Would wipe $TMP_DIR"
        echo " - Would sabotage tests in internal/"
        ;;

    *)
        usage
        ;;
esac
