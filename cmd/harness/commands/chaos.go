package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jadecobra/agbalumo/internal/util"
	"github.com/spf13/cobra"
)

var (
	flagStateCorrupt bool
	flagEnvWipe      bool
	flagTestSabotage bool
)

func ChaosCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "chaos",
		Short: "Inject chaos into the harness environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagStateCorrupt {
				state, err := getState()
				if err != nil {
					return err
				}
				// Corrupt the signature
				state.Signature = state.Signature + "_CORRUPTED"

				// Standard library JSON marshal to bypass the automatic signature calculation in agent.SaveState
				b, _ := json.MarshalIndent(state, "", "  ")
				b = append(b, '\n')
				if err := util.SafeWriteFile(StateFile, b); err != nil {
					return err
				}
				fmt.Println("💀 State signature corrupted.")
			}

			if flagEnvWipe {
				tmpDir := ".tester/tmp"
				// Walk and remove what we can, ignore errors (e.g. read-only module cache)
				_ = filepath.Walk(tmpDir, func(path string, info os.FileInfo, wErr error) error {
					if wErr != nil {
						return nil
					}
					if path == tmpDir {
						return nil
					}
					// #nosec G304,G122 - Internal harness chaos tool intentionally performs destructive cleanup
					_ = os.RemoveAll(path)
					return nil
				})
				fmt.Println("🧹 .tester/tmp/ has been partially wiped (read-only files skipped).")
			}

			if flagTestSabotage {
				files, _ := filepath.Glob("*_test.go")
				if len(files) == 0 {
					return nil
				}

				// Randomly select 1-3 files, excluding core harness tests to maintain CI integrity
				count := 1
				if len(files) > 1 {
					count = 1 + int(time.Now().Unix()%3)
					if count > len(files) {
						count = len(files)
					}
				}

				sabotagedCount := 0
				for _, target := range files {
					if sabotagedCount >= count {
						break
					}

					// Skip core harness infrastructure to avoid blocking the development pipeline
					isHarnessTest := target == "chaos_test.go" || target == "root_test.go" ||
						target == "verify_test.go" || target == "commands_test.go" ||
						target == "cost_test.go" || target == "gate_test.go" ||
						target == "status_test.go" || target == "init_test.go" ||
						target == "update_coverage_test.go" || target == "set_phase_test.go"

					if isHarnessTest {
						continue
					}

					// #nosec G304 - Internal harness chaos tool intentionally sabotages test logic for stress testing
					content, err := util.SafeReadFile(target)
					if err != nil {
						continue
					}

					// Sabotage by appending error to the first test function
					lines := strings.Split(string(content), "\n")
					sabotaged := false
					for j, line := range lines {
						if strings.HasPrefix(line, "func Test") && strings.Contains(line, "(t *testing.T)") {
							lines[j] = line + "\n\tt.Errorf(\"CHAOS_SABOTAGE: Intentionally failing test\")"
							sabotaged = true
							break
						}
					}

					if sabotaged {
						// #nosec G304 - Internal harness chaos tool intentionally sabotages test logic
						if err := util.SafeWriteFile(target, []byte(strings.Join(lines, "\n"))); err != nil {
							return err
						}
						fmt.Printf("🔥 Sabotaged %s\n", target)
						sabotagedCount++
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&flagStateCorrupt, "state-corrupt", false, "Randomly alter signatures in state.json")
	cmd.Flags().BoolVar(&flagEnvWipe, "env-wipe", false, "Clean .tester/tmp/")
	cmd.Flags().BoolVar(&flagTestSabotage, "test-sabotage", false, "Temporarily inject logic failures into *_test.go files")

	return cmd
}
