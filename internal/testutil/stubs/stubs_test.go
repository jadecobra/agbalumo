package stubs

import (
	"errors"
	"testing"
)

func TestNoopLogger(t *testing.T) {
	n := &NoopLogger{}
	n.Log("test")
	n.Logf("test %s", "format")
}

func TestFailingLogger(t *testing.T) {
	f := &FailingLogger{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got nil")
		}
	}()
	f.Log("test")
}

func TestFailingLogger_Logf(t *testing.T) {
	f := &FailingLogger{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got nil")
		}
	}()
	f.Logf("test %s", "format")
}

func TestStubResolver(t *testing.T) {
	// Case 1: Custom ResolveFunc
	s := &StubResolver{
		ResolveFunc: func(id string) (interface{}, error) {
			return id, nil
		},
	}
	res, err := s.Resolve("id1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "id1" {
		t.Errorf("expected id1, got %v", res)
	}

	// Case 2: Nil ResolveFunc
	s2 := &StubResolver{}
	_, err = s2.Resolve("id2")
	if !errors.Is(err, ErrDefault) {
		t.Errorf("expected ErrDefault, got %v", err)
	}
}
