package cmd

import (
	"testing"
)

func TestResolveSeedConfig(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		env      map[string]string
		expected string
	}{
		{
			name:     "Default",
			args:     []string{},
			env:      map[string]string{},
			expected: "agbalumo.db",
		},
		{
			name:     "WithArgument",
			args:     []string{"custom.db"},
			env:      map[string]string{},
			expected: "custom.db",
		},
		{
			name:     "WithEnvVar",
			args:     []string{},
			env:      map[string]string{"DATABASE_URL": "env.db"},
			expected: "env.db",
		},
		{
			name:     "ArgOverridesEnv",
			args:     []string{"arg.db"},
			env:      map[string]string{"DATABASE_URL": "env.db"},
			expected: "arg.db",
		},
		{
			name:     "ExplicitDefaultArg",
			args:     []string{"agbalumo.db"},
			env:      map[string]string{"DATABASE_URL": "env.db"},
			expected: "env.db", // Logic: if dbPath == "agbalumo.db", check env.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getEnv := func(key string) string {
				return tt.env[key]
			}
			result := ResolveSeedConfig(tt.args, getEnv)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
