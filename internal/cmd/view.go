package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chuck/openspec-go/internal/output"
	"github.com/spf13/cobra"
)

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Display terminal dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			ospPath, err := findOpenSpecPath("")
			if err != nil {
				return err
			}
			return renderDashboard(ospPath)
		},
	}
}

func renderDashboard(ospPath string) error {
	changes, _ := loadChanges(ospPath)
	specs, _ := loadSpecs(ospPath)
	archived := loadArchivedChanges(ospPath)

	// Count totals
	totalReqs := 0
	for _, spec := range specs {
		totalReqs += len(spec.Requirements)
	}
	totalTasks := 0
	completedTasks := 0
	for _, ch := range changes {
		totalTasks += ch.Tasks.Total
		completedTasks += ch.Tasks.Completed
	}

	// Header
	fmt.Println(output.Bold("═══════════════════════════════════════"))
	fmt.Println(output.Bold("  OpenSpec Dashboard"))
	fmt.Println(output.Bold("═══════════════════════════════════════"))
	fmt.Println()

	// Summary
	fmt.Printf("  %s  %s    %s  %s\n",
		output.Cyan("Specs:"), output.Cyan(fmt.Sprintf("%d", len(specs))),
		output.Cyan("Requirements:"), output.Cyan(fmt.Sprintf("%d", totalReqs)),
	)
	fmt.Printf("  %s  %s    %s  %s\n",
		output.Yellow("Active:"), output.Yellow(fmt.Sprintf("%d", len(changes))),
		output.Green("Archived:"), output.Green(fmt.Sprintf("%d", len(archived))),
	)
	if totalTasks > 0 {
		pct := completedTasks * 100 / totalTasks
		fmt.Printf("  %s  %s\n",
			output.Magenta("Tasks:"),
			output.Magenta(fmt.Sprintf("%d/%d (%d%%)", completedTasks, totalTasks, pct)),
		)
	}
	fmt.Println()

	// Active changes
	if len(changes) > 0 {
		fmt.Println(output.Bold("─── Active Changes ───"))
		for _, ch := range changes {
			pct := ch.Tasks.Percent()
			bar := output.Green(output.ProgressBar(pct, 20))
			fmt.Printf("  %s  %s %s\n",
				output.Yellow(ch.ID),
				bar,
				output.Dim(output.FormatProgress(ch.Tasks.Completed, ch.Tasks.Total)),
			)
		}
		fmt.Println()
	}

	// Archived changes
	if len(archived) > 0 {
		fmt.Println(output.Bold("─── Completed Changes ───"))
		for _, id := range archived {
			fmt.Printf("  %s %s\n", output.Green("✓"), id)
		}
		fmt.Println()
	}

	// Specs
	if len(specs) > 0 {
		fmt.Println(output.Bold("─── Specifications ───"))
		type specEntry struct {
			name     string
			reqCount int
		}
		var entries []specEntry
		for name, spec := range specs {
			entries = append(entries, specEntry{name, len(spec.Requirements)})
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].reqCount > entries[j].reqCount
		})
		for _, e := range entries {
			fmt.Printf("  %s  %s\n",
				output.Cyan(e.name),
				output.Dim(fmt.Sprintf("%d requirements", e.reqCount)),
			)
		}
		fmt.Println()
	}

	if len(changes) == 0 && len(specs) == 0 && len(archived) == 0 {
		fmt.Println(output.Dim("  No specs or changes found."))
		fmt.Println()
	}

	fmt.Println(strings.Repeat("─", 39))
	return nil
}
