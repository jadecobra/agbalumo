package agent

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type Taskfile struct {
	Tasks map[string]Task `yaml:"tasks"`
}

type Task struct {
	Deps   []interface{} `yaml:"deps"`
	Status []interface{} `yaml:"status"`
	Cmds   []interface{} `yaml:"cmds"`
}

func TestTaskfileOptimization(t *testing.T) {
	data, err := os.ReadFile("../../Taskfile.yml")
	require.NoError(t, err)

	var tf Taskfile
	err = yaml.Unmarshal(data, &tf)
	require.NoError(t, err)

	tests := []struct {
		name        string
		targetTasks []string
		installTask string
		binary      string
	}{
		{
			name:        "vulncheck",
			targetTasks: []string{"vulncheck", "ci:vulncheck"},
			installTask: "vulncheck:install",
			binary:      "govulncheck",
		},
		{
			name:        "lint",
			targetTasks: []string{"lint", "ci:lint"},
			installTask: "lint:install",
			binary:      "golangci-lint",
		},
		{
			name:        "gitleaks",
			targetTasks: []string{"gitleaks"},
			installTask: "gitleaks:install",
			binary:      "gitleaks",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Verify install task
			installTask, exists := tf.Tasks[tc.installTask]
			require.True(t, exists, tc.installTask+" task must exist")

			foundStatus := false
			for _, si := range installTask.Status {
				if s, ok := si.(string); ok && strings.Contains(s, "test -f") && strings.Contains(s, tc.binary) {
					foundStatus = true
					break
				}
			}
			assert.True(t, foundStatus, tc.installTask+" must have a status check for binary existence")

			// Verify target tasks depend on install task
			for _, target := range tc.targetTasks {
				task, exists := tf.Tasks[target]
				require.True(t, exists, target+" task must exist")

				foundDep := false
				for _, d := range task.Deps {
					if ds, ok := d.(string); ok && ds == tc.installTask {
						foundDep = true
						break
					} else if dm, ok := d.(map[string]interface{}); ok {
						if depName, ok := dm["task"].(string); ok && depName == tc.installTask {
							foundDep = true
							break
						}
					}
				}
				assert.True(t, foundDep, target+" must depend on "+tc.installTask)
			}
		})
	}
}
