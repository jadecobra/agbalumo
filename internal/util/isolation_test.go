package util

import (
	"os"
	"runtime"
	"testing"
)

func TestContainerIsolation(t *testing.T) {
	// This test is designed to fail when running on the host but pass when running in a Linux container.
	// It is triggered only when CONTAINER_ISOLATION_REQUIRED is set to true.

	if os.Getenv("CONTAINER_ISOLATION_REQUIRED") != "true" {
		t.Skip("Container isolation test skipped (CONTAINER_ISOLATION_REQUIRED != true)")
	}

	t.Logf("Running isolation check on %s/%s", runtime.GOOS, runtime.GOARCH)

	// If we are on MacOS, we are definitely NOT in the ephemeral Linux container.
	if runtime.GOOS == "darwin" {
		t.Errorf("FAIL: Running on host MacOS instead of isolated Linux container.")
	}

	// Additional check: Host environment variables should not leak unless explicitly passed.
	if os.Getenv("HOST_SECRET_MARKER") == "exposed" {
		t.Errorf("FAIL: Host environment variable 'HOST_SECRET_MARKER' leaked into container.")
	}
}
