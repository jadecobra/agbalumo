package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	content, err := os.ReadFile(".github/workflows/ci.yml")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if strings.Contains(string(content), "trivy-version:") {
		fmt.Println("❌ REPRODUCED: Found invalid 'trivy-version' input in ci.yml. Production CI will fail.")
		os.Exit(1)
	}

	fmt.Println("✅ No invalid trivy-version detected.")
}
