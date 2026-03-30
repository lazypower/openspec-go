package editor

import (
	"os"
	"path/filepath"
	"strings"

	tmpl "github.com/chuck/openspec-go/internal/template"
)

const (
	openspecStartMarker = "<!-- OPENSPEC:START -->"
	openspecEndMarker   = "<!-- OPENSPEC:END -->"
)

// ClaudeCode configures Claude Code slash commands and CLAUDE.md.
type ClaudeCode struct{}

func (c *ClaudeCode) Name() string { return "claude-code" }

func (c *ClaudeCode) Configure(projectPath, openspecPath string) error {
	if err := c.writeSlashCommands(projectPath); err != nil {
		return err
	}
	return c.injectClaudeMd(projectPath)
}

func (c *ClaudeCode) UpdateExisting(projectPath, openspecPath string) error {
	if err := c.writeSlashCommands(projectPath); err != nil {
		return err
	}
	return c.injectClaudeMd(projectPath)
}

func (c *ClaudeCode) IsConfigured(projectPath string) bool {
	cmdDir := filepath.Join(projectPath, ".claude", "commands", "openspec")
	_, err := os.Stat(cmdDir)
	return err == nil
}

func (c *ClaudeCode) writeSlashCommands(projectPath string) error {
	cmdDir := filepath.Join(projectPath, ".claude", "commands", "openspec")
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

func (c *ClaudeCode) injectClaudeMd(projectPath string) error {
	claudeMdPath := filepath.Join(projectPath, "CLAUDE.md")
	block, err := tmpl.Raw("claude.md")
	if err != nil {
		return err
	}

	existing, err := os.ReadFile(claudeMdPath)
	if err != nil {
		// File doesn't exist, create with just the block
		return os.WriteFile(claudeMdPath, []byte(block), 0o644)
	}

	content := string(existing)
	content = replaceMarkerBlock(content, block)
	return os.WriteFile(claudeMdPath, []byte(content), 0o644)
}

// replaceMarkerBlock replaces the OPENSPEC:START..END block, or appends if not found.
func replaceMarkerBlock(content, block string) string {
	startIdx := strings.Index(content, openspecStartMarker)
	endIdx := strings.Index(content, openspecEndMarker)

	if startIdx >= 0 && endIdx >= 0 {
		return content[:startIdx] + block + content[endIdx+len(openspecEndMarker):]
	}

	// Append
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return content + "\n" + block
}
