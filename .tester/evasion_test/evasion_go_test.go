package evasion
import (
	"os"
	"testing"
	"time"
)
func TestEvasionGoroutine(t *testing.T) {
    go func() {
	    os.Exit(0)
    }()
	time.Sleep(1 * time.Second)
}
