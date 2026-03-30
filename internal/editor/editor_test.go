package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- Claude Code Tests ---

func TestClaudeCode_Configure(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	c := &ClaudeCode{}
	if err := c.Configure(dir, osp); err != nil {
		t.Fatal(err)
	}

	// Verify slash commands created
	for _, cmd := range []string{"proposal", "apply", "archive"} {
		path := filepath.Join(dir, ".claude", "commands", "openspec", cmd+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("slash command not created: %s", path)
		}
	}
}

func TestClaudeCode_ClaudeMdInjection(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	c := &ClaudeCode{}
	if err := c.Configure(dir, osp); err != nil {
		t.Fatal(err)
	}

	claudeMdPath := filepath.Join(dir, "CLAUDE.md")
	data, err := os.ReadFile(claudeMdPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, openspecStartMarker) {
		t.Error("CLAUDE.md missing start marker")
	}
	if !strings.Contains(content, openspecEndMarker) {
		t.Error("CLAUDE.md missing end marker")
	}
}

func TestClaudeCode_Update(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	c := &ClaudeCode{}
	c.Configure(dir, osp)

	// Update should work
	if err := c.UpdateExisting(dir, osp); err != nil {
		t.Fatal(err)
	}

	// Slash commands should still exist
	for _, cmd := range []string{"proposal", "apply", "archive"} {
		path := filepath.Join(dir, ".claude", "commands", "openspec", cmd+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("slash command missing after update: %s", cmd)
		}
	}
}

func TestClaudeCode_MarkerBlockReplace(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	// Write existing CLAUDE.md with old block
	claudeMdPath := filepath.Join(dir, "CLAUDE.md")
	existing := "# My Project\n\nSome content.\n\n<!-- OPENSPEC:START -->\nOLD CONTENT\n<!-- OPENSPEC:END -->\n\nMore content.\n"
	os.WriteFile(claudeMdPath, []byte(existing), 0o644)

	c := &ClaudeCode{}
	if err := c.injectClaudeMd(dir); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(claudeMdPath)
	content := string(data)
	if strings.Contains(content, "OLD CONTENT") {
		t.Error("old content still present")
	}
	if !strings.Contains(content, "# My Project") {
		t.Error("surrounding content was lost")
	}
	if !strings.Contains(content, "More content.") {
		t.Error("trailing content was lost")
	}
	if !strings.Contains(content, openspecStartMarker) {
		t.Error("new block not injected")
	}
}

// --- OpenCode Tests ---

func TestOpenCode_Configure(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	o := &OpenCode{}
	if err := o.Configure(dir, osp); err != nil {
		t.Fatal(err)
	}

	for _, cmd := range []string{"proposal", "apply", "archive"} {
		path := filepath.Join(dir, ".opencode", "commands", "openspec", cmd+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("opencode command not created: %s", cmd)
		}
	}
}

func TestOpenCode_Update(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	o := &OpenCode{}
	o.Configure(dir, osp)
	if err := o.UpdateExisting(dir, osp); err != nil {
		t.Fatal(err)
	}
	if !o.IsConfigured(dir) {
		t.Error("not configured after update")
	}
}

// --- Codex Tests ---

func TestCodex_Configure(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	c := &Codex{}
	if err := c.Configure(dir, osp); err != nil {
		t.Fatal(err)
	}

	for _, cmd := range []string{"proposal", "apply", "archive"} {
		path := filepath.Join(dir, ".codex", "prompts", "openspec", cmd+".md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("codex prompt not created: %s", cmd)
		}
	}
}

func TestCodex_Update(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	c := &Codex{}
	c.Configure(dir, osp)
	if err := c.UpdateExisting(dir, osp); err != nil {
		t.Fatal(err)
	}
	if !c.IsConfigured(dir) {
		t.Error("not configured after update")
	}
}

// --- Goose Tests ---

func TestGoose_Configure(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	g := &Goose{}
	if err := g.Configure(dir, osp); err != nil {
		t.Fatal(err)
	}

	// Check .goosehints
	hintsPath := filepath.Join(dir, ".goosehints")
	data, err := os.ReadFile(hintsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), gooseStartMarker) {
		t.Error(".goosehints missing marker")
	}

	// Check recipes
	for _, recipe := range []string{"proposal.yaml", "apply.yaml", "archive.yaml"} {
		path := filepath.Join(dir, ".goose", "recipes", "openspec", recipe)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("recipe not created: %s", recipe)
		}
	}
}

func TestGoose_Update(t *testing.T) {
	dir := t.TempDir()
	osp := filepath.Join(dir, "openspec")
	os.MkdirAll(osp, 0o755)

	g := &Goose{}
	g.Configure(dir, osp)
	if err := g.UpdateExisting(dir, osp); err != nil {
		t.Fatal(err)
	}
	if !g.IsConfigured(dir) {
		t.Error("not configured after update")
	}
}
