package agent

import (
	"encoding/json"
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
	Feature      string    `json:"feature"`
	WorkflowType string    `json:"workflow_type"`
	Phase        string    `json:"phase"`
	Gates        Gates     `json:"gates"`
	UpdatedAt    time.Time `json:"updated_at"`
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

	return &state, nil
}

// SaveState marshals and writes the state struct to the JSON file,
// updating the UpdatedAt timestamp.
func SaveState(path string, state *State) error {
	state.UpdatedAt = time.Now().UTC()

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
