package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveServerConfig(t *testing.T) {
	tests := []struct {
		name           string
		env            string
		port           string
		filesExist     bool
		expectedAddr   string
		expectedTLS    bool
		expectedConfig ServerConfig
	}{
		{
			name:         "Production Default",
			env:          "production",
			port:         "",
			filesExist:   false,
			expectedAddr: ":8080",
			expectedTLS:  false,
		},
		{
			name:         "Production With Port",
			env:          "production",
			port:         "9090",
			filesExist:   false,
			expectedAddr: ":9090",
			expectedTLS:  false,
		},
		{
			name:         "Dev Default (No Certs)",
			env:          "development",
			port:         "",
			filesExist:   false,
			expectedAddr: ":8080",
			expectedTLS:  false,
		},
		{
			name:         "Dev With Port (No Certs)",
			env:          "development",
			port:         "3000",
			filesExist:   false,
			expectedAddr: ":3000",
			expectedTLS:  false,
		},
		{
			name:         "Dev With Certs",
			env:          "development",
			port:         "3000",
			filesExist:   true,
			expectedAddr: ":8443", // Should override port and use 8443
			expectedTLS:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFileExists := func(path string) bool {
				return tt.filesExist
			}

			config := ResolveServerConfig(tt.env, tt.port, mockFileExists)

			assert.Equal(t, tt.expectedAddr, config.Addr)
			assert.Equal(t, tt.expectedTLS, config.TLS)
			if tt.expectedTLS {
				assert.Equal(t, "certs/cert.pem", config.CertFile)
				assert.Equal(t, "certs/key.pem", config.KeyFile)
			}
		})
	}
}
