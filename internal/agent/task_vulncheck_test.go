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

func TestTaskfileVulncheckOptimization(t *testing.T) {
	data, err := os.ReadFile("../../Taskfile.yml")
	require.NoError(t, err)

	var tf Taskfile
	err = yaml.Unmarshal(data, &tf)
	require.NoError(t, err)

	// Verify vulncheck:install task
	installTask, exists := tf.Tasks["vulncheck:install"]
	require.True(t, exists, "vulncheck:install task must exist")
	
	foundStatus := false
	for _, si := range installTask.Status {
		if s, ok := si.(string); ok && strings.Contains(s, "test -f") && strings.Contains(s, "govulncheck") {
			foundStatus = true
			break
		}
	}
	assert.True(t, foundStatus, "vulncheck:install must have a status check for binary existence")

	// Verify vulncheck task
	vulncheck, exists := tf.Tasks["vulncheck"]
	require.True(t, exists, "vulncheck task must exist")
	
	foundDep := false
	for _, d := range vulncheck.Deps {
		if ds, ok := d.(string); ok && ds == "vulncheck:install" {
			foundDep = true
			break
		} else if dm, ok := d.(map[string]interface{}); ok {
			if task, ok := dm["task"].(string); ok && task == "vulncheck:install" {
				foundDep = true
				break
			}
		}
	}
	assert.True(t, foundDep, "vulncheck must depend on vulncheck:install")
	
	// Verify ci:vulncheck task
	ciVulncheck, exists := tf.Tasks["ci:vulncheck"]
	require.True(t, exists, "ci:vulncheck task must exist")
	
	foundCIDep := false
	for _, d := range ciVulncheck.Deps {
		if ds, ok := d.(string); ok && ds == "vulncheck:install" {
			foundCIDep = true
			break
		} else if dm, ok := d.(map[string]interface{}); ok {
			if task, ok := dm["task"].(string); ok && task == "vulncheck:install" {
				foundCIDep = true
				break
			}
		}
	}
	assert.True(t, foundCIDep, "ci:vulncheck must depend on vulncheck:install")
}
