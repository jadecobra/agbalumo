package stubs

import (
	"errors"
	"fmt"
)

// ErrDefault is a standard error used for testing faliure paths.
var ErrDefault = errors.New("test: default stub error")

// NoopLogger is a logger that does nothing.
type NoopLogger struct{}

func (n *NoopLogger) Log(v ...interface{})                 {}
func (n *NoopLogger) Logf(format string, v ...interface{}) {}

// FailingLogger is a logger that panics on log.
type FailingLogger struct{}

func (f *FailingLogger) Log(v ...interface{}) { panic(errors.New("failing: log not allowed")) }
func (f *FailingLogger) Logf(format string, v ...interface{}) {
	panic(errors.New("failing: log not allowed"))
}

// StubResolver is a generic stub for any resolver pattern.
type StubResolver struct {
	ResolveFunc func(id string) (interface{}, error)
}

func (s *StubResolver) Resolve(id string) (interface{}, error) {
	if s.ResolveFunc != nil {
		return s.ResolveFunc(id)
	}
	return nil, fmt.Errorf("%w: no resolve function provided", ErrDefault)
}
