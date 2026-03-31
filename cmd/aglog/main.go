package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jadecobra/agbalumo/internal/history"
)

func main() {
	if err := run(os.Args, os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout io.Writer) error {
	fs := flag.NewFlagSet("aglog", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // Silencing standard output for flags

	fFeature := fs.String("feature", "", "Feature name")
	fArch := fs.String("arch", "", "Systems Architect")
	fPO := fs.String("po", "", "Product Owner")
	fSDET := fs.String("sdet", "", "SDET")
	fBE := fs.String("be", "", "Backend Engineer")
	fSummary := fs.String("summary", "", "Decision summary")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	var decision history.SquadDecision

	// If flags are provided, use them. Otherwise, try reading from stdin.
	if *fFeature != "" {
		decision = history.SquadDecision{
			FeatureName:      *fFeature,
			SystemsArchitect: *fArch,
			ProductOwner:     *fPO,
			SDET:             *fSDET,
			BackendEngineer:  *fBE,
			DecisionSummary:  *fSummary,
		}
	} else if stdin != nil {
		data, err := io.ReadAll(stdin)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}

		if len(data) > 0 {
			if err := json.Unmarshal(data, &decision); err != nil {
				return fmt.Errorf("failed to unmarshal JSON: %w", err)
			}
		}
	}

	if decision.FeatureName == "" {
		return fmt.Errorf("FeatureName is required via flags or JSON input")
	}

	path, err := history.Store(decision)
	if err != nil {
		return fmt.Errorf("failed to store decision: %w", err)
	}

	_, err = fmt.Fprintln(stdout, path)
	return err
}
