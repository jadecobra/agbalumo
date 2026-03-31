package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jadecobra/agbalumo/internal/history"
	"github.com/spf13/cobra"
)

func main() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	var fFeature, fArch, fPO, fSDET, fBE, fSummary string

	cmd := &cobra.Command{
		Use:          "aglog",
		Short:        "Capture squad decisions",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.InOrStdin(), cmd.OutOrStdout(), fFeature, fArch, fPO, fSDET, fBE, fSummary)
		},
	}

	cmd.Flags().StringVar(&fFeature, "feature", "", "Feature name")
	cmd.Flags().StringVar(&fArch, "arch", "", "Systems Architect")
	cmd.Flags().StringVar(&fPO, "po", "", "Product Owner")
	cmd.Flags().StringVar(&fSDET, "sdet", "", "SDET")
	cmd.Flags().StringVar(&fBE, "be", "", "Backend Engineer")
	cmd.Flags().StringVar(&fSummary, "summary", "", "Decision summary")

	return cmd
}

func run(stdin io.Reader, stdout io.Writer, feature, arch, po, sdet, be, summary string) error {
	var decision history.SquadDecision

	if feature != "" {
		decision = history.SquadDecision{
			FeatureName:      feature,
			SystemsArchitect: arch,
			ProductOwner:     po,
			SDET:             sdet,
			BackendEngineer:  be,
			DecisionSummary:  summary,
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
