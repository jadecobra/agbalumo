package agent_test

func replaceValue(s, old, new string) string {
	res := ""
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			res = s[:i] + new + s[i+len(old):]
			break
		}
	}
	return res
}

func contains(s, substr string) bool {
	return ind(s, substr) != -1
}

func ind(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
