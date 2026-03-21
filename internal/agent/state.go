package agent

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type GateStatus string

const (
	GatePending GateStatus = "PENDING"
	GatePassed  GateStatus = "PASSED"
	GateFailed  GateStatus = "FAILED"
)

// Gates represents the state of various validation checks.
type Gates struct {
	RedTest             GateStatus `json:"red-test"`
	ApiSpec             GateStatus `json:"api-spec"`
	Implementation      GateStatus `json:"implementation"`
	Lint                GateStatus `json:"lint"`
	Coverage            GateStatus `json:"coverage"`
	BrowserVerification GateStatus `json:"browser-verification"`
}

// State represents the contents of .agent/state.json
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
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state State
	if err := json.Unmarshal(b, &state); err != nil {
		return nil, err
	}

	// Validate signature to prevent manual edits
	if state.Signature != "" {
		expected := calculateSignature(&state)
		if state.Signature != expected {
			return nil, errors.New("ANTI-CHEAT TRIGGERED: Manual modification of .agent/state.json detected. You must use the ./scripts/agent-exec.sh commands to manage state")
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

	return os.WriteFile(path, b, 0644)
}

// IsNotExist is a helper utility mapped to os.IsNotExist
func IsNotExist(err error) bool {
	return os.IsNotExist(err)
}
