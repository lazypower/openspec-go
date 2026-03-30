package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chuck/openspec-go/internal/model"
	"github.com/chuck/openspec-go/internal/parser"
)

// findOpenSpecPath walks up from the given path to find an openspec/ directory.
func findOpenSpecPath(startPath string) (string, error) {
	if startPath == "" {
		var err error
		startPath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	abs, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	// Check if the path itself is an openspec directory
	candidate := filepath.Join(abs, "openspec")
	if info, err := os.Stat(candidate); err == nil && info.IsDir() {
		return candidate, nil
	}

	// Walk up
	for {
		parent := filepath.Dir(abs)
		if parent == abs {
			break
		}
		abs = parent
		candidate = filepath.Join(abs, "openspec")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("openspec directory not found")
}

// loadChanges loads all active changes from the openspec directory.
func loadChanges(ospPath string) ([]model.Change, error) {
	changesDir := filepath.Join(ospPath, "changes")
	entries, err := os.ReadDir(changesDir)
	if err != nil {
		return nil, nil
	}

	var changes []model.Change
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "archive" {
			continue
		}

		change, err := loadChange(ospPath, entry.Name())
		if err != nil {
			continue
		}
		changes = append(changes, change)
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ID < changes[j].ID
	})
	return changes, nil
}

// loadChange loads a single change by ID.
func loadChange(ospPath, id string) (model.Change, error) {
	changePath := filepath.Join(ospPath, "changes", id)
	change := model.Change{ID: id}

	// Parse proposal
	proposalData, err := os.ReadFile(filepath.Join(changePath, "proposal.md"))
	if err != nil {
		return change, err
	}
	parsed := parser.ParseChange(string(proposalData))
	change.Title = parsed.Title
	change.Why = parsed.Why
	change.WhatChanges = parsed.WhatChanges
	change.Impact = parsed.Impact

	// Parse tasks
	tasksData, err := os.ReadFile(filepath.Join(changePath, "tasks.md"))
	if err == nil {
		change.Tasks = parser.ParseTaskProgress(string(tasksData))
	}

	// Parse deltas
	specsDir := filepath.Join(changePath, "specs")
	if specEntries, err := os.ReadDir(specsDir); err == nil {
		for _, specEntry := range specEntries {
			if !specEntry.IsDir() {
				continue
			}
			specFile := filepath.Join(specsDir, specEntry.Name(), "spec.md")
			deltaContent, err := os.ReadFile(specFile)
			if err != nil {
				continue
			}
			deltas := parser.ParseDeltas(string(deltaContent))
			for i := range deltas {
				deltas[i].SpecName = specEntry.Name()
			}
			change.Deltas = append(change.Deltas, deltas...)
		}
	}

	return change, nil
}

// loadSpecs loads all specs from the openspec directory.
func loadSpecs(ospPath string) (map[string]model.Spec, error) {
	specsDir := filepath.Join(ospPath, "specs")
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil, nil
	}

	specs := make(map[string]model.Spec)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		specFile := filepath.Join(specsDir, entry.Name(), "spec.md")
		data, err := os.ReadFile(specFile)
		if err != nil {
			continue
		}
		specs[entry.Name()] = parser.ParseSpec(string(data))
	}
	return specs, nil
}

// loadArchivedChanges returns IDs of archived changes.
func loadArchivedChanges(ospPath string) []string {
	archiveDir := filepath.Join(ospPath, "changes", "archive")
	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		return nil
	}
	var ids []string
	for _, e := range entries {
		if e.IsDir() {
			ids = append(ids, e.Name())
		}
	}
	sort.Strings(ids)
	return ids
}

// suggestNearestMatch returns the closest matching item name.
func suggestNearestMatch(target string, candidates []string) string {
	target = strings.ToLower(target)
	var best string
	bestScore := 0
	for _, c := range candidates {
		cl := strings.ToLower(c)
		score := 0
		if strings.Contains(cl, target) || strings.Contains(target, cl) {
			score = 10
		}
		// Simple prefix match scoring
		for i := 0; i < len(target) && i < len(cl); i++ {
			if target[i] == cl[i] {
				score++
			} else {
				break
			}
		}
		if score > bestScore {
			bestScore = score
			best = c
		}
	}
	return best
}
