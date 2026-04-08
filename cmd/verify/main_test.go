package main

import (
	"testing"
)

func TestCICmdHasWithDockerFlag(t *testing.T) {
	flag := ciCmd.Flags().Lookup("with-docker")
	if flag == nil {
		t.Fatal("ciCmd should have a --with-docker flag")
	}
	if flag.DefValue != "false" {
		t.Errorf("expected default false, got %s", flag.DefValue)
	}
}
