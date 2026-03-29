package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Persona struct {
	Name           string   `yaml:"name"`
	Role           string   `yaml:"role"`
	PreferredModel string   `yaml:"preferred_model"`
	Description    string   `yaml:"description"`
	Tools          []string `yaml:"tools"`
	Instructions   string   `yaml:"instructions"`
}

type Config struct {
	Name        string   `yaml:"name"`
	GlobalRules []string `yaml:"global_rules"`
	Personas    []struct {
		Name string `yaml:"name"`
		File string `yaml:"file"`
	} `yaml:"personas"`
}

func main() {
	configPath := ".agents/config.yaml"
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- Verifying Squad: %s ---\n", config.Name)

	hasErrors := false
	var yamlPersonas []string

	for _, pReg := range config.Personas {
		yamlPersonas = append(yamlPersonas, pReg.Name)
		pPath := filepath.Join(".agents", pReg.File)
		// ... existing persona validation ...
		pData, err := os.ReadFile(pPath)
		if err != nil {
			fmt.Printf("[FAIL] Could not read persona file %s: %v\n", pPath, err)
			hasErrors = true
			continue
		}

		var p Persona
		err = yaml.Unmarshal(pData, &p)
		if err != nil {
			fmt.Printf("[FAIL] Error parsing YAML for %s: %v\n", pPath, err)
			hasErrors = true
			continue
		}

		// Validation: Mandatory Fields
		if p.Name == "" || p.Role == "" || p.Instructions == "" {
			fmt.Printf("[FAIL] %s: Missing mandatory fields (name, role, or instructions)\n", pPath)
			hasErrors = true
		}

		// Validation: Forbidden Patterns
		if strings.Contains(strings.ToLower(p.Name), "sdet") && strings.Contains(strings.ToLower(p.Instructions), "production code") {
			if !strings.Contains(strings.ToLower(p.Instructions), "never write production code") {
				fmt.Printf("[FAIL] %s: SDET instructions might allow production code writing\n", pPath)
				hasErrors = true
			}
		}

		// Validation: Brand Constants (UIUX)
		if strings.Contains(strings.ToLower(p.Name), "uiux") {
			if !strings.Contains(p.Instructions, "#FF5E0E") || !strings.Contains(p.Instructions, "#2D5A27") {
				fmt.Printf("[FAIL] %s: UIUXDesigner missing mandatory brand colors (#FF5E0E, #2D5A27)\n", pPath)
				hasErrors = true
			}
		}
		
		// Validation: Chaos Monkey Safety Gates
		if strings.Contains(strings.ToLower(p.Name), "chaos") {
			if p.PreferredModel != "GEMINI_3_1_PRO" {
				fmt.Printf("[FAIL] %s: ChaosMonkey MUST use GEMINI_3_1_PRO for deep reasoning\n", pPath)
				hasErrors = true
			}
			if !strings.Contains(strings.ToLower(p.Instructions), "anti-cheat") && !strings.Contains(strings.ToLower(p.Instructions), "bypass") {
				fmt.Printf("[FAIL] %s: ChaosMonkey missing anti-cheat/bypass verification instructions\n", pPath)
				hasErrors = true
			}
		}

		if !hasErrors {
			fmt.Printf("[PASS] %s validated\n", p.Name)
		}
	}

	// Drift Check: Compare with CODING_STANDARDS.md
	fmt.Println("\n--- Checking for Drift (CODING_STANDARDS.md) ---")
	mdData, err := os.ReadFile("docs/CODING_STANDARDS.md")
	if err != nil {
		fmt.Printf("[FAIL] Could not read docs/CODING_STANDARDS.md: %v\n", err)
		hasErrors = true
	} else {
		mdContent := string(mdData)
		var mdPersonas []string
		
		// Simple regex-like search for "*   **PersonaName**" in Section 5
		lines := strings.Split(mdContent, "\n")
		inSection5 := false
		for _, line := range lines {
			if strings.Contains(line, "## 5. Agent Protocol") {
				inSection5 = true
				continue
			}
			if inSection5 && strings.HasPrefix(line, "##") {
				break
			}
			if inSection5 && strings.HasPrefix(line, "*   **") {
				parts := strings.Split(line, "**")
				if len(parts) >= 2 {
					mdPersonas = append(mdPersonas, parts[1])
				}
			}
		}

		// Compare
		if len(yamlPersonas) != len(mdPersonas) {
			fmt.Printf("[FAIL] Persona count mismatch: YAML=%d, MD=%d\n", len(yamlPersonas), len(mdPersonas))
			hasErrors = true
		}

		for _, yp := range yamlPersonas {
			found := false
			for _, mp := range mdPersonas {
				if yp == mp {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("[FAIL] Persona '%s' found in config but missing in CODING_STANDARDS.md\n", yp)
				hasErrors = true
			}
		}
	}

	if hasErrors {
		fmt.Println("\n--- Squad Verification FAILED ---")
		os.Exit(1)
	} else {
		fmt.Println("\n--- Squad Verification SUCCESS ---")
	}
}
