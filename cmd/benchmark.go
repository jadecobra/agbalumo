package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"
	"time"

	"github.com/jadecobra/agbalumo/internal/domain"
	"github.com/jadecobra/agbalumo/internal/repository/sqlite"
	"github.com/spf13/cobra"
)

var warmup bool

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Run read-heavy queries to benchmark SQLite performance",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := ResolveSeedConfig(args, os.Getenv)

		repo, err := sqlite.NewSQLiteRepository(dbPath)
		if err != nil {
			slog.Error(domain.MsgFailedToOpenDB, "path", dbPath, "error", err)
			os.Exit(1)
		}

		ctx := context.Background()

		fmt.Println("Running benchmark scenarios on SQLite data model...")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "Scenario\tExecution Time (ms)\tResults Count")
		_, _ = fmt.Fprintln(w, "--------\t-------------------\t-------------")

		scenarios := []struct {
			name       string
			filterType string
			queryText  string
			sortField  string
			sortOrder  string
			limit      int
			offset     int
		}{
			{"Page 1 (No Filters)", "", "", "", "", 20, 0},
			{"Page 500 (Deep Pagination)", "", "", "", "", 20, 10000},
			{"Category Filter ('Business')", "Business", "", "", "", 20, 0},
		}

		for _, s := range scenarios {
			if warmup {
				for i := 0; i < 5; i++ {
					_, _, _ = repo.FindAll(ctx, s.filterType, s.queryText, "", s.sortField, s.sortOrder, false, s.limit, s.offset)
				}
			}

			start := time.Now()
			listings, _, err := repo.FindAll(ctx, s.filterType, s.queryText, "", s.sortField, s.sortOrder, false, s.limit, s.offset)
			duration := time.Since(start)

			if err != nil {
				slog.Error("Query failed", "scenario", s.name, "error", err)
				continue
			}

			_, _ = fmt.Fprintf(w, "%s\t%.2f\t%d\n", s.name, float64(duration.Microseconds())/1000.0, len(listings))
		}

		_ = w.Flush()
	},
}

func init() {
	benchmarkCmd.Flags().BoolVar(&warmup, "warmup", false, "Execute queries 5 times in a loop before measuring time")
	rootCmd.AddCommand(benchmarkCmd)
}
