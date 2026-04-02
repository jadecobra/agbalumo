package main


// MockCommandRunner
type MockRunner struct {
	CapturedDir  string
	CapturedName string
	CapturedArgs []string
	Output       string
	Err          error
}

func (m *MockRunner) Run(dir string, name string, args ...string) (string, error) {
	m.CapturedDir = dir
	m.CapturedName = name
	m.CapturedArgs = args
	return m.Output, m.Err
}

type RunnerResponse struct {
	Output string
	Err    error
}

type SmartMockRunner struct {
	Responses map[string]RunnerResponse
}

func (m *SmartMockRunner) Run(dir string, name string, args ...string) (string, error) {
	if resp, ok := m.Responses[name]; ok {
		return resp.Output, resp.Err
	}
	return "", nil
}
