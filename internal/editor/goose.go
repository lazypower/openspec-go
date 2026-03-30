package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	gooseStartMarker = "# OPENSPEC:START"
	gooseEndMarker   = "# OPENSPEC:END"
)

// Goose configures Goose with .goosehints and recipe YAML files.
type Goose struct{}

func (g *Goose) Name() string { return "goose" }

func (g *Goose) Configure(projectPath, openspecPath string) error {
	if err := g.writeGooseHints(projectPath); err != nil {
		return err
	}
	return g.writeRecipes(projectPath)
}

func (g *Goose) UpdateExisting(projectPath, openspecPath string) error {
	if err := g.writeGooseHints(projectPath); err != nil {
		return err
	}
	return g.writeRecipes(projectPath)
}

func (g *Goose) IsConfigured(projectPath string) bool {
	_, err := os.Stat(filepath.Join(projectPath, ".goose", "recipes", "openspec"))
	return err == nil
}

func (g *Goose) writeGooseHints(projectPath string) error {
	hintsPath := filepath.Join(projectPath, ".goosehints")
	block := fmt.Sprintf("%s\n# OpenSpec: See @openspec/AGENTS.md for spec-driven workflow instructions\n%s\n", gooseStartMarker, gooseEndMarker)

	existing, err := os.ReadFile(hintsPath)
	if err != nil {
		return os.WriteFile(hintsPath, []byte(block), 0o644)
	}

	content := string(existing)
	startIdx := strings.Index(content, gooseStartMarker)
	endIdx := strings.Index(content, gooseEndMarker)

	if startIdx >= 0 && endIdx >= 0 {
		content = content[:startIdx] + block + content[endIdx+len(gooseEndMarker):]
	} else {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		content += "\n" + block
	}

	return os.WriteFile(hintsPath, []byte(content), 0o644)
}

func (g *Goose) writeRecipes(projectPath string) error {
	recipeDir := filepath.Join(projectPath, ".goose", "recipes", "openspec")
	if err := os.MkdirAll(recipeDir, 0o755); err != nil {
		return err
	}

	recipes := map[string]string{
		"proposal.yaml": `version: 1
title: OpenSpec Proposal
description: Create a new OpenSpec change proposal
instructions: |
  Open @openspec/AGENTS.md and follow the instructions for creating a new change proposal.
`,
		"apply.yaml": `version: 1
title: OpenSpec Apply
description: Implement an approved OpenSpec change
instructions: |
  Open @openspec/AGENTS.md and follow the instructions for implementing an approved change.
`,
		"archive.yaml": `version: 1
title: OpenSpec Archive
description: Archive a completed OpenSpec change
instructions: |
  Open @openspec/AGENTS.md and follow the instructions for archiving a completed change.
`,
	}

	for name, content := range recipes {
		if err := os.WriteFile(filepath.Join(recipeDir, name), []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}
