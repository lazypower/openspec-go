package editor

import (
	"os"
	"path/filepath"

	tmpl "github.com/chuck/openspec-go/internal/template"
)

// Codex configures Codex with OpenSpec prompts.
type Codex struct{}

func (c *Codex) Name() string { return "codex" }

func (c *Codex) Configure(projectPath, openspecPath string) error {
	return c.writePrompts(projectPath)
}

func (c *Codex) UpdateExisting(projectPath, openspecPath string) error {
	return c.writePrompts(projectPath)
}

func (c *Codex) IsConfigured(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, ".codex", "prompts", "openspec"))
	return err == nil
}

func (c *Codex) writePrompts(projectPath string) error {
	promptDir := filepath.Join(projectPath, ".codex", "prompts", "openspec")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		return err
	}

	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		content, err := tmpl.Raw(cmd + ".md")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(promptDir, cmd+".md"), []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}
