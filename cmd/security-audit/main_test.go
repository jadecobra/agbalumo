package main

import (
	"testing"
)



func TestSecurityAuditCoverage(t *testing.T) {
	runner := &MockRunner{}
	checkGoVet(".", runner)
	checkFlyConfig("nothing to see here", false)
	
	if containsSensitive("hello SECRET=value", "SECRET") == false {
		t.Error("expected true")
	}
	if containsSensitive("safe config", "SECRET") == true {
		t.Error("expected false")
	}
}
