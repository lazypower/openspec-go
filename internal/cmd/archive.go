package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	archivePkg "github.com/chuck/openspec-go/internal/archive"
	"github.com/chuck/openspec-go/internal/output"
	"github.com/chuck/openspec-go/internal/parser"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newArchiveCmd() *cobra.Command {
	var (
		yes        bool
		skipSpecs  bool
		noValidate bool
	)

	cmd := &cobra.Command{
		Use:   "archive <change-id>",
		Short: "Archive a completed change and merge deltas",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ospPath, err := findOpenSpecPath("")
			if err != nil {
				return err
			}

			changeID := args[0]
			changePath := filepath.Join(ospPath, "changes", changeID)
			if _, err := os.Stat(changePath); os.IsNotExist(err) {
				return fmt.Errorf("change %q not found", changeID)
			}

			// Check incomplete tasks
			tasksData, err := os.ReadFile(filepath.Join(changePath, "tasks.md"))
			if err == nil {
				tasks := parser.ParseTaskProgress(string(tasksData))
				incomplete := tasks.Total - tasks.Completed
				if incomplete > 0 {
					fmt.Fprintf(os.Stderr, "%s %d incomplete task(s)\n", output.WarnStyle.Render("WARNING:"), incomplete)
					if !yes {
						if !confirmAction("Proceed with archiving despite incomplete tasks?") {
							return fmt.Errorf("aborted")
						}
					}
				}
			}

			// Confirmation
			if !yes {
				if !confirmAction(fmt.Sprintf("Archive change %q?", changeID)) {
					return fmt.Errorf("aborted")
				}
			}

			err = archivePkg.Archive(archivePkg.ArchiveOptions{
				OpenSpecPath: ospPath,
				ChangeID:     changeID,
				SkipSpecs:    skipSpecs,
				SkipValidate: noValidate,
				Yes:          yes,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Archived %s\n", output.Green(changeID))
			return nil
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Skip confirmation prompts")
	cmd.Flags().BoolVar(&skipSpecs, "skip-specs", false, "Archive without merging deltas into specs")
	cmd.Flags().BoolVar(&noValidate, "no-validate", false, "Skip validation before archiving")
	return cmd
}

func confirmAction(prompt string) bool {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return false
	}

	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(prompt).
				Value(&confirmed),
		),
	)
	if err := form.Run(); err != nil {
		return false
	}
	return confirmed
}
