package maintenance

import (
	"slices"
	"testing"
)

func TestBuildLinterCommand_Precommit(t *testing.T) {
	opts := ChiefCriticOptions{Full: false, NewFromRev: "HEAD"}
	cmd := buildLinterCommand(opts)
	assertContains(t, cmd, "--new-from-rev")
	assertContains(t, cmd, "HEAD")
	assertContains(t, cmd, "--whole-files")
}

func TestBuildLinterCommand_DefaultRev(t *testing.T) {
	opts := ChiefCriticOptions{Full: false}
	cmd := buildLinterCommand(opts)
	assertContains(t, cmd, "--new-from-rev")
	assertContains(t, cmd, "HEAD~1")
	assertContains(t, cmd, "--whole-files")
}

func TestBuildLinterCommand_FullAudit(t *testing.T) {
	opts := ChiefCriticOptions{Full: true}
	cmd := buildLinterCommand(opts)
	assertNotContains(t, cmd, "--whole-files")
	assertNotContains(t, cmd, "--new-from-rev")
}

func assertContains(t *testing.T, cmd []string, arg string) {
	t.Helper()
	if !slices.Contains(cmd, arg) {
		t.Errorf("buildLinterCommand missing argument %q; got: %v", arg, cmd)
	}
}

func assertNotContains(t *testing.T, cmd []string, arg string) {
	t.Helper()
	if slices.Contains(cmd, arg) {
		t.Errorf("buildLinterCommand unexpectedly contained %q; got: %v", arg, cmd)
	}
}
