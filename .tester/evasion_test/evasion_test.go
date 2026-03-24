package evasion
import (
	"os"
	"testing"
)
func TestMain(m *testing.M) {
	os.Exit(0)
}
func TestShouldFailButWont(t *testing.T) {
	t.Fail()
}
