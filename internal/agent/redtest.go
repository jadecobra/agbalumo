package agent

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type TestResult struct {
	Success           bool
	CompilationFailed bool
	Failures          []TestFailure
}

type TestFailure struct {
	TestName string
	Output   string
}

type testEvent struct {
	Time       string  `json:"Time"`
	Action     string  `json:"Action"`
	Package    string  `json:"Package"`
	Test       string  `json:"Test"`
	Output     string  `json:"Output"`
	Elapsed    float64 `json:"Elapsed"`
	ImportPath string  `json:"ImportPath"`
}

func ParseTestJSON(r io.Reader) (*TestResult, error) {
	scanner := bufio.NewScanner(r)
	result := &TestResult{
		Success:  true, // assume success until proven otherwise
		Failures: []TestFailure{},
	}

	testOutputs := make(map[string]*strings.Builder)
	buildOutputs := &strings.Builder{}

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 || line[0] != '{' {
			continue
		}

		var event testEvent
		if err := json.Unmarshal(line, &event); err != nil {
			continue
		}

		switch event.Action {
		case "build-output":
			buildOutputs.WriteString(event.Output)
		case "build-fail":
			result.Success = false
			result.CompilationFailed = true
		case "output":
			if strings.Contains(event.Output, "[setup failed]") {
				result.Success = false
				result.CompilationFailed = true
			}
		}

		if event.Test != "" {
			switch event.Action {
			case "output":
				b, ok := testOutputs[event.Test]
				if !ok {
					b = &strings.Builder{}
					testOutputs[event.Test] = b
				}
				b.WriteString(event.Output)
			case "fail":
				result.Success = false
				result.Failures = append(result.Failures, TestFailure{
					TestName: event.Test,
					Output:   testOutputs[event.Test].String(),
				})
			}
		} else if event.Package != "" && event.Action == "fail" {
			result.Success = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if result.CompilationFailed && buildOutputs.Len() > 0 {
		result.Failures = append(result.Failures, TestFailure{
			TestName: "build",
			Output:   buildOutputs.String(),
		})
	}

	return result, nil
}
