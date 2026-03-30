package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chuck/openspec-go/internal/editor"
	tmpl "github.com/chuck/openspec-go/internal/template"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newInitCmd() *cobra.Command {
	var tools string

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Initialize OpenSpec directory structure",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			ospPath := filepath.Join(abs, "openspec")
			if _, err := os.Stat(ospPath); err == nil {
				return fmt.Errorf("openspec directory already exists at %s", ospPath)
			}

			var selectedTools []string
			if tools != "" {
				selectedTools = resolveTools(tools)
			} else if term.IsTerminal(int(os.Stdout.Fd())) {
				selectedTools, err = interactiveToolSelect()
				if err != nil {
					return err
				}
			} else {
				fmt.Fprintln(cmd.ErrOrStderr(), "Non-TTY detected. Use --tools flag to specify editors. Proceeding with no editor configuration.")
			}

			return initProject(abs, ospPath, selectedTools)
		},
	}

	cmd.Flags().StringVar(&tools, "tools", "", "Comma-separated editors: claude-code,opencode,codex,goose,all,none")
	return cmd
}

func resolveTools(tools string) []string {
	tools = strings.TrimSpace(strings.ToLower(tools))
	if tools == "all" {
		return editor.AllNames()
	}
	if tools == "none" {
		return nil
	}
	var result []string
	for _, t := range strings.Split(tools, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}

func interactiveToolSelect() ([]string, error) {
	var selected []string
	options := make([]huh.Option[string], len(editor.AllNames()))
	for i, name := range editor.AllNames() {
		options[i] = huh.NewOption(name, name)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select AI editors to configure").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}
	return selected, nil
}

func initProject(projectPath, ospPath string, tools []string) error {
	// Create directory structure
	dirs := []string{
		ospPath,
		filepath.Join(ospPath, "specs"),
		filepath.Join(ospPath, "changes"),
		filepath.Join(ospPath, "changes", "archive"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	// Write project.md
	projectContent, err := tmpl.Render("project.md", nil)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(ospPath, "project.md"), []byte(projectContent), 0o644); err != nil {
		return err
	}

	// Write AGENTS.md
	agentsContent, err := tmpl.Render("agents.md", nil)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(ospPath, "AGENTS.md"), []byte(agentsContent), 0o644); err != nil {
		return err
	}

	// Configure editors
	for _, toolName := range tools {
		ed, err := editor.Get(toolName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: unknown editor %q, skipping\n", toolName)
			continue
		}
		if err := ed.Configure(projectPath, ospPath); err != nil {
			return fmt.Errorf("configuring %s: %w", toolName, err)
		}
	}

	fmt.Printf("Initialized OpenSpec at %s\n", ospPath)
	return nil
}
