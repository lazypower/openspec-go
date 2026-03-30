package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/chuck/openspec-go/internal/output"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var showSpecs, showChanges bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List changes or specs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ospPath, err := findOpenSpecPath("")
			if err != nil {
				return err
			}

			if showSpecs {
				return listSpecs(ospPath)
			}
			return listChanges(ospPath)
		},
	}

	cmd.Flags().BoolVar(&showSpecs, "specs", false, "List specifications")
	cmd.Flags().BoolVar(&showChanges, "changes", false, "List changes (default)")
	return cmd
}

func listChanges(ospPath string) error {
	changes, err := loadChanges(ospPath)
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		fmt.Println("No active changes found.")
		return nil
	}

	for _, ch := range changes {
		progress := output.FormatProgress(ch.Tasks.Completed, ch.Tasks.Total)
		fmt.Printf("  %s  %s\n", output.Yellow(ch.ID), output.Dim(progress))
	}
	return nil
}

func listSpecs(ospPath string) error {
	specs, err := loadSpecs(ospPath)
	if err != nil {
		return err
	}

	if len(specs) == 0 {
		fmt.Println("No specifications found.")
		return nil
	}

	// Sort by requirement count descending
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
		fmt.Fprintf(os.Stdout, "  %s  %s\n", output.Cyan(e.name), output.Dim(fmt.Sprintf("%d requirements", e.reqCount)))
	}
	return nil
}
