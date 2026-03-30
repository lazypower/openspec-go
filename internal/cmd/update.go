package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chuck/openspec-go/internal/editor"
	tmpl "github.com/chuck/openspec-go/internal/template"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [path]",
		Short: "Refresh managed instruction files",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := ""
			if len(args) > 0 {
				path = args[0]
			}

			ospPath, err := findOpenSpecPath(path)
			if err != nil {
				return err
			}
			projectPath := filepath.Dir(ospPath)

			// Update AGENTS.md
			agentsContent, err := tmpl.Render("agents.md", nil)
			if err != nil {
				return err
			}
			agentsPath := filepath.Join(ospPath, "AGENTS.md")
			if err := os.WriteFile(agentsPath, []byte(agentsContent), 0o644); err != nil {
				return err
			}
			fmt.Println("Updated", agentsPath)

			// Update configured editors
			for _, name := range editor.AllNames() {
				ed, _ := editor.Get(name)
				if ed.IsConfigured(projectPath) {
					if err := ed.UpdateExisting(projectPath, ospPath); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to update %s: %v\n", name, err)
						continue
					}
					fmt.Printf("Updated %s configuration\n", name)
				}
			}

			return nil
		},
	}
}
