package editor

import (
	"os"
	"path/filepath"

	tmpl "github.com/chuck/openspec-go/internal/template"
)

// OpenCode configures OpenCode with OpenSpec commands.
type OpenCode struct{}

func (o *OpenCode) Name() string { return "opencode" }

func (o *OpenCode) Configure(projectPath, openspecPath string) error {
	return o.writeConfig(projectPath)
}

func (o *OpenCode) UpdateExisting(projectPath, openspecPath string) error {
	return o.writeConfig(projectPath)
}

func (o *OpenCode) IsConfigured(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, ".opencode", "commands", "openspec"))
	return err == nil
}

func (o *OpenCode) writeConfig(projectPath string) error {
	cmdDir := filepath.Join(projectPath, ".opencode", "commands", "openspec")
	if err := os.MkdirAll(cmdDir, 0o755); err != nil {
		return err
	}

	commands := []string{"proposal", "apply", "archive"}
	for _, cmd := range commands {
		content, err := tmpl.Raw(cmd + ".md")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(cmdDir, cmd+".md"), []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}
