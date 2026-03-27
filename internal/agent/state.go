package agent

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jadecobra/agbalumo/internal/util"
)

type GateStatus string

const (
	GatePending GateStatus = "PENDING"
	GatePassed  GateStatus = "PASSED"
	GateFailed  GateStatus = "FAILED"
)

// Standard gate IDs
const (
	GateRedTest             = "red-test"
	GateApiSpec             = "api-spec"
	GateImplementation      = "implementation"
	GateLint                = "lint"
	GateCoverage            = "coverage"
	GateBrowserVerification  = "browser-verification"
	GateTemplateDrift       = "template-drift"
)


// Standard workflow types
const (
	WorkflowFeature  = "feature"
	WorkflowBugfix   = "bugfix"
	WorkflowRefactor = "refactor"
)

// Gates represents the state of various validation checks.
type Gates struct {
	RedTest             GateStatus `json:"red-test"`
	ApiSpec             GateStatus `json:"api-spec"`
	Implementation      GateStatus `json:"implementation"`
	Lint                GateStatus `json:"lint"`
	Coverage            GateStatus `json:"coverage"`
	BrowserVerification GateStatus `json:"browser-verification"`
	TemplateDrift       GateStatus `json:"template-drift"`
}

// State represents the contents of .agents/state.json
type State struct {
	Warning      string    `json:"_DO_NOT_EDIT_"`
	Feature      string    `json:"feature"`
	WorkflowType string    `json:"workflow_type"`
	Phase        string    `json:"phase"`
	Gates        Gates     `json:"gates"`
	UpdatedAt    time.Time `json:"updated_at"`
	Signature    string    `json:"signature"`
}

func calculateSignature(s *State) string {
	copy := *s
	copy.Signature = "" // exclude signature itself from hash
	
	// predictable hashing by marshalling
	b, _ := json.Marshal(copy)
	hash := sha256.Sum256(b)
	return fmt.Sprintf("%x", hash)
}

// LoadState reads and parses the JSON state file.
func LoadState(path string) (*State, error) {
	// #nosec G304 - Internal harness state file
	b, err := util.SafeReadFile(path)
	if err != nil {
		return nil, err
	}

	var state State
	if unmarshalErr := json.Unmarshal(b, &state); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	// Validate structural equivalence to prevent case-insensitive JSON bypasses
	canonicalBytes, err := json.MarshalIndent(state, "", "  ")
	if err == nil {
		canonicalBytes = append(canonicalBytes, '\n')
		if string(b) != string(canonicalBytes) {
			return nil, errors.New("ANTI-CHEAT TRIGGERED: Manual modification of .agents/state.json detected (structural mismatch). You must use the ./scripts/agent-exec.sh commands to manage state")
		}
	}

	// Validate signature to prevent manual edits
	if state.Signature != "" {
		expected := calculateSignature(&state)
		if state.Signature != expected {
			return nil, errors.New("ANTI-CHEAT TRIGGERED: Manual modification of .agents/state.json detected. You must use the ./scripts/agent-exec.sh commands to manage state")
		}
	}

	return &state, nil
}

// SaveState marshals and writes the state struct to the JSON file,
// updating the UpdatedAt timestamp and Signature.
func SaveState(path string, state *State) error {
	state.Warning = "MANUAL EDITS WILL INVALIDATE SIGNATURE. USE ./scripts/agent-exec.sh TO MANAGE STATE."
	state.UpdatedAt = time.Now().UTC()
	state.Signature = calculateSignature(state)

	b, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	// Make sure we end with a newline for consistency.
	b = append(b, '\n')

	return util.SafeWriteFile(path, b)
}

// IsNotExist is a helper utility mapped to util.SafeIsNotExist
func IsNotExist(err error) bool {
	return util.SafeIsNotExist(err)
}
